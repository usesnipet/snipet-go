package llm

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
	return &Handler{
		service:          service,
		apiKeyMiddleware: apiKeyMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/llm", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Get("/", serve(h.filter))
			r.Get("/drivers", serve(h.listDrivers))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Put("/{id}", serve(h.update))
			r.Delete("/{id}", serve(h.deleteByID))
		})
	})
}

func (h *Handler) llmID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindLLMsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.llmID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateLLMDTO
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
	var dto UpdateLLMDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.Update(r.Context(), h.llmID(r), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.llmID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) listDrivers(w http.ResponseWriter, r *http.Request) error {
	drivers, err := h.service.ListDrivers(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, drivers)
}
