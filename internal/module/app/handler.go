package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/app", func(r chi.Router) {
		r.Get("/config", serve(h.config))
		r.Get("/system-info", serve(h.systemInfo))
	})
}

func (h *Handler) config(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, h.service.Config())
}

func (h *Handler) systemInfo(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, h.service.SystemInfo())
}
