package knowledge

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
	r.Route("/knowledge", func(r chi.Router) {
		r.Use(h.apiKeyMiddleware)
		r.Get("/", serve(h.filter))
		r.Get("/{id}/items", serve(h.filterItems))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
		r.Post("/test-connection", serve(h.testConnection))
		r.Post("/{id}/sync", serve(h.sync))
	})
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var dto FilterKnowledgeDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), dto.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) filterItems(w http.ResponseWriter, r *http.Request) error {
	var dto FilterKnowledgeItemDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.FilterItems(r.Context(), chi.URLParam(r, "id"), dto.ToFilter())
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
	var dto CreateKnowledgeDTO
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
	var dto UpdateKnowledgeDTO
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

func (h *Handler) testConnection(w http.ResponseWriter, r *http.Request) error {
	var dto TestConnectionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.TestConnection(r.Context(), dto.Driver, dto.Configuration)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) sync(w http.ResponseWriter, r *http.Request) error {
	var dto SyncKnowledgeQueryDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	knowledgeID := chi.URLParam(r, "id")
	err := h.service.Sync(r.Context(), knowledgeID, dto.Force)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
