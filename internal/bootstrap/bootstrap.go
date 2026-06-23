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

	// service
	botService := bot.NewService(botRepo)

	// handler
	botHandler := bot.NewHandler(botService)

	// register handlers
	api := api.New()
	api.Router.Route(config.APIPrefix, func(r chi.Router) {
		botHandler.RegisterRoutes(r, api.Serve)
	})

	logger.Infof("server started on port %d", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), api.Router)
	if err != nil {
		logger.Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
