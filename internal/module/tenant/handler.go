package tenant

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

// Handler routes tenant CRUD under /tenants. Platform-admin routes
// (GET /tenants list, RequirePlatformAdmin) are intentionally not wired up
// yet — every route here just requires a valid bearer token, with
// finer-grained membership/admin checks done in the service.
type Handler struct {
	service               *Service
	platformJWTMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, platformJWTMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, platformJWTMiddleware: platformJWTMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/tenants", func(r chi.Router) {
		r.Use(h.platformJWTMiddleware)

		r.Get("/me", serve(h.findMine))
		r.Post("/", serve(h.create))
		r.Get("/slug/{slug}", serve(h.findBySlug))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
		r.Post("/{id}/leave", serve(h.leave))
	})
}

// @Summary		My tenants
// @Description	Returns the tenants the authenticated user is a member of.
// @Tags			tenant
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	TenantsPage
// @Failure		401	{object}	api.Error
// @Router			/tenants/me [get]
func (h *Handler) findMine(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindMine(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get tenant
// @Description	Returns a tenant. Caller must be a member.
// @Tags			tenant
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"Tenant ID"
// @Success		200	{object}	TenantResponse
// @Failure		403	{object}	api.Error
// @Failure		404	{object}	api.Error
// @Router			/tenants/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	data, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get tenant by slug
// @Description	Returns a tenant by its slug. Caller must be a member.
// @Tags			tenant
// @Produce		json
// @Security		BearerAuth
// @Param			slug	path		string	true	"Tenant slug"
// @Success		200	{object}	TenantResponse
// @Failure		403	{object}	api.Error
// @Failure		404	{object}	api.Error
// @Router			/tenants/slug/{slug} [get]
func (h *Handler) findBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := chi.URLParam(r, "slug")
	data, err := h.service.FindBySlug(r.Context(), slug)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create tenant
// @Description	Creates a new tenant. Blocked without a valid Multi-Tenant Use license once a tenant already exists. The creator becomes the tenant's admin.
// @Tags			tenant
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			body	body		CreateTenantDTO	true	"Tenant data"
// @Success		201		{object}	TenantResponse
// @Failure		403		{object}	api.Error
// @Failure		409		{object}	api.Error
// @Router			/tenants [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateTenantDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update tenant
// @Description	Updates a tenant. Caller must be able to manage it (tenant admin or platform admin).
// @Tags			tenant
// @Accept			json
// @Security		BearerAuth
// @Param			id		path	string			true	"Tenant ID"
// @Param			body	body	UpdateTenantDTO	true	"Tenant data"
// @Success		204
// @Failure		403	{object}	api.Error
// @Router			/tenants/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var dto UpdateTenantDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.UpdateByID(r.Context(), id, dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete tenant
// @Description	Deletes a tenant. Caller must be able to manage it (tenant admin or platform admin).
// @Tags			tenant
// @Security		BearerAuth
// @Param			id	path	string	true	"Tenant ID"
// @Success		204
// @Failure		403	{object}	api.Error
// @Router			/tenants/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteByID(r.Context(), id); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Leave tenant
// @Description	Removes the caller's own membership. Blocked if the caller is the tenant's last active admin.
// @Tags			tenant
// @Security		BearerAuth
// @Param			id	path	string	true	"Tenant ID"
// @Success		204
// @Failure		409	{object}	api.Error
// @Router			/tenants/{id}/leave [post]
func (h *Handler) leave(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.service.Leave(r.Context(), id); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
