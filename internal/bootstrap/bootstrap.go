package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/config"
	_ "github.com/usesnipet/snipet/docs/swagger"
	"github.com/usesnipet/snipet/drivers/index"
	"github.com/usesnipet/snipet/drivers/llm"
	"github.com/usesnipet/snipet/drivers/source"
	"github.com/usesnipet/snipet/drivers/tool"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/middleware"
	"github.com/usesnipet/snipet/internal/module/agent"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
	app_module "github.com/usesnipet/snipet/internal/module/app"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/module/clientauth"
	auth_provider "github.com/usesnipet/snipet/internal/module/clientauth/auth-provider"
	"github.com/usesnipet/snipet/internal/module/clientuser"
	"github.com/usesnipet/snipet/internal/module/knowledge"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	llmmodule "github.com/usesnipet/snipet/internal/module/llm"
	"github.com/usesnipet/snipet/internal/module/session"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/web"
)

func Bootstrap(cfg *config.Config, logger *logger.Logger) error {
	// database
	db, _, embeddedDB, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Errorf("failed to create database: %v", err)
		return err
	}

	if embeddedDB != nil {
		defer func() {
			logger.Infof("stopping embedded database...")
			if err := embeddedDB.Stop(); err != nil {
				logger.Errorf("failed to stop embedded database: %v", err)
				return
			}
			logger.Infof("embedded database stopped successfully")
		}()
	}

	// repositories
	txManager := repository.NewTxManager(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	llmRepo := repository.NewLLMRepository(db)
	clientRepo := repository.NewClientRepository(db)
	sessionRepo := repository.NewSessionRepository(db, clientRepo)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	knowledgeIndexRepo := repository.NewKnowledgeIndexRepository(db)
	knowledgeItemRepo := repository.NewKnowledgeItemRepository(db)
	indexedKnowledgeItemRepo := repository.NewIndexedKnowledgeItemRepository(db)
	clientUserRepo := repository.NewClientUserRepository(db, clientRepo)
	refreshTokenRepo := repository.NewClientUserRefreshTokenRepository(db)
	executionRepo := repository.NewExecutionRepository(db)
	messageRepo := repository.NewExecutionMessageRepository(db)

	// runtime
	sourceRegistry := source.Registry(logger)
	sourceManager := manager.NewDriver(sourceRegistry)

	indexRegistry := index.Registry(logger)
	indexManager := manager.NewDriver(indexRegistry)

	llmRegistry := llm.Registry(logger)
	llmManager := manager.NewDriver(llmRegistry)

	toolRegistry := tool.Registry(logger)
	toolManager := manager.NewTool(manager.NewDriver(toolRegistry))

	engine := runtime.NewEngine(llmManager, toolManager, logger)

	syncPool := queue.NewPool(cfg.Sync.Workers, logger)
	syncPool.Start(context.Background())

	indexSyncWorker := knowledgeindex.NewSyncIndexWorker(
		indexManager,
		sourceManager,
		knowledgeRepo,
		knowledgeItemRepo,
		knowledgeIndexRepo,
		indexedKnowledgeItemRepo,
		logger,
	)
	knowledgeSyncWorker := knowledge.NewSyncWorker(
		sourceManager,
		knowledgeRepo,
		knowledgeItemRepo,
		knowledgeIndexRepo,
		indexSyncWorker.Sync,
		100,
		logger,
	)

	// services
	apiKeyGenerator := auth.NewAPIKeyGenerator()
	apiKeyHasher := auth.NewKeyHasher()
	jwtService := auth.NewJWTService(cfg.Auth, func() *clientauth.UserClaims { return &clientauth.UserClaims{} })

	refreshTokenService := auth.NewTokenService()

	authRegistry := auth_provider.NewRegistry()

	authService := clientauth.NewService(
		authRegistry,
		clientRepo,
		clientUserRepo,
		refreshTokenRepo,
		jwtService,
		refreshTokenService,
		cfg.Auth,
	)

	apiKeyService := apikey.NewService(logger, apiKeyRepo, apiKeyGenerator, apiKeyHasher)
	apiKeyService.Init(context.Background())

	clientService := client.NewService(clientRepo, agentRepo, logger)
	if err := clientService.Init(context.Background(), &cfg.App); err != nil {
		logger.Errorf("failed to init client service: %v", err)
		return err
	}

	agentService := agent.NewService(agentRepo, llmRepo, txManager, engine, executionRepo, messageRepo, logger)
	llmService := llmmodule.NewService(llmRepo, llmManager)

	knowledgeService := knowledge.NewService(
		txManager,
		knowledgeRepo,
		knowledgeItemRepo,
		sourceManager,
		syncPool,
		knowledgeSyncWorker,
	)
	knowledgeIndexService := knowledgeindex.NewService(
		knowledgeIndexRepo,
		indexedKnowledgeItemRepo,
		indexManager,
		syncPool,
		indexSyncWorker,
		txManager,
	)

	sessionService := session.NewService(sessionRepo, messageRepo, clientService, agentService)

	clientUserService := clientuser.NewService(clientUserRepo)
	appService := app_module.NewService(&cfg.App)

	// cache
	apiKeyCache := cache.NewMemoryCache(1000, 1*time.Hour)

	// middlewares
	apiKeyMiddleware := middleware.APIKeyMiddleware(apiKeyService, apiKeyCache)
	anyAuthMiddleware := middleware.AnyAuth(jwtService, apiKeyService, apiKeyCache)
	jwtMiddleware := middleware.JWT(jwtService)

	// handlers
	clientAuthHandler := clientauth.NewHandler(authService)
	appHandler := app_module.NewHandler(appService)
	apiKeyHandler := apikey.NewHandler(apiKeyService, apiKeyMiddleware)
	agentHandler := agent.NewHandler(agentService, apiKeyMiddleware)
	llmHandler := llmmodule.NewHandler(llmService, apiKeyMiddleware)
	clientHandler := client.NewHandler(clientService, apiKeyMiddleware, anyAuthMiddleware)
	sessionHandler := session.NewHandler(sessionService, anyAuthMiddleware)
	knowledgeHandler := knowledge.NewHandler(knowledgeService, apiKeyMiddleware)
	knowledgeIndexHandler := knowledgeindex.NewHandler(knowledgeIndexService, apiKeyMiddleware)
	clientUserHandler := clientuser.NewHandler(clientUserService, apiKeyMiddleware, anyAuthMiddleware, jwtMiddleware)

	// register routes
	api := api.New()
	api.Router.Handle("/*", web.Handler())
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		clientAuthHandler.RegisterRoutes(r, api.Serve)
		appHandler.RegisterRoutes(r, api.Serve)
		apiKeyHandler.RegisterRoutes(r, api.Serve)
		agentHandler.RegisterRoutes(r, api.Serve)
		llmHandler.RegisterRoutes(r, api.Serve)
		clientHandler.RegisterRoutes(r, api.Serve)
		sessionHandler.RegisterRoutes(r, api.Serve)
		knowledgeHandler.RegisterRoutes(r, api.Serve)
		knowledgeIndexHandler.RegisterRoutes(r, api.Serve)
		clientUserHandler.RegisterRoutes(r, api.Serve)
	})

	logger.Infof("sync worker pool started with %d workers", cfg.Sync.Workers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: api.Router,
	}

	go func() {
		logger.Infof("server started on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("failed to start server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Infof("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("failed to shutdown server: %v", err)
	}

	return nil
}
