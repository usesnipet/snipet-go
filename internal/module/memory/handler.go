package memory

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
	r.Route("/memory", func(r chi.Router) {
		r.Get("/", serve(h.filterBy))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
		r.Post("/default", serve(h.setAsDefault))
	})
}

func (h *Handler) filterBy(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FilterBy(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateMemoryDTO
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
	var dto UpdateMemoryDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), chi.URLParam(r, "id"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), chi.URLParam(r, "id")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) setAsDefault(w http.ResponseWriter, r *http.Request) error {
	var dto SetAsDefaultMemoryDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.SetAsDefault(r.Context(), dto)
	if err != nil {
		return err
	}

	return api.WriteNoContent(w)
}
