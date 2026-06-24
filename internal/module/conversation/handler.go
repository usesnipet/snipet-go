package conversation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/module/client"
)

type Handler struct {
	service       *Service
	clientService *client.Service
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/client/{client_id}/conversation", func(r chi.Router) {
		// r.Use(ClientExistsMiddleware(h.clientService))

		r.Get("/", serve(h.findBy))
		r.Post("/", serve(h.create))
		r.Get("/{id}", serve(h.findByID))
		r.Delete("/{id}", serve(h.deleteByID))

		r.Get("/{id}/messages", serve(h.findMessages))
	})
}

func (h *Handler) findMessages(w http.ResponseWriter, r *http.Request) error {
	clientID := chi.URLParam(r, "client_id")
	conversationID := chi.URLParam(r, "id")
	var filter FindMessagesFilterDTO
	if err := api.ParseQuery(r, &filter); err != nil {
		return err
	}
	data, err := h.service.FindMessages(r.Context(), clientID, conversationID, filter)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findBy(w http.ResponseWriter, r *http.Request) error {
	clientID := chi.URLParam(r, "client_id")

	data, err := h.service.FindBy(r.Context(), clientID)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) error {
	clientID := chi.URLParam(r, "client_id")
	conversationID := chi.URLParam(r, "id")
	data, err := h.service.FindByID(r.Context(), clientID, conversationID)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) error {
	var dto CreateConversationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Create(r.Context(), chi.URLParam(r, "client_id"), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

func (h *Handler) deleteByID(w http.ResponseWriter, r *http.Request) error {
	clientID := chi.URLParam(r, "client_id")
	conversationID := chi.URLParam(r, "id")
	if err := h.service.DeleteByID(r.Context(), clientID, conversationID); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

func NewHandler(service *Service) api.Handler {
	return &Handler{service: service}
}
