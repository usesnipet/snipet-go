package client

import (
	"github.com/usesnipet/snipet/internal/filter"
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

type FindClientsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindClientsFilterDTO) ToFilter() *filter.Options[model.Client] {
	return filter.New[model.Client](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
