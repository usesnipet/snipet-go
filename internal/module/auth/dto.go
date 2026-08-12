package auth

import (
	"time"

	"github.com/usesnipet/snipet/internal/model"
)

// UserResponse exists so swagger annotations in this package can reference
// it without importing internal/model directly.
type UserResponse = model.User

type AuthenticateResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type RegisterDTO struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SetPasswordDTO struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ResetPasswordDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ActivateAccountDTO struct {
	Token string `json:"token" validate:"required"`
}

type ResendActivationDTO struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ProviderCallbackQueryDTO struct {
	Code  string `form:"code" validate:"required"`
	State string `form:"state" validate:"required"`
}
