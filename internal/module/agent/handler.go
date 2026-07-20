package agent

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/runtime"
)

type Handler struct {
	service          *Service
	apiKeyMiddleware api.MiddlewareFunc
}

func NewHandler(
	service *Service,
	apiKeyMiddleware api.MiddlewareFunc,
) api.Handler {
	return &Handler{
		service:          service,
		apiKeyMiddleware: apiKeyMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/agent", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Put("/{id}", serve(h.update))
			r.Delete("/{id}", serve(h.deleteByID))
			r.Post("/{id}/run", serve(h.run))
		})
	})
}
func (h *Handler) agentID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Filter(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.agentID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateAgentDTO
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
	var dto UpdateAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), h.agentID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.agentID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) run(w http.ResponseWriter, r *http.Request) error {
	var dto RunAgentDTO
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

	err := h.service.Run(r.Context(), RunInput{
		AgentID: h.agentID(r),
		Message: dto.Message,
	}, func(event runtime.IEvent) error {
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
