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

func (h *Handler) me(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Me(r.Context(), h.clientCode(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

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
