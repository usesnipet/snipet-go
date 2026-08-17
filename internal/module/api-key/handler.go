package apikey

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service          *Service
	userMiddleware   api.MiddlewareFunc
	apiKeyMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, userMiddleware, apiKeyMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, userMiddleware: userMiddleware, apiKeyMiddleware: apiKeyMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/tenants/{tenant_id}/api-keys", func(r chi.Router) {
		r.Use(h.userMiddleware)
		r.Get("/", serve(h.filter))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Post("/{id}/roll", serve(h.roll))
		r.Put("/{id}/expiration", serve(h.updateExpiration))
		r.Put("/{id}/active", serve(h.toggleActive))
		r.Put("/{id}/disabled", serve(h.toggleDisabled))
		r.Delete("/{id}", serve(h.delete))
	})

	// /api-key/me only makes sense called WITH an API key — it introspects
	// whichever key authenticated the request — so it stays on
	// apiKeyMiddleware instead of moving under /tenants/{tenant_id}/api-keys.
	r.Route("/api-key/me", func(r chi.Router) {
		r.Use(h.apiKeyMiddleware)
		r.Get("/", serve(h.me))
	})
}

func (h *Handler) tenantID(r *http.Request) string {
	return chi.URLParam(r, "tenant_id")
}

// @Summary		List API keys
// @Description	Lists API keys, with optional pagination. Caller must be a tenant admin.
// @Tags			api-key
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	APIKeysPage
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindAPIKeysFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.tenantID(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get current API key
// @Description	Returns the API key used to authenticate the current request.
// @Tags			api-key
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{object}	APIKeyResponse
// @Failure		401	{object}	api.Error
// @Router			/api-key/me [get]
func (h *Handler) me(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Me(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get API key
// @Description	Returns an API key by ID. Caller must be a tenant admin.
// @Tags			api-key
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			id			path		string	true	"API key ID"
// @Success		200			{object}	APIKeyResponse
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.tenantID(r), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create API key
// @Description	Creates a new API key. The plaintext key is only returned once, on creation. Caller must be a tenant admin.
// @Tags			api-key
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string			true	"Tenant ID"
// @Param			body		body		CreateAPIKeyDTO	true	"API key data"
// @Success		201			{object}	APIKeyWithSecret
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateAPIKeyDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.tenantID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Roll API key
// @Description	Rotates an API key, invalidating the previous secret. The new plaintext key is only returned once. Caller must be a tenant admin.
// @Tags			api-key
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			id			path		string	true	"API key ID"
// @Success		200			{object}	APIKeyWithSecret
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys/{id}/roll [post]
func (h *Handler) roll(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Roll(r.Context(), h.tenantID(r), chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Update API key expiration
// @Description	Updates the expiration date of an API key. Caller must be a tenant admin.
// @Tags			api-key
// @Accept			json
// @Security		BearerAuth
// @Param			tenant_id	path	string				true	"Tenant ID"
// @Param			id			path	string				true	"API key ID"
// @Param			body		body	UpdateExpirationDTO	true	"Expiration data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys/{id}/expiration [put]
func (h *Handler) updateExpiration(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateExpirationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.UpdateExpiration(r.Context(), h.tenantID(r), chi.URLParam(r, "id"), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Activate API key
// @Description	Marks an API key as active. Caller must be a tenant admin.
// @Tags			api-key
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"API key ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys/{id}/active [put]
func (h *Handler) toggleActive(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), h.tenantID(r), chi.URLParam(r, "id"), true); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Disable API key
// @Description	Marks an API key as disabled. Caller must be a tenant admin.
// @Tags			api-key
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"API key ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys/{id}/disabled [put]
func (h *Handler) toggleDisabled(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), h.tenantID(r), chi.URLParam(r, "id"), false); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete API key
// @Description	Deletes an API key by ID. Caller must be a tenant admin.
// @Tags			api-key
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"API key ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/api-keys/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Delete(r.Context(), h.tenantID(r), chi.URLParam(r, "id")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
