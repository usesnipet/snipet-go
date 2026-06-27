package c_auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	auth_provider "github.com/usesnipet/snipet/internal/module/c-auth/auth-provider"
)

type Handler struct {
	service *Service
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_code}", func(r chi.Router) {
		r.Post("/authenticate/{provider_name}", serve(h.authenticate))
	})
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) error {
	clientCode := chi.URLParam(r, "client_code")
	providerName := auth_provider.ProviderName(chi.URLParam(r, "provider_name"))
	res, err := h.service.Authenticate(r.Context(), clientCode, providerName, r)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, res)
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}
