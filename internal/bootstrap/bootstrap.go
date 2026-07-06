package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riverqueue/river"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/drivers/index/rag"
	fsdriver "github.com/usesnipet/snipet/drivers/source/fs"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/infra/cache"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/middleware"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
	auth_module "github.com/usesnipet/snipet/internal/module/auth"
	auth_provider "github.com/usesnipet/snipet/internal/module/auth/auth-provider"
	"github.com/usesnipet/snipet/internal/module/bot"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/module/knowledge"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	"github.com/usesnipet/snipet/internal/module/session"
	"github.com/usesnipet/snipet/internal/module/user"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
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
	botRepo := repository.NewBotRepository(db)
	clientRepo := repository.NewClientRepository(db)
	sessionRepo := repository.NewSessionRepository(db, clientRepo)
	sessionMessageRepo := repository.NewSessionMessageRepository(db, clientRepo)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	knowledgeIndexRepo := repository.NewKnowledgeIndexRepository(db)
	knowledgeItemRepo := repository.NewKnowledgeItemRepository(db)
	indexedKnowledgeItemRepo := repository.NewIndexedKnowledgeItemRepository(db)
	userRepo := repository.NewUserRepository(db, clientRepo)

	// runtime
	sourceRegistry := runtime.NewRegistry[runtime.ISourceDriver]()
	sourceRegistry.MustRegister("fs", fsdriver.NewDriver())
	sourceManager := runtime.NewSourceManager(sourceRegistry)

	indexRegistry := runtime.NewRegistry[runtime.IIndexDriver]()
	indexRegistry.MustRegister("rag", rag.NewDriver())
	indexManager := runtime.NewIndexManager(indexRegistry)

	workers := river.NewWorkers()
	river.AddWorker(
		workers,
		knowledge.NewSyncWorker(sourceManager, knowledgeRepo, knowledgeItemRepo, 100, logger),
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

	authService := auth_module.NewService(authRegistry, clientRepo, userRepo, jwtService)

	apiKeyService := apikey.NewService(logger, apiKeyRepo, apiKeyGenerator, apiKeyHasher)
	apiKeyService.Init(context.Background())

	clientService := client.NewService(clientRepo)

	botService := bot.NewService(botRepo)

	knowledgeService := knowledge.NewService(txManager, knowledgeRepo, knowledgeItemRepo, sourceManager, riverClient)
	knowledgeIndexService := knowledgeindex.NewService(knowledgeIndexRepo, indexedKnowledgeItemRepo, indexManager, riverClient, txManager)

	sessionService := session.NewService(sessionRepo, sessionMessageRepo, clientService)

	userService := user.NewService(userRepo)

	// cache
	apiKeyCache := cache.NewMemoryCache(1000, 1*time.Hour)

	// middlewares
	apiKeyMiddleware := middleware.APIKeyMiddleware(apiKeyService, apiKeyCache)
	jwtAuthMiddleware := middleware.JWT(jwtService)
	anyAuthMiddleware := middleware.AnyAuth(jwtService, apiKeyService, apiKeyCache)

	// handlers
	authHandler := auth_module.NewHandler(authService)
	apiKeyHandler := apikey.NewHandler(apiKeyService, apiKeyMiddleware)
	botHandler := bot.NewHandler(botService, apiKeyMiddleware)
	clientHandler := client.NewHandler(clientService, apiKeyMiddleware)
	sessionHandler := session.NewHandler(sessionService, anyAuthMiddleware, jwtAuthMiddleware)
	knowledgeHandler := knowledge.NewHandler(knowledgeService, apiKeyMiddleware)
	knowledgeIndexHandler := knowledgeindex.NewHandler(knowledgeIndexService, apiKeyMiddleware)
	userHandler := user.NewHandler(userService, apiKeyMiddleware, anyAuthMiddleware)

	// register routes
	api := api.New()
	api.Router.Handle("/*", web.Handler())
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		authHandler.RegisterRoutes(r, api.Serve)
		apiKeyHandler.RegisterRoutes(r, api.Serve)
		botHandler.RegisterRoutes(r, api.Serve)
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
