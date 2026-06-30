package client

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

type Handler struct {
	service           *Service
	apiKeyMiddleware  api.MiddlewareFunc
	anyAuthMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, apiKeyMiddleware api.MiddlewareFunc, anyAuthMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, apiKeyMiddleware: apiKeyMiddleware, anyAuthMiddleware: anyAuthMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.anyAuthMiddleware)
			r.Get("/", serve(h.filterBy))
		})

		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Post("/", serve(h.create))
			r.Get("/{code}", serve(h.findByCode))
			r.Put("/{code}", serve(h.update))
			r.Delete("/{code}", serve(h.delete))
		})
	})
}

func (h *Handler) filterBy(w http.ResponseWriter, r *http.Request) error {
	var query FindClientsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}

	principal := auth.GetPrincipal(r.Context())
	var clients *page.Paginated[model.Client]
	var err error
	if principal.GetType() == auth.PrincipalTypeAPIKey {
		clients, err = h.service.FilterBy(r.Context(), query.ToFilter())
	} else if principal.GetType() == auth.PrincipalTypeJWT {
		clients, err = h.service.FilterByUser(r.Context(), principal.GetJWTClaims().Subject, query.ToFilter())
	}

	if err != nil {
		return err
	}

	return api.WriteJSON(w, http.StatusOK, clients)
}

func (h *Handler) findByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateClientDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateClientDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.UpdateByCode(r.Context(), chi.URLParam(r, "code"), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByCode(r.Context(), chi.URLParam(r, "code")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
