package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	auth_provider "github.com/usesnipet/snipet/internal/module/auth/auth-provider"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_code}", func(r chi.Router) {
		r.Post("/authenticate/{provider_name}", serve(h.authenticate))
		r.Post("/authenticate/anonymous", serve(h.authenticateAnonymous))
	})
}

func (h *Handler) clientCode(r *http.Request) string {
	return chi.URLParam(r, "client_code")
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) error {
	providerName := auth_provider.ProviderName(chi.URLParam(r, "provider_name"))
	res, err := h.service.Authenticate(r.Context(), h.clientCode(r), providerName, r)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) authenticateAnonymous(w http.ResponseWriter, r *http.Request) error {
	var dto AuthenticateAnonymousDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	res, err := h.service.AuthenticateAnonymous(r.Context(), h.clientCode(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, res)
}
