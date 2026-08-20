package session

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/module/agent/subscriber"
)

type Handler struct {
	service           *Service
	anyAuthMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, anyAuthMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, anyAuthMiddleware: anyAuthMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/apps/{code}/session", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.anyAuthMiddleware)
			r.Get("/", serve(h.filter))
			r.Post("/", serve(h.create))
			r.Get("/{id}", serve(h.findByID))
			r.Put("/{id}", serve(h.updateByID))
			r.Delete("/{id}", serve(h.deleteByID))
			r.Get("/{id}/messages", serve(h.findMessages))
			r.Post("/{id}/run", serve(h.run))
		})
	})
}

func (h *Handler) appCode(r *http.Request) string {
	return chi.URLParam(r, "code")
}

func (h *Handler) sessionID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

// @Summary		Run session
// @Description	Runs an agent within a session, streaming the execution via SSE.
// @Tags			session
// @Accept			json
// @Produce		text/event-stream
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path	string			true	"App code"
// @Param			id		path	string			true	"Session ID"
// @Param			body	body	RunSessionDTO	true	"Run input"
// @Success		200 {object}  any
// @Failure		400	{object}	api.Error
// @Router			/apps/{code}/session/{id}/run [post]
func (h *Handler) run(w http.ResponseWriter, r *http.Request) error {
	var dto RunSessionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}

	sse := subscriber.NewSSE(w)

	err := h.service.Run(r.Context(), h.appCode(r), h.sessionID(r), dto, sse)
	if err != nil {
		return sse.HandleError(err)
	}
	return nil
}

// @Summary		List session messages
// @Description	Lists the execution messages of a session, with optional ordering/pagination.
// @Tags			session
// @Produce		json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path		string	true	"App code"
// @Param			id		path		string	true	"Session ID"
// @Param			sort	query		string	false	"Order by timestamp"	Enums(asc, desc)
// @Param			take	query		int		false	"Page size"
// @Param			skip	query		int		false	"Page offset"
// @Success		200		{object}	SessionMessagesPage
// @Failure		400		{object}	api.Error
// @Router			/apps/{code}/session/{id}/messages [get]
func (h *Handler) findMessages(w http.ResponseWriter, r *http.Request) error {
	var query MessagesFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FindMessages(r.Context(), h.appCode(r), h.sessionID(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		List sessions
// @Description	Lists sessions of an app, with optional pagination and includes.
// @Tags			session
// @Produce		json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path		string		true	"App code"
// @Param			take	query		int			false	"Page size"
// @Param			skip	query		int			false	"Page offset"
// @Param			include	query		[]string	false	"Related resources to include"	Enums(agent)
// @Success		200		{object}	SessionsPage
// @Failure		400		{object}	api.Error
// @Router			/apps/{code}/session [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query SessionsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), h.appCode(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get session
// @Description	Returns a session by ID.
// @Tags			session
// @Produce		json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path		string		true	"App code"
// @Param			id		path		string		true	"Session ID"
// @Param			include	query		[]string	false	"Related resources to include"	Enums(agent)
// @Success		200		{object}	SessionResponse
// @Failure		404		{object}	api.Error
// @Router			/apps/{code}/session/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	var query SessionIncludeDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.FindByID(r.Context(), h.appCode(r), h.sessionID(r), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create session
// @Description	Creates a new session for an app.
// @Tags			session
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path		string				true	"App code"
// @Param			body	body		CreateSessionDTO	true	"Session data"
// @Success		201		{object}	SessionResponse
// @Failure		400		{object}	api.Error
// @Router			/apps/{code}/session [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateSessionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), h.appCode(r), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update session
// @Description	Updates a session by ID.
// @Tags			session
// @Accept			json
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path	string				true	"App code"
// @Param			id		path	string				true	"Session ID"
// @Param			body	body	UpdateSessionDTO	true	"Session data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/apps/{code}/session/{id} [put]
func (h *Handler) updateByID(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateSessionDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.UpdateByID(r.Context(), h.appCode(r), h.sessionID(r), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete session
// @Description	Deletes a session by ID.
// @Tags			session
// @Security		BearerAuth
// @Security		ApiKeyAuth
// @Param			code	path	string	true	"App code"
// @Param			id		path	string	true	"Session ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/apps/{code}/session/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.appCode(r), h.sessionID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
