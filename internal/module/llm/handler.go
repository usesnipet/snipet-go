package llm

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service       *Service
	basicAuthGate api.Gate
}

func NewHandler(service *Service, basicAuthGate api.Gate) api.Handler {
	return &Handler{
		service:       service,
		basicAuthGate: basicAuthGate,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/llm", func(r chi.Router) {
		r.Use(h.basicAuthGate.Handler())
		r.Get("/", serve(h.filter))
		r.Get("/drivers", serve(h.listDrivers))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Put("/{id}", serve(h.update))
		r.Delete("/{id}", serve(h.deleteByID))
	})
}

func (h *Handler) llmID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

// @Summary		List LLMs
// @Description	Lists configured LLMs, with optional pagination.
// @Tags			llm
// @Produce		json
// @Security		BasicAuth
// @Param			take		query		int		false	"Page size"
// @Param			skip		query		int		false	"Page offset"
// @Success		200			{object}	LLMsPage
// @Failure		400			{object}	api.Error
// @Router			/llm [get]
func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
	var query FindLLMsFilterDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.Filter(r.Context(), query.ToFilter())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Get LLM
// @Description	Returns an LLM by ID.
// @Tags			llm
// @Produce		json
// @Security		BasicAuth
// @Param			id			path		string	true	"LLM ID"
// @Success		200			{object}	LLMResponse
// @Failure		404			{object}	api.Error
// @Router			/llm/{id} [get]
func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.FindByID(r.Context(), h.llmID(r))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Create LLM
// @Description	Creates a new LLM configuration.
// @Tags			llm
// @Accept			json
// @Produce		json
// @Security		BasicAuth
// @Param			body		body		CreateLLMDTO	true	"LLM data"
// @Success		201			{object}	LLMResponse
// @Failure		400			{object}	api.Error
// @Router			/llm [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateLLMDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Update LLM
// @Description	Updates an LLM by ID.
// @Tags			llm
// @Accept			json
// @Security		BasicAuth
// @Param			id			path	string			true	"LLM ID"
// @Param			body		body	UpdateLLMDTO	true	"LLM data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/llm/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateLLMDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.Update(r.Context(), h.llmID(r), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Delete LLM
// @Description	Deletes an LLM by ID.
// @Tags			llm
// @Security		BasicAuth
// @Param			id			path	string	true	"LLM ID"
// @Success		204
// @Failure		404	{object}	api.Error
// @Router			/llm/{id} [delete]
func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.DeleteByID(r.Context(), h.llmID(r)); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		List LLM drivers
// @Description	Lists the available LLM provider drivers.
// @Tags			llm
// @Produce		json
// @Security		BasicAuth
// @Success		200			{array}		DriverInfo
// @Failure		400			{object}	api.Error
// @Router			/llm/drivers [get]
func (h *Handler) listDrivers(w http.ResponseWriter, r *http.Request) error {
	drivers, err := h.service.ListDrivers(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, drivers)
}
