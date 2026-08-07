package client

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
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
	// Use /clients (plural) so CRUD does not conflict with /client/{client_code}/...
	// routes used by auth, session, and user modules.
	r.Route("/clients", func(r chi.Router) {
		r.Get("/{code}/public", serve(h.findPublicByCode))

		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{code}", serve(h.findByCode))
			r.Put("/{code}", serve(h.update))
			r.Delete("/{code}", serve(h.delete))
		})

		r.Group(func(r chi.Router) {
			r.Use(h.anyAuthMiddleware)

			r.Get("/{code}/agents", serve(h.getAgents))
		})
	})
}

// @Summary		List clients
// @Description	Lists clients, with optional pagination.
// @Tags			client
// @Produce		json
// @Security		ApiKeyAuth
// @Param			take	query		int	false	"Page size"
// @Param			skip	query		int	false	"Page offset"
// @Success		200		{object}	ClientsPage
// @Failure		400		{object}	api.Error
// @Router			/clients [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindClientsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}

	clients, err := h.service.Filter(r.Context(), query.ToFilter())
	if err != nil {
		return err
	}

	return api.WriteJSON(w, http.StatusOK, clients)
}

// @Summary		Get client
// @Description	Returns a client by code.
// @Tags			client
// @Produce		json
// @Security		ApiKeyAuth
// @Param			code	path		string	true	"Client code"
// @Success		200		{object}	ClientResponse
// @Failure		404		{object}	api.Error
// @Router			/clients/{code} [get]
func (h *Handler) findByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get public client info
// @Description	Returns public information for a client by code.
// @Tags			client
// @Produce		json
// @Param			code	path		string	true	"Client code"
// @Success		200		{object}	ClientPublicDTO
// @Failure		404		{object}	api.Error
// @Router			/clients/{code}/public [get]
func (h *Handler) findPublicByCode(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindPublicByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create client
// @Description	Creates a new client.
// @Tags			client
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			body	body		CreateClientDTO	true	"Client data"
// @Success		201		{object}	ClientResponse
// @Failure		400		{object}	api.Error
// @Router			/clients [post]
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

// @Summary		Update client
// @Description	Updates a client by code.
// @Tags			client
// @Accept			json
// @Security		ApiKeyAuth
// @Param			code	path	string			true	"Client code"
// @Param			body	body	UpdateClientDTO	true	"Client data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/clients/{code} [put]
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

// @Summary		Delete client
// @Description	Deletes a client by code.
// @Tags			client
// @Security		ApiKeyAuth
// @Param			code	path	string	true	"Client code"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/clients/{code} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByCode(r.Context(), chi.URLParam(r, "code")); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		List client agents
// @Description	Lists agents linked to a client, with optional pagination.
// @Tags			client
// @Produce		json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path		string	true	"Client code"
// @Param			take	query		int		false	"Page size"
// @Param			skip	query		int		false	"Page offset"
// @Success		200		{object}	ClientAgentsPage
// @Failure		400		{object}	api.Error
// @Router			/clients/{code}/agents [get]
func (h *Handler) getAgents(w http.ResponseWriter, r *http.Request) error {
	var query FindClientAgentsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	agents, err := h.service.GetAgents(r.Context(), chi.URLParam(r, "code"), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, agents)
}
