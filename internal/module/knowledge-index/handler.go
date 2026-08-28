package knowledgeindex

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service       *Service
	basicAuthGate api.Gate
}

func NewHandler(
	service *Service,
	basicAuthGate api.Gate,
) api.Handler {
	return &Handler{
		service:       service,
		basicAuthGate: basicAuthGate,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/knowledge/index", func(r chi.Router) {
		r.Use(h.basicAuthGate.Handler())
		r.Get("/drivers", serve(h.listDrivers))
	})
	r.Route("/knowledge/{knowledge_id}/index", func(r chi.Router) {
		r.Use(h.basicAuthGate.Handler())
		r.Get("/", serve(h.filter))
		r.Get("/{id}/items", serve(h.filterItems))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
	})
}

func (h *Handler) knowledgeID(r *http.Request) string {
	return chi.URLParam(r, "knowledge_id")
}

// @Summary		List knowledge indexes
// @Description	Lists indexes of a knowledge base, with optional pagination.
// @Tags			knowledge-index
// @Produce		json
// @Security		BasicAuth
// @Param			knowledge_id	path		string	true	"Knowledge ID"
// @Param			take			query		int		false	"Page size"
// @Param			skip			query		int		false	"Page offset"
// @Success		200				{object}	KnowledgeIndexesPage
// @Failure		400				{object}	api.Error
// @Router			/knowledge/{knowledge_id}/index [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var dto FilterKnowledgeIndexDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.knowledgeID(r), dto.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		List indexed knowledge items
// @Description	Lists items of a knowledge index, with optional pagination.
// @Tags			knowledge-index
// @Produce		json
// @Security		BasicAuth
// @Param			knowledge_id	path		string	true	"Knowledge ID"
// @Param			id				path		string	true	"Index ID"
// @Param			take			query		int		false	"Page size"
// @Param			skip			query		int		false	"Page offset"
// @Success		200				{object}	IndexedKnowledgeItemsPage
// @Failure		400				{object}	api.Error
// @Router			/knowledge/{knowledge_id}/index/{id}/items [get]
func (h *Handler) filterItems(w http.ResponseWriter, r *http.Request) error {
	var dto FilterIndexedKnowledgeItemDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.FilterItems(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id"), dto.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get knowledge index
// @Description	Returns a knowledge index by ID.
// @Tags			knowledge-index
// @Produce		json
// @Security		BasicAuth
// @Param			knowledge_id	path		string	true	"Knowledge ID"
// @Param			id				path		string	true	"Index ID"
// @Success		200				{object}	KnowledgeIndexResponse
// @Failure		404				{object}	api.Error
// @Router			/knowledge/{knowledge_id}/index/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create knowledge index
// @Description	Creates a new index for a knowledge base.
// @Tags			knowledge-index
// @Accept			json
// @Produce		json
// @Security		BasicAuth
// @Param			knowledge_id	path		string					true	"Knowledge ID"
// @Param			body			body		CreateKnowledgeIndexDTO	true	"Index data"
// @Success		201				{object}	KnowledgeIndexResponse
// @Failure		400				{object}	api.Error
// @Router			/knowledge/{knowledge_id}/index [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateKnowledgeIndexDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.knowledgeID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update knowledge index
// @Description	Updates a knowledge index by ID.
// @Tags			knowledge-index
// @Accept			json
// @Security		BasicAuth
// @Param			knowledge_id	path	string					true	"Knowledge ID"
// @Param			id				path	string					true	"Index ID"
// @Param			body			body	UpdateKnowledgeIndexDTO	true	"Index data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/knowledge/{knowledge_id}/index/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateKnowledgeIndexDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete knowledge index
// @Description	Deletes a knowledge index by ID.
// @Tags			knowledge-index
// @Security		BasicAuth
// @Param			knowledge_id	path	string	true	"Knowledge ID"
// @Param			id				path	string	true	"Index ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/knowledge/{knowledge_id}/index/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.knowledgeID(r), chi.URLParam(r, "id")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		List knowledge index drivers
// @Description	Lists the available knowledge index drivers.
// @Tags			knowledge-index
// @Produce		json
// @Security		BasicAuth
// @Success		200			{object}	DriversDTO
// @Failure		400			{object}	api.Error
// @Router			/knowledge/index/drivers [get]
func (h *Handler) listDrivers(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.ListDrivers(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}
