package knowledge

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service        *Service
	userMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, userMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, userMiddleware: userMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/tenants/{tenant_id}/knowledge", func(r chi.Router) {
		r.Use(h.userMiddleware)
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

func (h *Handler) tenantID(r *http.Request) string {
	return chi.URLParam(r, "tenant_id")
}

// @Summary		List knowledge bases
// @Description	Lists knowledge bases, with optional pagination. Caller must be a member of the tenant.
// @Tags			knowledge
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	KnowledgesPage
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var dto FilterKnowledgeDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.tenantID(r), dto.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		List knowledge items
// @Description	Lists items of a knowledge base, with optional pagination. Caller must be a member of the tenant.
// @Tags			knowledge
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			id			path		string	true	"Knowledge ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	KnowledgeItemsPage
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/{id}/items [get]
func (h *Handler) filterItems(w http.ResponseWriter, r *http.Request) error {
	var dto FilterKnowledgeItemDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	data, err := h.service.FilterItems(r.Context(), h.tenantID(r), chi.URLParam(r, "id"), dto.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get knowledge base
// @Description	Returns a knowledge base by ID. Caller must be a member of the tenant.
// @Tags			knowledge
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			id			path		string	true	"Knowledge ID"
// @Success		200			{object}	KnowledgeResponse
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.tenantID(r), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create knowledge base
// @Description	Creates a new knowledge base. Caller must be a member of the tenant.
// @Tags			knowledge
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string				true	"Tenant ID"
// @Param			body		body		CreateKnowledgeDTO	true	"Knowledge data"
// @Success		201			{object}	CreateKnowledgeResponseDTO
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateKnowledgeDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.tenantID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, map[string]any{
		"knowledge": data,
	})
}

// @Summary		Update knowledge base
// @Description	Updates a knowledge base by ID. Caller must be a member of the tenant.
// @Tags			knowledge
// @Accept			json
// @Security		BearerAuth
// @Param			tenant_id	path	string				true	"Tenant ID"
// @Param			id			path	string				true	"Knowledge ID"
// @Param			body		body	UpdateKnowledgeDTO	true	"Knowledge data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateKnowledgeDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), h.tenantID(r), chi.URLParam(r, "id"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete knowledge base
// @Description	Deletes a knowledge base by ID. Caller must be a member of the tenant.
// @Tags			knowledge
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"Knowledge ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.tenantID(r), chi.URLParam(r, "id")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Test knowledge connection
// @Description	Tests a knowledge source driver configuration without persisting it. Caller must be a member of the tenant.
// @Tags			knowledge
// @Accept			json
// @Security		BearerAuth
// @Param			tenant_id	path	string				true	"Tenant ID"
// @Param			body		body	TestConnectionDTO	true	"Connection data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/test-connection [post]
func (h *Handler) testConnection(w http.ResponseWriter, r *http.Request) error {
	var dto TestConnectionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.TestConnection(r.Context(), h.tenantID(r), dto.Driver, dto.Configuration)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		List knowledge source drivers
// @Description	Lists the available knowledge source drivers. Caller must be a member of the tenant.
// @Tags			knowledge
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Success		200			{object}	DriversDTO
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/drivers [get]
func (h *Handler) listDrivers(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.ListDrivers(r.Context(), h.tenantID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Sync knowledge base
// @Description	Triggers a sync of a knowledge base's items from its source. Caller must be a member of the tenant.
// @Tags			knowledge
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"Knowledge ID"
// @Param			force		query	bool	false	"Force a full resync"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/knowledge/{id}/sync [post]
func (h *Handler) sync(w http.ResponseWriter, r *http.Request) error {
	var dto SyncKnowledgeQueryDTO
	if err := api.ParseQuery(r, &dto); err != nil {
		return err
	}
	knowledgeID := chi.URLParam(r, "id")
	if err := h.service.Sync(r.Context(), h.tenantID(r), knowledgeID, dto.Force); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
