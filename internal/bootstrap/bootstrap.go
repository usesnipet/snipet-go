package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riverqueue/river"
	"github.com/usesnipet/snipet/config"
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
	auth_module "github.com/usesnipet/snipet/internal/module/auth"
	auth_provider "github.com/usesnipet/snipet/internal/module/auth/auth-provider"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/module/knowledge"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	"github.com/usesnipet/snipet/internal/module/session"
	"github.com/usesnipet/snipet/internal/module/user"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/web"
)

func Bootstrap(cfg *config.Config, logger *logger.Logger) error {
	// database
	db, sqlDB, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Errorf("failed to create database: %v", err)
		return err
	}

	// repositories
	txManager := repository.NewTxManager(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	clientRepo := repository.NewClientRepository(db)
	sessionRepo := repository.NewSessionRepository(db, clientRepo)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	knowledgeIndexRepo := repository.NewKnowledgeIndexRepository(db)
	knowledgeItemRepo := repository.NewKnowledgeItemRepository(db)
	indexedKnowledgeItemRepo := repository.NewIndexedKnowledgeItemRepository(db)
	userRepo := repository.NewUserRepository(db, clientRepo)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	executionRepo := repository.NewExecutionRepository(db)
	messageRepo := repository.NewExecutionMessageRepository(db)

	// runtime
	sourceRegistry := source.Registry()
	sourceManager := driver.NewManager(sourceRegistry)

	indexRegistry := index.Registry()
	indexManager := driver.NewManager(indexRegistry)

	llmRegistry := llm.Registry()
	llmManager := driver.NewManager(llmRegistry)

	toolRegistry := tool.Registry()
	toolManager := driver.NewManager(toolRegistry)

	engine := runtime.NewEngine(
		toolManager,
		llmManager,
		logger,
	)

	workers := river.NewWorkers()
	river.AddWorker(
		workers,
		knowledge.NewSyncWorker(txManager, sourceManager, knowledgeRepo, knowledgeItemRepo, knowledgeIndexRepo, 100, logger),
	)
	riverClient, err := queue.NewRiver(sqlDB, workers)
	if err != nil {
		logger.Errorf("failed to create river client: %v", err)
		return err
	}

	// services
	apiKeyGenerator := auth.NewAPIKeyGenerator()
	apiKeyHasher := auth.NewKeyHasher()
	jwtService := auth.NewJWTService(cfg.Auth)

	authRegistry := auth_provider.NewRegistry()

	authService := auth_module.NewService(authRegistry, clientRepo, userRepo, refreshTokenRepo, jwtService, cfg.Auth)

	apiKeyService := apikey.NewService(logger, apiKeyRepo, apiKeyGenerator, apiKeyHasher)
	apiKeyService.Init(context.Background())

	clientService := client.NewService(clientRepo, logger)
	if err := clientService.Init(context.Background(), &cfg.App); err != nil {
		logger.Errorf("failed to init client service: %v", err)
		return err
	}

	agentService := agent.NewService(agentRepo, engine, llmManager, toolManager, executionRepo, messageRepo, logger)

	knowledgeService := knowledge.NewService(txManager, knowledgeRepo, knowledgeItemRepo, sourceManager, riverClient)
	knowledgeIndexService := knowledgeindex.NewService(knowledgeIndexRepo, indexedKnowledgeItemRepo, indexManager, riverClient, txManager)

	sessionService := session.NewService(sessionRepo, messageRepo, clientService, agentService)

	userService := user.NewService(userRepo)
	appService := app_module.NewService(&cfg.App)

	// cache
	apiKeyCache := cache.NewMemoryCache(1000, 1*time.Hour)

	// middlewares
	apiKeyMiddleware := middleware.APIKeyMiddleware(apiKeyService, apiKeyCache)
	anyAuthMiddleware := middleware.AnyAuth(jwtService, apiKeyService, apiKeyCache)

	// handlers
	authHandler := auth_module.NewHandler(authService)
	appHandler := app_module.NewHandler(appService)
	apiKeyHandler := apikey.NewHandler(apiKeyService, apiKeyMiddleware)
	agentHandler := agent.NewHandler(agentService, apiKeyMiddleware)
	clientHandler := client.NewHandler(clientService, apiKeyMiddleware)
	sessionHandler := session.NewHandler(sessionService, anyAuthMiddleware)
	knowledgeHandler := knowledge.NewHandler(knowledgeService, apiKeyMiddleware)
	knowledgeIndexHandler := knowledgeindex.NewHandler(knowledgeIndexService, apiKeyMiddleware)
	userHandler := user.NewHandler(userService, apiKeyMiddleware, anyAuthMiddleware)

	// register routes
	api := api.New()
	api.Router.Handle("/*", web.Handler())
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		authHandler.RegisterRoutes(r, api.Serve)
		appHandler.RegisterRoutes(r, api.Serve)
		apiKeyHandler.RegisterRoutes(r, api.Serve)
		agentHandler.RegisterRoutes(r, api.Serve)
		clientHandler.RegisterRoutes(r, api.Serve)
		sessionHandler.RegisterRoutes(r, api.Serve)
		knowledgeHandler.RegisterRoutes(r, api.Serve)
		knowledgeIndexHandler.RegisterRoutes(r, api.Serve)
		userHandler.RegisterRoutes(r, api.Serve)
	})

	err = riverClient.Start(context.Background())
	if err != nil {
		logger.Errorf("failed to start river client: %v", err)
		return err
	}

	logger.Infof("server started on port %d", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), api.Router)
	if err != nil {
		logger.Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
