package session

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/runtime"
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
		r.Group(func(r chi.Router) {
			r.Use(h.anyAuthMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Delete("/{id}", serve(h.deleteByID))
			r.Get("/{id}/messages", serve(h.findMessages))
			r.Post("/{id}/run", serve(h.run))
		})
	})
}

func (h *Handler) clientCode(r *http.Request) string {
	return chi.URLParam(r, "client_code")
}

func (h *Handler) sessionID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) run(w http.ResponseWriter, r *http.Request) error {
	var dto RunSessionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}

	var sse *api.SSEWriter
	ensureSSE := func() error {
		if sse != nil {
			return nil
		}
		var err error
		sse, err = api.NewSSEWriter(w)
		return err
	}

	err := h.service.Run(r.Context(), h.clientCode(r), h.sessionID(r), dto, func(event runtime.IEvent) error {
		if err := ensureSSE(); err != nil {
			return err
		}
		switch event := event.(type) {
		case runtime.ExecutionStatusChangedEvent:
			return sse.Write("status_changed", event)
		case runtime.ExecutionMessageAddedEvent:
			return sse.Write("message_added", event)
		default:
			return nil
		}
	})
	if err != nil {
		if sse == nil {
			return err
		}
		_ = sse.Write("error", map[string]string{"message": err.Error()})
		return nil
	}

	if err := ensureSSE(); err != nil {
		return err
	}
	return sse.Write("close", map[string]string{"status": "done"})
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
	var query SessionIncludeDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FindByID(r.Context(), h.clientCode(r), h.sessionID(r), query.ToFilter())
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
