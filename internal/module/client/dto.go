package client

import (
	"github.com/usesnipet/snipet/internal/model"
)

type CreateClientDTO struct {
	Name   string             `json:"name" validate:"required,max=255"`
	Config model.ClientConfig `json:"config" validate:"required"`
}

type UpdateClientDTO struct {
	Name   *string             `json:"name" validate:"omitempty,max=255"`
	Config *model.ClientConfig `json:"config" validate:"omitempty"`
}
