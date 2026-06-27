package session

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
	r.Route("/client/{client_code}/session", func(r chi.Router) {
		r.Get("/", serve(h.findBy))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Delete("/{id}", serve(h.deleteByID))
		r.Get("/{id}/messages", serve(h.findMessages))
	})
}

func (h *Handler) findMessages(w http.ResponseWriter, r *http.Request) error {
	clientCode := chi.URLParam(r, "client_code")
	sessionID := chi.URLParam(r, "id")
	var filter FindMessagesFilterDTO
	if err := api.ParseQuery(r, &filter); err != nil {
		return err
	}
	data, err := h.service.FindMessages(r.Context(), clientCode, sessionID, filter)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findBy(w http.ResponseWriter, r *http.Request) error {
	clientCode := chi.URLParam(r, "client_code")

	data, err := h.service.FindBy(r.Context(), clientCode)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	clientCode := chi.URLParam(r, "client_code")
	sessionID := chi.URLParam(r, "id")
	data, err := h.service.FindByID(r.Context(), clientCode, sessionID)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateSessionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), chi.URLParam(r, "client_code"), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	clientCode := chi.URLParam(r, "client_code")
	sessionID := chi.URLParam(r, "id")
	if err := h.service.DeleteByID(r.Context(), clientCode, sessionID); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
