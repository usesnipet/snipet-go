package agent

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/module/agent/subscriber"
)

type Handler struct {
	service          *Service
	apiKeyMiddleware api.MiddlewareFunc
}

func NewHandler(
	service *Service,
	apiKeyMiddleware api.MiddlewareFunc,
) api.Handler {
	return &Handler{
		service:          service,
		apiKeyMiddleware: apiKeyMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/agent", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.apiKeyMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Put("/{id}", serve(h.update))
			r.Delete("/{id}", serve(h.deleteByID))
			r.Post("/{id}/run", serve(h.run))
		})
	})
}

func (h *Handler) agentID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

// @Summary		List agents
// @Description	Lists agents.
// @Tags			agent
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{object}	AgentsPage
// @Failure		400	{object}	api.Error
// @Router			/agent [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Filter(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get agent
// @Description	Returns an agent by ID.
// @Tags			agent
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id	path		string	true	"Agent ID"
// @Success		200	{object}	AgentResponse
// @Failure		404	{object}	api.Error
// @Router			/agent/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.agentID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create agent
// @Description	Creates a new agent.
// @Tags			agent
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			body	body		CreateAgentDTO	true	"Agent data"
// @Success		201		{object}	AgentResponse
// @Failure		400		{object}	api.Error
// @Router			/agent [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update agent
// @Description	Updates an agent by ID.
// @Tags			agent
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id		path	string			true	"Agent ID"
// @Param			body	body	UpdateAgentDTO	true	"Agent data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/agent/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	err := h.service.Update(r.Context(), h.agentID(r), dto)
	if err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete agent
// @Description	Deletes an agent by ID.
// @Tags			agent
// @Security		ApiKeyAuth
// @Param			id	path	string	true	"Agent ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/agent/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.agentID(r)); err != nil {
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
// @Success		200
// @Failure		400	{object}	api.Error
// @Router			/agent/{id}/run [post]
func (h *Handler) run(w http.ResponseWriter, r *http.Request) error {
	var dto RunAgentDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}

	sse := subscriber.NewSSE(w)

	err := h.service.Run(r.Context(), RunInput{AgentID: h.agentID(r), Message: dto.Message}, sse)
	if err != nil {
		return sse.HandleError(err)
	}
	return nil
}
