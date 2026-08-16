package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
)

type Handler struct {
	service        *Service
	userMiddleware api.MiddlewareFunc
}

func NewHandler(service *Service, userMiddleware api.MiddlewareFunc) api.Handler {
	return &Handler{service: service, userMiddleware: userMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", serve(h.register))
		r.Post("/login", serve(h.login))
		r.Get("/{provider}", serve(h.authorizationURL))
		r.Get("/{provider}/callback", serve(h.callback))
		r.Post("/refresh", serve(h.refresh))
		r.Post("/logout", serve(h.logout))
		r.Post("/password/forgot", serve(h.forgotPassword))
		r.Post("/password/reset", serve(h.resetPassword))
		r.Post("/activate", serve(h.activate))
		r.Post("/activate/resend", serve(h.resendActivation))

		r.Group(func(r chi.Router) {
			r.Use(h.userMiddleware)
			r.Put("/password", serve(h.setPassword))
		})
	})
}

// @Summary		Register
// @Description	Registers a new tenant-staff user. Account starts inactive; an activation email is sent.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			body	body		RegisterDTO	true	"Registration data"
// @Success		201		{object}	RegisterResponse
// @Failure		400		{object}	api.Error
// @Router			/auth/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) error {
	var dto RegisterDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Register(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, data)
}

// @Summary		Login
// @Description	Logs in with email + password. Blocked until the account is activated.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			body	body		LoginDTO	true	"Credentials"
// @Success		200		{object}	AuthenticateResponse
// @Failure		401		{object}	api.Error
// @Router			/auth/login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) error {
	var dto LoginDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Login(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		OAuth authorization URL
// @Description	Returns the authorization URL to start an OAuth login with the given provider.
// @Tags			auth
// @Produce		json
// @Param			provider	path		string	true	"Provider name"	Enums(google, github)
// @Success		200			{object}	map[string]string
// @Failure		404			{object}	api.Error
// @Router			/auth/{provider} [get]
func (h *Handler) authorizationURL(w http.ResponseWriter, r *http.Request) error {
	provider := ProviderName(chi.URLParam(r, "provider"))
	url, err := h.service.GetAuthorizationURL(r.Context(), provider)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, map[string]string{"url": url})
}

// @Summary		OAuth callback
// @Description	Completes an OAuth login/link for the given provider.
// @Tags			auth
// @Produce		json
// @Param			provider	path		string	true	"Provider name"	Enums(google, github)
// @Param			code		query		string	true	"Authorization code"
// @Param			state		query		string	true	"State"
// @Success		200			{object}	AuthenticateResponse
// @Failure		400			{object}	api.Error
// @Router			/auth/{provider}/callback [get]
func (h *Handler) callback(w http.ResponseWriter, r *http.Request) error {
	provider := ProviderName(chi.URLParam(r, "provider"))
	var query ProviderCallbackQueryDTO
	if err := api.ParseQuery(r, &query); err != nil {
		return err
	}
	data, err := h.service.AuthenticateCallback(r.Context(), provider, query.Code)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Refresh
// @Description	Exchanges a refresh token for a new access/refresh token pair.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			body	body		RefreshDTO	true	"Refresh token"
// @Success		200		{object}	AuthenticateResponse
// @Failure		401		{object}	api.Error
// @Router			/auth/refresh [post]
func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) error {
	var dto RefreshDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	data, err := h.service.Refresh(r.Context(), dto)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, data)
}

// @Summary		Logout
// @Description	Revokes a refresh token.
// @Tags			auth
// @Accept			json
// @Param			body	body	RefreshDTO	true	"Refresh token"
// @Success		204
// @Failure		401	{object}	api.Error
// @Router			/auth/logout [post]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) error {
	var dto RefreshDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.Logout(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Set password
// @Description	Sets the current user's password. No current-password check.
// @Tags			auth
// @Accept			json
// @Security		BearerAuth
// @Param			body	body	SetPasswordDTO	true	"New password"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/auth/password [put]
func (h *Handler) setPassword(w http.ResponseWriter, r *http.Request) error {
	var dto SetPasswordDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.SetPassword(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Forgot password
// @Description	Always succeeds; emails a reset link if the email is registered.
// @Tags			auth
// @Accept			json
// @Param			body	body	ForgotPasswordDTO	true	"Email"
// @Success		204
// @Router			/auth/password/forgot [post]
func (h *Handler) forgotPassword(w http.ResponseWriter, r *http.Request) error {
	var dto ForgotPasswordDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.ForgotPassword(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Reset password
// @Description	Resets a password via an emailed reset token.
// @Tags			auth
// @Accept			json
// @Param			body	body	ResetPasswordDTO	true	"Reset token + new password"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/auth/password/reset [post]
func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) error {
	var dto ResetPasswordDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.ResetPassword(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Activate account
// @Description	Activates an account via an emailed activation token.
// @Tags			auth
// @Accept			json
// @Param			body	body	ActivateAccountDTO	true	"Activation token"
// @Success		204
// @Failure		400	{object}	api.Error
// @Router			/auth/activate [post]
func (h *Handler) activate(w http.ResponseWriter, r *http.Request) error {
	var dto ActivateAccountDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.Activate(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}

// @Summary		Resend activation
// @Description	Always succeeds; re-issues the activation email if the email is registered.
// @Tags			auth
// @Accept			json
// @Param			body	body	ResendActivationDTO	true	"Email"
// @Success		204
// @Router			/auth/activate/resend [post]
func (h *Handler) resendActivation(w http.ResponseWriter, r *http.Request) error {
	var dto ResendActivationDTO
	if err := api.ParseBody(r, &dto); err != nil {
		return err
	}
	if err := h.service.ResendActivation(r.Context(), dto); err != nil {
		return err
	}
	return api.WriteNoContent(w)
}
