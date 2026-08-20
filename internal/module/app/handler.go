package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service        *Service
	userMiddleware api.MiddlewareFunc
	appMiddleware  api.MiddlewareFunc
}

func NewHandler(service *Service, userMiddleware api.MiddlewareFunc, appMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, userMiddleware: userMiddleware, appMiddleware: appMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	// Public/connected-app surface — authenticated with the app's own key
	// (see guard.RequireAppKey), not tenant-staff bearer auth.
	r.Route("/apps/{code}/ping", func(r chi.Router) {
		r.Use(h.appMiddleware)
		r.Post("/", serve(h.ping))
	})

	// Admin/CRUD surface — tenant-staff bearer auth + membership, mirroring
	// member/tenant.
	r.Route("/tenants/{tenant_id}/apps", func(r chi.Router) {
		r.Use(h.userMiddleware)
		r.Get("/", serve(h.filter))
		r.Post("/", serve(h.create))
		r.Get("/{code}", serve(h.findByCode))
		r.Put("/{code}", serve(h.update))
		r.Post("/{code}/roll", serve(h.roll))
		r.Put("/{code}/active", serve(h.toggleActive))
		r.Put("/{code}/disabled", serve(h.toggleDisabled))
		r.Delete("/{code}", serve(h.delete))
	})
}

func (h *Handler) tenantID(r *http.Request) string {
	return chi.URLParam(r, "tenant_id")
}

// @Summary		List apps
// @Description	Lists apps, with optional pagination. Caller must be a tenant admin.
// @Tags			app
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	AppsPage
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/apps [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindAppsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}

	apps, err := h.service.Filter(r.Context(), h.tenantID(r), query.ToFilter())
	if err != nil {
		return err
	}

	return api.WriteJSON(w, http.StatusOK, apps)
}

// @Summary		Get app
// @Description	Returns an app by code. Caller must be a tenant admin.
// @Tags			app
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			code		path		string	true	"App code"
// @Success		200			{object}	AppResponse
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/apps/{code} [get]
func (h *Handler) findByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByCodeInTenant(r.Context(), h.tenantID(r), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create app
// @Description	Creates a new app. The plaintext key is only returned once, on creation. Caller must be a tenant admin.
// @Tags			app
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string			true	"Tenant ID"
// @Param			body		body		CreateAppDTO	true	"App data"
// @Success		201			{object}	AppWithSecret
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/apps [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateAppDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.tenantID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update app
// @Description	Updates an app's name/description by code. Caller must be a tenant admin.
// @Tags			app
// @Accept			json
// @Security		BearerAuth
// @Param			tenant_id	path	string			true	"Tenant ID"
// @Param			code		path	string			true	"App code"
// @Param			body		body	UpdateAppDTO	true	"App data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/apps/{code} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateAppDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.UpdateByCode(r.Context(), h.tenantID(r), chi.URLParam(r, "code"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Roll app key
// @Description	Rotates an app's key, invalidating the previous secret. The new plaintext key is only returned once. Caller must be a tenant admin.
// @Tags			app
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			code		path		string	true	"App code"
// @Success		200			{object}	AppWithSecret
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/apps/{code}/roll [post]
func (h *Handler) roll(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Roll(r.Context(), h.tenantID(r), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Activate app
// @Description	Marks an app as active. Caller must be a tenant admin.
// @Tags			app
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			code		path	string	true	"App code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/apps/{code}/active [put]
func (h *Handler) toggleActive(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), h.tenantID(r), chi.URLParam(r, "code"), true); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Disable app
// @Description	Marks an app as disabled. Caller must be a tenant admin.
// @Tags			app
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			code		path	string	true	"App code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/apps/{code}/disabled [put]
func (h *Handler) toggleDisabled(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), h.tenantID(r), chi.URLParam(r, "code"), false); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete app
// @Description	Deletes an app by code. Caller must be a tenant admin.
// @Tags			app
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			code		path	string	true	"App code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/apps/{code} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByCode(r.Context(), h.tenantID(r), chi.URLParam(r, "code")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Ping
// @Description	Called by a connected app, authenticated with its own key, to confirm it's alive. Promotes a pending app to active on first success.
// @Tags			app
// @Security		AppKeyAuth
// @Param			code	path	string	true	"App code"
// @Success		204
// @Failure		401	{object}	api.Error
// @Failure		403	{object}	api.Error
// @Router			/apps/{code}/ping [post]
func (h *Handler) ping(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Ping(r.Context(), chi.URLParam(r, "code")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
