package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

// Handler routes tenant-staff user management under /users.
//
// bearerMiddleware/adminMiddleware are injected rather than constructed
// here because they depend on the tenant-staff auth module (see
// plan-ee/module/auth.md and plan-ee/module/middleware.md), which does not
// exist yet. Bootstrap wiring of this handler is deferred until that module
// lands.
type Handler struct {
	service          *Service
	bearerMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, bearerMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{
		service:          service,
		bearerMiddleware: bearerMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.bearerMiddleware)
			r.Get("/me", serve(h.me))
			r.Put("/me/picture", serve(h.updatePicture))
		})
	})
}

// @Summary		Get current user
// @Description	Returns the authenticated tenant-staff user.
// @Tags			user
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	UserResponse
// @Failure		401	{object}	api.Error
// @Router			/users/me [get]
func (h *Handler) me(w http.ResponseWriter, r *http.Request) error {
	data, err := h.service.Me(r.Context())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Update my picture
// @Description	Updates the authenticated user's profile picture.
// @Tags			user
// @Accept			json
// @Security		BearerAuth
// @Param			body	body	UpdateProfilePictureDTO	true	"Picture data"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/users/me/picture [put]
func (h *Handler) updatePicture(w http.ResponseWriter, r *http.Request) error {
	var dto UpdateProfilePictureDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.UpdateMyPicture(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
