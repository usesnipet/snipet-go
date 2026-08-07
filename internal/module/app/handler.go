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

// @Summary		Get app config
// @Description	Returns public application configuration.
// @Tags			app
// @Produce		json
// @Success		200	{object}	AppConfigDTO
// @Router			/app/config [get]
func (h *Handler) config(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, h.service.Config())
}

// @Summary		Get system info
// @Description	Returns application system information.
// @Tags			app
// @Produce		json
// @Success		200	{object}	SystemInfoDTO
// @Router			/app/system-info [get]
func (h *Handler) systemInfo(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, h.service.SystemInfo())
}
