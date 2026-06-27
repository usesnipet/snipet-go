package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/config"
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
	"github.com/usesnipet/snipet/internal/module/memory"
	"github.com/usesnipet/snipet/internal/module/session"
	"github.com/usesnipet/snipet/internal/module/user"
	"github.com/usesnipet/snipet/internal/repository"
)

func Bootstrap(cfg *config.Config, logger *logger.Logger) error {
	// database
	db, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Errorf("failed to create database: %v", err)
		return err
	}

	// repositories
	apiKeyRepo := repository.NewApiKeyRepository(db)
	botRepo := repository.NewBotRepository(db)
	clientRepo := repository.NewClientRepository(db)
	sessionRepo := repository.NewSessionRepository(db, clientRepo)
	sessionMessageRepo := repository.NewSessionMessageRepository(db, clientRepo)
	memoryRepo := repository.NewMemoryRepository(db)
	userRepo := repository.NewUserRepository(db, clientRepo)

	// services
	apiKeyGenerator := auth.NewAPIKeyGenerator()
	apiKeyHasher := auth.NewKeyHasher()
	jwtService := auth.NewJWTService(cfg.Auth)

	authRegistry := auth_provider.NewRegistry()

	authService := auth_module.NewService(authRegistry, clientRepo, userRepo, jwtService)
	apiKeyService := apikey.NewService(logger, apiKeyRepo, apiKeyGenerator, apiKeyHasher)
	apiKeyService.Init(context.Background())
	botService := bot.NewService(botRepo, clientRepo)
	clientService := client.NewService(clientRepo)
	memoryService := memory.NewService(memoryRepo)
	sessionService := session.NewService(sessionRepo, sessionMessageRepo, memoryRepo)
	userService := user.NewService(userRepo)

	// cache
	apiKeyCache := cache.NewMemoryCache(1000, 1*time.Hour)

	// middlewares
	apiKeyMiddleware := middleware.APIKeyMiddleware(apiKeyService, apiKeyCache)
	anyAuthMiddleware := middleware.AnyAuth(jwtService, apiKeyService, apiKeyCache)

	// handlers
	authHandler := auth_module.NewHandler(authService)
	apiKeyHandler := apikey.NewHandler(apiKeyService, apiKeyMiddleware)
	botHandler := bot.NewHandler(botService, apiKeyMiddleware)
	clientHandler := client.NewHandler(clientService, apiKeyMiddleware)
	sessionHandler := session.NewHandler(sessionService, apiKeyMiddleware)
	memoryHandler := memory.NewHandler(memoryService, apiKeyMiddleware)
	userHandler := user.NewHandler(userService, apiKeyMiddleware, anyAuthMiddleware)

	// register routes
	api := api.New()
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		authHandler.RegisterRoutes(r, api.Serve)
		apiKeyHandler.RegisterRoutes(r, api.Serve)
		botHandler.RegisterRoutes(r, api.Serve)
		clientHandler.RegisterRoutes(r, api.Serve)
		sessionHandler.RegisterRoutes(r, api.Serve)
		memoryHandler.RegisterRoutes(r, api.Serve)
		userHandler.RegisterRoutes(r, api.Serve)
	})

	logger.Infof("server started on port %d", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), api.Router)
	if err != nil {
		logger.Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
