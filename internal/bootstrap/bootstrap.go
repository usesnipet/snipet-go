package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/infra/cache"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/middleware"
	apikey "github.com/usesnipet/snipet/internal/module/api-key"
	"github.com/usesnipet/snipet/internal/module/bot"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/module/memory"
	"github.com/usesnipet/snipet/internal/module/session"
	"github.com/usesnipet/snipet/internal/repository"
)

func Bootstrap(cfg *config.Config, logger *logger.Logger) error {
	// database
	db, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Errorf("failed to create database: %v", err)
		return err
	}

	// repository
	apiKeyRepo := repository.NewApiKeyRepository(db)
	botRepo := repository.NewBotRepository(db)
	clientRepo := repository.NewClientRepository(db)
	sessionRepo := repository.NewSessionRepository(db, clientRepo)
	sessionMessageRepo := repository.NewSessionMessageRepository(db, clientRepo)
	memoryRepo := repository.NewMemoryRepository(db)

	// service
	apiKeyService := apikey.NewService(logger, apiKeyRepo)
	apiKeyService.Init(context.Background())
	botService := bot.NewService(botRepo, clientRepo)
	clientService := client.NewService(clientRepo)
	memoryService := memory.NewService(memoryRepo)
	sessionService := session.NewService(sessionRepo, sessionMessageRepo, memoryRepo)

	// middleware
	apiKeyMiddleware := middleware.APIKeyMiddleware(
		apiKeyService,
		cache.NewMemoryCache(1000, 1*time.Hour),
	)

	// handler
	apiKeyHandler := apikey.NewHandler(apiKeyService, apiKeyMiddleware)
	botHandler := bot.NewHandler(botService, apiKeyMiddleware)
	clientHandler := client.NewHandler(clientService, apiKeyMiddleware)
	sessionHandler := session.NewHandler(sessionService, apiKeyMiddleware)
	memoryHandler := memory.NewHandler(memoryService, apiKeyMiddleware)

	// register handlers
	api := api.New()
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		apiKeyHandler.RegisterRoutes(r, api.Serve)
		botHandler.RegisterRoutes(r, api.Serve)
		clientHandler.RegisterRoutes(r, api.Serve)
		sessionHandler.RegisterRoutes(r, api.Serve)
		memoryHandler.RegisterRoutes(r, api.Serve)
	})

	logger.Infof("server started on port %d", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), api.Router)
	if err != nil {
		logger.Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
