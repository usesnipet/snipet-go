package apikey

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service *Service
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/api-key", func(r chi.Router) {
		r.Get("/", serve(h.filterBy))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Post("/{id}/roll", serve(h.roll))
		r.Put("/{id}/expiration", serve(h.updateExpiration))
		r.Put("/{id}/active", serve(h.toggleActive))
		r.Put("/{id}/disabled", serve(h.toggleDisabled))
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
	var dto CreateAPIKeyDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) roll(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Roll(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) updateExpiration(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateExpirationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.UpdateExpiration(r.Context(), chi.URLParam(r, "id"), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) toggleActive(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), chi.URLParam(r, "id"), true); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) toggleDisabled(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), chi.URLParam(r, "id"), false); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}
