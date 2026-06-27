package client

import (
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type CreateClientDTO struct {
	Name   string             `json:"name" validate:"required,max=255"`
	Config model.ClientConfig `json:"config" validate:"required"`
}

type UpdateClientDTO struct {
	Name   *string             `json:"name" validate:"omitempty,max=255"`
	Config *model.ClientConfig `json:"config" validate:"omitempty"`
}

type AuthenticateClientUserDTO struct {
	ExternalID string       `json:"external_id" validate:"required,max=255"`
	Name       string       `json:"name" validate:"omitempty,max=255"`
	Email      string       `json:"email" validate:"omitempty,email"`
	Metadata   util.JSONMap `json:"metadata" validate:"omitempty"`
}
