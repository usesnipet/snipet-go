package session

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service           *Service
	anyAuthMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, anyAuthMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, anyAuthMiddleware: anyAuthMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_code}/session", func(r chi.Router) {
		r.Use(h.anyAuthMiddleware)
		r.Get("/", serve(h.filter))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Delete("/{id}", serve(h.deleteByID))
		r.Get("/{id}/messages", serve(h.findMessages))
	})
}

func (h *Handler) clientCode(r *http.Request) string {
	return chi.URLParam(r, "client_code")
}

func (h *Handler) sessionID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) findMessages(w http.ResponseWriter, r *http.Request) error {
	query := &FindMessagesFilterDTO{}
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FindMessages(r.Context(), h.clientCode(r), h.sessionID(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Filter(r.Context(), h.clientCode(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.clientCode(r), h.sessionID(r))
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
	data, err := h.service.Create(r.Context(), h.clientCode(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.clientCode(r), h.sessionID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
