package session

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	apperr "github.com/usesnipet/snipet/internal/app-err"
)

type Handler struct {
	service           *Service
	anyAuthMiddleware api.MiddlewareFunc
	jwtAuthMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, anyAuthMiddleware api.MiddlewareFunc, jwtAuthMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, anyAuthMiddleware: anyAuthMiddleware, jwtAuthMiddleware: jwtAuthMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_code}/session", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.anyAuthMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Delete("/{id}", serve(h.deleteByID))
			r.Get("/{id}/messages", serve(h.findMessages))
		})
		r.Group(func(r chi.Router) {
			r.Use(h.jwtAuthMiddleware)
			r.Get("/{id}/send", serve(h.sendMessage))
		})
	})
}

func (h *Handler) clientCode(r *http.Request) string {
	return chi.URLParam(r, "client_code")
}

func (h *Handler) sessionID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) sendMessage(w http.ResponseWriter, r *http.Request) error {
	var dto SendMessageDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return apperr.InternalServerError("streaming unsupported")
	}

	fmt.Fprintf(w, "event: message\ndata: %s\n\n", dto.Message)
	flusher.Flush()

	fmt.Fprintf(w, "event: close\ndata: %s\n\n", "done")
	flusher.Flush()

	return nil
}

func (h *Handler) findMessages(w http.ResponseWriter, r *http.Request) error {
	query := &MessagesFilterDTO{}
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
	var query SessionsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.clientCode(r), query.ToFilter())
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
