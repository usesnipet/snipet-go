package agent

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/module/agent/subscriber"
	_ "github.com/usesnipet/snipet/internal/runtime/execution"
)

type Handler struct {
	service          *Service
	userMiddleware   api.MiddlewareFunc
	apiKeyMiddleware api.MiddlewareFunc
}

func NewHandler(
	service *Service,
	userMiddleware api.MiddlewareFunc,
	apiKeyMiddleware api.MiddlewareFunc,
) api.Handler {
	return &Handler{
		service:          service,
		userMiddleware:   userMiddleware,
		apiKeyMiddleware: apiKeyMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/tenants/{tenant_id}/agents", func(r chi.Router) {
		r.Use(h.userMiddleware)
		r.Get("/", serve(h.filter))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
	})

	// /agent/{id}/run is not restructured under /tenants/{tenant_id}/... —
	// it's the one runtime surface API keys still call directly (playground
	// runs, external integrations), so it stays on apiKeyMiddleware with
	// tenant derived server-side from the calling key (see Service.Run).
	r.Route("/agent", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Post("/{id}/run", serve(h.run))
		})
	})
}

func (h *Handler) agentID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func (h *Handler) tenantID(r *http.Request) string {
	return chi.URLParam(r, "tenant_id")
}

// @Summary		List agents
// @Description	Lists agents. Caller must be a member of the tenant.
// @Tags			agent
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Success		200			{object}	AgentsPage
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/agents [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Filter(r.Context(), h.tenantID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get agent
// @Description	Returns an agent by ID. Caller must be a member of the tenant.
// @Tags			agent
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string	true	"Tenant ID"
// @Param			id			path		string	true	"Agent ID"
// @Success		200			{object}	AgentResponse
// @Failure		404			{object}	api.Error
// @Router			/tenants/{tenant_id}/agents/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.tenantID(r), h.agentID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create agent
// @Description	Creates a new agent. Caller must be a member of the tenant.
// @Tags			agent
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path		string			true	"Tenant ID"
// @Param			body		body		CreateAgentDTO	true	"Agent data"
// @Success		201			{object}	AgentResponse
// @Failure		400			{object}	api.Error
// @Router			/tenants/{tenant_id}/agents [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.tenantID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update agent
// @Description	Updates an agent by ID. Caller must be a member of the tenant.
// @Tags			agent
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			tenant_id	path	string			true	"Tenant ID"
// @Param			id			path	string			true	"Agent ID"
// @Param			body		body	UpdateAgentDTO	true	"Agent data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/tenants/{tenant_id}/agents/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), h.tenantID(r), h.agentID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete agent
// @Description	Deletes an agent by ID. Caller must be a member of the tenant.
// @Tags			agent
// @Security		BearerAuth
// @Param			tenant_id	path	string	true	"Tenant ID"
// @Param			id			path	string	true	"Agent ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/tenants/{tenant_id}/agents/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.tenantID(r), h.agentID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Run agent
// @Description	Runs an agent in playground mode, streaming the execution via SSE.
// @Tags			agent
// @Accept			json
// @Produce		text/event-stream
// @Security		ApiKeyAuth
// @Param			id		path	string		true	"Agent ID"
// @Param			body	body	RunAgentDTO	true	"Run input"
// @Success		200 {object}  any
// @Failure		400	{object}	api.Error
// @Router			/agent/{id}/run [post]
func (h *Handler) run(w http.ResponseWriter, r *http.Request) error {
	var dto RunAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}

	sse := subscriber.NewSSE(w)

	err := h.service.Run(
		r.Context(),
		RunInput{AgentID: h.agentID(r), Message: dto.Message, StreamMessages: dto.StreamMessages},
		sse,
	)
	if err != nil {
		return sse.HandleError(err)
	}
	return nil
}
