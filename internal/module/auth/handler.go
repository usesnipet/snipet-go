package auth

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	auth_provider "github.com/usesnipet/snipet/internal/module/auth/auth-provider"
	"github.com/usesnipet/snipet/pkg/jsonx"
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
		r.Post("/refresh", serve(h.refresh))
	})
}

func (h *Handler) clientCode(r *http.Request) string {
	return chi.URLParam(r, "client_code")
}

func requestMetadata(r *http.Request) jsonx.JSONMap {
	ip := clientIP(r)
	return jsonx.JSONMap{
		"ip":         ip,
		"user_agent": r.UserAgent(),
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) error {
	providerName := auth_provider.ProviderName(chi.URLParam(r, "provider_name"))
	res, err := h.service.Authenticate(r.Context(), h.clientCode(r), providerName, r, requestMetadata(r))
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
	res, err := h.service.AuthenticateAnonymous(r.Context(), h.clientCode(r), dto, requestMetadata(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) error {
	var dto RefreshDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	res, err := h.service.Refresh(r.Context(), h.clientCode(r), dto, requestMetadata(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, res)
}
