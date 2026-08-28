package knowledge

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service       *Service
	basicAuthGate api.Gate
}

func NewHandler(service *Service, basicAuthGate api.Gate) api.Handler {
	return &Handler{service: service, basicAuthGate: basicAuthGate}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/knowledge", func(r chi.Router) {
		r.Use(h.basicAuthGate.Handler())
		r.Get("/", serve(h.filter))
		r.Get("/drivers", serve(h.listDrivers))
		r.Get("/{id}/items", serve(h.filterItems))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
		r.Post("/test-connection", serve(h.testConnection))
		r.Post("/{id}/sync", serve(h.sync))
	})
}

// @Summary		List knowledge bases
// @Description	Lists knowledge bases, with optional pagination.
// @Tags			knowledge
// @Produce		json
// @Security		BasicAuth
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	KnowledgesPage
// @Failure		400			{object}	api.Error
// @Router			/knowledge [get]
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

// @Summary		List knowledge items
// @Description	Lists items of a knowledge base, with optional pagination.
// @Tags			knowledge
// @Produce		json
// @Security		BasicAuth
// @Param			id			path		string	true	"Knowledge ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	KnowledgeItemsPage
// @Failure		400			{object}	api.Error
// @Router			/knowledge/{id}/items [get]
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

// @Summary		Get knowledge base
// @Description	Returns a knowledge base by ID.
// @Tags			knowledge
// @Produce		json
// @Security		BasicAuth
// @Param			id			path		string	true	"Knowledge ID"
// @Success		200			{object}	KnowledgeResponse
// @Failure		404			{object}	api.Error
// @Router			/knowledge/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create knowledge base
// @Description	Creates a new knowledge base.
// @Tags			knowledge
// @Accept			json
// @Produce		json
// @Security		BasicAuth
// @Param			body		body		CreateKnowledgeDTO	true	"Knowledge data"
// @Success		201			{object}	CreateKnowledgeResponseDTO
// @Failure		400			{object}	api.Error
// @Router			/knowledge [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateKnowledgeDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, map[string]any{
		"knowledge": data,
	})
}

// @Summary		Update knowledge base
// @Description	Updates a knowledge base by ID.
// @Tags			knowledge
// @Accept			json
// @Security		BasicAuth
// @Param			id			path	string				true	"Knowledge ID"
// @Param			body		body	UpdateKnowledgeDTO	true	"Knowledge data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/knowledge/{id} [put]
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

// @Summary		Delete knowledge base
// @Description	Deletes a knowledge base by ID.
// @Tags			knowledge
// @Security		BasicAuth
// @Param			id			path	string	true	"Knowledge ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/knowledge/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), chi.URLParam(r, "id")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Test knowledge connection
// @Description	Tests a knowledge source driver configuration without persisting it.
// @Tags			knowledge
// @Accept			json
// @Security		BasicAuth
// @Param			body		body	TestConnectionDTO	true	"Connection data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/knowledge/test-connection [post]
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

// @Summary		List knowledge source drivers
// @Description	Lists the available knowledge source drivers.
// @Tags			knowledge
// @Produce		json
// @Security		BasicAuth
// @Success		200			{object}	DriversDTO
// @Failure		400			{object}	api.Error
// @Router			/knowledge/drivers [get]
func (h *Handler) listDrivers(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.ListDrivers(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Sync knowledge base
// @Description	Triggers a sync of a knowledge base's items from its source.
// @Tags			knowledge
// @Security		BasicAuth
// @Param			id			path	string	true	"Knowledge ID"
// @Param			force		query	bool	false	"Force a full resync"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/knowledge/{id}/sync [post]
func (h *Handler) sync(w http.ResponseWriter, r *http.Request) error {
	var dto SyncKnowledgeQueryDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	knowledgeID := chi.URLParam(r, "id")
	if err := h.service.Sync(r.Context(), knowledgeID, dto.Force); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
