package auth

import "github.com/usesnipet/snipet/internal/util"

type AuthenticateResponse struct {
	Token string `json:"token"`
}

type AuthenticateAnonymousDTO struct {
	Name     *string      `json:"name" validate:"omitempty,max=255"`
	Picture  *string      `json:"picture" validate:"omitempty,url"`
	Email    *string      `json:"email" validate:"omitempty,email"`
	Metadata util.JSONMap `json:"metadata"`
}
