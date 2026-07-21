package client

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service          *Service
	apiKeyMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, apiKeyMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, apiKeyMiddleware: apiKeyMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	// Use /clients (plural) so CRUD does not conflict with /client/{client_code}/...
	// routes used by auth, session, and user modules.
	r.Route("/clients", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{code}", serve(h.findByCode))
			r.Put("/{code}", serve(h.update))
			r.Delete("/{code}", serve(h.delete))
		})
	})
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindClientsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}

	clients, err := h.service.Filter(r.Context(), query.ToFilter())
	if err != nil {
		return err
	}

	return api.WriteJSON(w, http.StatusOK, clients)
}

func (h *Handler) findByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateClientDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateClientDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.UpdateByCode(r.Context(), chi.URLParam(r, "code"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByCode(r.Context(), chi.URLParam(r, "code")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
