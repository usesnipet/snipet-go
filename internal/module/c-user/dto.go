package c_user

import (
	"github.com/usesnipet/snipet/internal/util"
)

type CreateAnonymousClientUserDTO struct {
	Name     *string      `json:"name" validate:"omitempty,max=255"`
	Metadata util.JSONMap `json:"metadata" validate:"omitempty"`
}

type CreateAuthenticatedClientUserDTO struct {
	ExternalID string       `json:"external_id" validate:"required,max=255"`
	Name       string       `json:"name" validate:"omitempty,max=255"`
	Metadata   util.JSONMap `json:"metadata" validate:"omitempty"`
}
