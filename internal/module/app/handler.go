package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service       *Service
	basicAuthGate api.Gate
	appKeyGate    api.Gate
}

func NewHandler(service *Service, basicAuthGate api.Gate, appKeyGate api.Gate) api.Handler {
	return &Handler{service: service, basicAuthGate: basicAuthGate, appKeyGate: appKeyGate}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	// A single mount for the whole surface. Splitting this into sibling
	// r.Route("/apps/{code}", ...) and r.Route("/apps", ...) blocks makes chi
	// mount two subrouters at overlapping paths, and the more specific
	// "/apps/{code}" mount then shadows every "/apps/{code}/*" admin route.
	r.Route("/apps", func(r chi.Router) {
		// Public surface — no auth at all. findPublicByCode only ever exposes
		// apps explicitly marked Public.
		r.Get("/{code}/public", serve(h.findPublicByCode))

		// ping is authenticated with the app's own key (see
		// guard.RequireAppKey), not admin basic auth.
		r.Group(func(r chi.Router) {
			r.Use(h.appKeyGate.Handler())
			r.Post("/{code}/ping", serve(h.ping))
		})

		// Admin/CRUD surface.
		r.Group(func(r chi.Router) {
			r.Use(h.basicAuthGate.Handler())
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{code}", serve(h.findByCode))
			r.Put("/{code}", serve(h.update))
			r.Put("/{code}/agents", serve(h.linkAgents))
			r.Put("/{code}/auth-config", serve(h.updateAuthConfig))
			r.Post("/{code}/roll", serve(h.roll))
			r.Put("/{code}/active", serve(h.toggleActive))
			r.Put("/{code}/disabled", serve(h.toggleDisabled))
			r.Delete("/{code}", serve(h.delete))
		})
	})
}

// @Summary		List apps
// @Description	Lists apps, with optional pagination.
// @Tags			app
// @Produce		json
// @Security		BasicAuth
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	AppsPage
// @Failure		400			{object}	api.Error
// @Router			/apps [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindAppsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}

	apps, err := h.service.Filter(r.Context(), query.ToFilter())
	if err != nil {
		return err
	}

	return api.WriteJSON(w, http.StatusOK, apps)
}

// @Summary		Get app
// @Description	Returns an app by code.
// @Tags			app
// @Produce		json
// @Security		BasicAuth
// @Param			code		path		string	true	"App code"
// @Success		200			{object}	AppResponse
// @Failure		404			{object}	api.Error
// @Router			/apps/{code} [get]
func (h *Handler) findByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create app
// @Description	Creates a new app. The plaintext key is only returned once, on creation.
// @Tags			app
// @Accept			json
// @Produce		json
// @Security		BasicAuth
// @Param			body		body		CreateAppDTO	true	"App data"
// @Success		201			{object}	AppWithSecret
// @Failure		400			{object}	api.Error
// @Router			/apps [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateAppDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update app
// @Description	Updates an app's name/description by code.
// @Tags			app
// @Accept			json
// @Security		BasicAuth
// @Param			code		path	string			true	"App code"
// @Param			body		body	UpdateAppDTO	true	"App data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/apps/{code} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateAppDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.UpdateByCode(r.Context(), chi.URLParam(r, "code"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Link agents to app
// @Description	Replaces the whole set of agents linked to an app.
// @Tags			app
// @Accept			json
// @Security		BasicAuth
// @Param			code		path	string				true	"App code"
// @Param			body		body	LinkAppAgentsDTO	true	"Agent ids"
// @Success		204
// @Failure		400	{object}	api.Error
// @Failure		404	{object}	api.Error
// @Router			/apps/{code}/agents [put]
func (h *Handler) linkAgents(w http.ResponseWriter, r *http.Request) error {
	var dto LinkAppAgentsDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.LinkAgents(r.Context(), chi.URLParam(r, "code"), dto.AgentIDs); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Update app auth config
// @Description	Replaces how the app federates its end-users' identity (OIDC/webhook).
// @Tags			app
// @Accept			json
// @Security		BasicAuth
// @Param			code		path	string					true	"App code"
// @Param			body		body	UpdateAppAuthConfigDTO	true	"Auth config data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/apps/{code}/auth-config [put]
func (h *Handler) updateAuthConfig(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateAppAuthConfigDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.UpdateAuthConfig(r.Context(), chi.URLParam(r, "code"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Roll app key
// @Description	Rotates an app's key, invalidating the previous secret. The new plaintext key is only returned once.
// @Tags			app
// @Produce		json
// @Security		BasicAuth
// @Param			code		path		string	true	"App code"
// @Success		200			{object}	AppWithSecret
// @Failure		404			{object}	api.Error
// @Router			/apps/{code}/roll [post]
func (h *Handler) roll(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Roll(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Activate app
// @Description	Marks an app as active.
// @Tags			app
// @Security		BasicAuth
// @Param			code		path	string	true	"App code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/apps/{code}/active [put]
func (h *Handler) toggleActive(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), chi.URLParam(r, "code"), true); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Disable app
// @Description	Marks an app as disabled.
// @Tags			app
// @Security		BasicAuth
// @Param			code		path	string	true	"App code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/apps/{code}/disabled [put]
func (h *Handler) toggleDisabled(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.ToggleActive(r.Context(), chi.URLParam(r, "code"), false); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete app
// @Description	Deletes an app by code.
// @Tags			app
// @Security		BasicAuth
// @Param			code		path	string	true	"App code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/apps/{code} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByCode(r.Context(), chi.URLParam(r, "code")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Get public app info
// @Description	Returns an app's name, code, and description. Only exposed when the app is marked public; otherwise 404. Unauthenticated.
// @Tags			app
// @Produce		json
// @Param			code	path		string	true	"App code"
// @Success		200		{object}	PublicAppDTO
// @Failure		404		{object}	api.Error
// @Router			/apps/{code}/public [get]
func (h *Handler) findPublicByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindPublicByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
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
