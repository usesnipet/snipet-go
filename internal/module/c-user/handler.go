package c_user

import (
	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/module/client"
)

type Handler struct {
	service       *Service
	clientService *client.Service
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_code}/user", func(r chi.Router) {

	})
}
