package llm

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
	return &Handler{
		service:        service,
		userMiddleware: userMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/tenants/{tenant_id}/llm", func(r chi.Router) {
		r.Use(h.userMiddleware)
		r.Get("/", serve(h.filter))
		r.Get("/drivers", serve(h.listDrivers))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
	})
}

func (h *Handler) llmID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) tenantID(r *http.Request) string {
	return chi.URLParam(r, "tenant_id")
}

// @Summary		List LLMs
// @Description	Lists configured LLMs, with optional pagination. Caller must be a member of the tenant.
// @Tags			llm
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	LLMsPage
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/llm [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindLLMsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.tenantID(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get LLM
// @Description	Returns an LLM by ID. Caller must be a member of the tenant.
// @Tags			llm
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			id			path		string	true	"LLM ID"
// @Success		200			{object}	LLMResponse
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/llm/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.tenantID(r), h.llmID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create LLM
// @Description	Creates a new LLM configuration. Caller must be a member of the tenant.
// @Tags			llm
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string			true	"Tenant ID"
// @Param			body		body		CreateLLMDTO	true	"LLM data"
// @Success		201			{object}	LLMResponse
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/llm [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateLLMDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.tenantID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update LLM
// @Description	Updates an LLM by ID. Caller must be a member of the tenant.
// @Tags			llm
// @Accept			json
// @Security		BearerAuth
// @Param			tenant_id	path	string			true	"Tenant ID"
// @Param			id			path	string			true	"LLM ID"
// @Param			body		body	UpdateLLMDTO	true	"LLM data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/llm/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateLLMDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.Update(r.Context(), h.tenantID(r), h.llmID(r), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete LLM
// @Description	Deletes an LLM by ID. Caller must be a member of the tenant.
// @Tags			llm
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"LLM ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/llm/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.tenantID(r), h.llmID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		List LLM drivers
// @Description	Lists the available LLM provider drivers. Caller must be a member of the tenant.
// @Tags			llm
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Success		200			{array}		DriverInfo
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/llm/drivers [get]
func (h *Handler) listDrivers(w http.ResponseWriter, r *http.Request) error {
	drivers, err := h.service.ListDrivers(r.Context(), h.tenantID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, drivers)
}
