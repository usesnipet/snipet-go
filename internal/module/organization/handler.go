package organization

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service *Service
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/organization", func(r chi.Router) {
		r.Get("/", serve(h.findBy))
		r.Post("/", serve(h.create))
	})
}

func (h *Handler) findBy(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindBy(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateOrganizationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	return h.service.Create(r.Context(), dto)
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}
