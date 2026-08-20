package client

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

// ClientResponse, ClientsPage and ClientAgentsPage exist so swagger annotations
// in this package can reference them without importing internal/model or
// internal/page directly.
type ClientResponse = model.App

type ClientsPage = page.Paginated[model.App]

type ClientAgentsPage = page.Paginated[model.Agent]

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

func (dto *FindClientsFilterDTO) ToFilter() *filter.Options[model.App] {
	return filter.New[model.App](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type ClientPublicDTO struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	AllowAnonymous bool   `json:"allow_anonymous"`
}

type FindClientAgentsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindClientAgentsFilterDTO) ToFilter() *filter.Options[model.Agent] {
	return filter.New[model.Agent](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
