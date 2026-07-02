package knowledgeitem

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
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
	r.Route("/knowledge/{knowledge_id}/items", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Put("/{id}", serve(h.update))
			r.Delete("/{id}", serve(h.deleteByID))
		})
	})
}
func (h *Handler) knowledgeID(r *http.Request) string {
	return chi.URLParam(r, "knowledge_id")
}
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var dto FilterKnowledgeItemDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.knowledgeID(r), dto.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateKnowledgeItemDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.knowledgeID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateKnowledgeItemDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
