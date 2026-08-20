package appuser

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service           *Service
	appKeyMiddleware  api.MiddlewareFunc
	anyAuthMiddleware api.MiddlewareFunc
	jwtMiddleware     api.MiddlewareFunc
}

func NewHandler(
	service *Service,
	appKeyMiddleware api.MiddlewareFunc,
	jwtMiddleware api.MiddlewareFunc,
) api.Handler {
	return &Handler{
		service:          service,
		appKeyMiddleware: appKeyMiddleware,
		jwtMiddleware:    jwtMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/apps/{code}/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.appKeyMiddleware)
			r.Get("/", serve(h.filterBy))
			r.Post("/", serve(h.create))
		})

		r.Group(func(r chi.Router) {
			r.Use(h.jwtMiddleware)
			r.Get("/me", serve(h.me))
		})
	})
}

func (h *Handler) appCode(r *http.Request) string {
	return chi.URLParam(r, "code")
}

// @Summary		Get current user
// @Description	Returns the authenticated user for the given app.
// @Tags			user
// @Produce		json
// @Security	BearerAuth
// @Param			code	path		string	true	"App code"
// @Success		200		{object}	AppUserResponse
// @Failure		401		{object}	api.Error
// @Router			/apps/{code}/user/me [get]
func (h *Handler) me(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Me(r.Context(), h.appCode(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		List users
// @Description	Lists users of an app, with optional ordering/pagination.
// @Tags			user
// @Produce		json
// @Security	AppKeyAuth
// @Param			code		path		string	true	"App code"
// @Param			name_order	query		string	false	"Order by name"	Enums(asc, desc)
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	AppUsersPage
// @Failure		400			{object}	api.Error
// @Router			/apps/{code}/user [get]
func (h *Handler) filterBy(w http.ResponseWriter, r *http.Request) error {
	var query FindAppUsersFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FilterInApp(r.Context(), h.appCode(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create user
// @Description	Creates an app user via the app's own key (server-to-server).
// @Tags			user
// @Accept			json
// @Produce		json
// @Security		AppKeyAuth
// @Param			code	path	string				  	true	"App code"
// @Param			body	body	CreateAppUserDTO	true	"User data"
// @Success		201
// @Failure		400	{object}	api.Error
// @Router			/apps/{code}/user/authenticated [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var body CreateAppUserDTO
	if err := api.ParseBody(r, &body); err != nil {
		return err
	}
	err := h.service.Create(r.Context(), h.appCode(r), body)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, nil)
}
