package appauth

import (
	"time"

	"github.com/usesnipet/snipet/internal/model"
)

// UserResponse exists so swagger annotations in this package can reference it
// without importing internal/model directly.
type UserResponse = model.AppUser

type AuthenticateResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
