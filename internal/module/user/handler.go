package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service           *Service
	apiKeyMiddleware  api.MiddlewareFunc
	anyAuthMiddleware api.MiddlewareFunc
	jwtMiddleware     api.MiddlewareFunc
}

func NewHandler(
	service *Service,
	apiKeyMiddleware api.MiddlewareFunc,
	anyAuthMiddleware api.MiddlewareFunc,
	jwtMiddleware api.MiddlewareFunc,
) api.Handler {
	return &Handler{
		service:           service,
		apiKeyMiddleware:  apiKeyMiddleware,
		anyAuthMiddleware: anyAuthMiddleware,
		jwtMiddleware:     jwtMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_code}/user", func(r chi.Router) {
		r.Post("/anonymous", serve(h.createAnonymous))

		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Post("/authenticated", serve(h.createAuthenticated))
		})

		r.Group(func(r chi.Router) {
			r.Use(h.jwtMiddleware)
			r.Get("/me", serve(h.me))
		})

		r.Group(func(r chi.Router) {
			r.Use(h.anyAuthMiddleware)
			r.Get("/", serve(h.filterBy))
		})
	})
}

func (h *Handler) clientCode(r *http.Request) string {
	return chi.URLParam(r, "client_code")
}

// @Summary		Get current user
// @Description	Returns the authenticated user for the given client.
// @Tags			user
// @Produce		json
// @Security		BearerAuth
// @Param			client_code	path		string	true	"Client code"
// @Success		200			{object}	UserResponse
// @Failure		401			{object}	api.Error
// @Router			/client/{client_code}/user/me [get]
func (h *Handler) me(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Me(r.Context(), h.clientCode(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		List users
// @Description	Lists users of a client, with optional ordering/pagination.
// @Tags			user
// @Produce		json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			client_code	path		string	true	"Client code"
// @Param			name_order	query		string	false	"Order by name"	Enums(asc, desc)
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	UsersPage
// @Failure		400			{object}	api.Error
// @Router			/client/{client_code}/user [get]
func (h *Handler) filterBy(w http.ResponseWriter, r *http.Request) error {
	var query FindUsersFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FilterInClient(r.Context(), h.clientCode(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create anonymous user
// @Description	Creates an anonymous user for a client.
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			client_code	path	string							true	"Client code"
// @Param			body		body	CreateAnonymousClientUserDTO	true	"Anonymous user data"
// @Success		201
// @Failure		400	{object}	api.Error
// @Router			/client/{client_code}/user/anonymous [post]
func (h *Handler) createAnonymous(w http.ResponseWriter, r *http.Request) error {
	var body CreateAnonymousClientUserDTO
	if err := api.ParseBody(r, &body); err != nil {
		return err
	}
	err := h.service.CreateAnonymous(r.Context(), h.clientCode(r), body)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, nil)
}

// @Summary		Create authenticated user
// @Description	Creates a client user via API key (server-to-server).
// @Tags			user
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			client_code	path	string								true	"Client code"
// @Param			body		body	CreateAuthenticatedClientUserDTO	true	"User data"
// @Success		201
// @Failure		400	{object}	api.Error
// @Router			/client/{client_code}/user/authenticated [post]
func (h *Handler) createAuthenticated(w http.ResponseWriter, r *http.Request) error {
	var body CreateAuthenticatedClientUserDTO
	if err := api.ParseBody(r, &body); err != nil {
		return err
	}
	err := h.service.CreateAuthenticated(r.Context(), h.clientCode(r), body)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, nil)
}
