package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/module/bot"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/module/conversation"
)

func Bootstrap(cfg *config.Config, logger *logger.Logger) error {
	// database
	db, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Errorf("failed to create database: %v", err)
		return err
	}

	// repository
	botRepo := bot.NewRepository(db)
	clientRepo := client.NewRepository(db)
	conversationRepo := conversation.NewRepository(db)

	// service
	botService := bot.NewService(botRepo, clientRepo)
	clientService := client.NewService(clientRepo)
	conversationService := conversation.NewService(conversationRepo)

	// handler
	botHandler := bot.NewHandler(botService)
	clientHandler := client.NewHandler(clientService)
	conversationHandler := conversation.NewHandler(conversationService)

	// register handlers
	api := api.New()
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		botHandler.RegisterRoutes(r, api.Serve)
		clientHandler.RegisterRoutes(r, api.Serve)
		conversationHandler.RegisterRoutes(r, api.Serve)
	})

	logger.Infof("server started on port %d", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), api.Router)
	if err != nil {
		logger.Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
