package auth

import (
	"time"

	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type AuthenticateResponse struct {
	AccessToken           string     `json:"access_token"`
	AccessTokenExpiresAt  time.Time  `json:"access_token_expires_at"`
	RefreshToken          string     `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time  `json:"refresh_token_expires_at"`
	User                  model.User `json:"user"`
}

type AuthenticateAnonymousDTO struct {
	Name     *string      `json:"name" validate:"omitempty,max=255"`
	Picture  *string      `json:"picture" validate:"omitempty,url"`
	Email    *string      `json:"email" validate:"omitempty,email"`
	Metadata util.JSONMap `json:"metadata"`
}

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
