package session

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type RunSessionDTO struct {
	Message string `json:"message" validate:"required"`
}

type CreateSessionDTO struct {
	AgentID  string        `json:"agent_id" validate:"required,uuid"`
	Metadata jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type UpdateSessionDTO struct {
	AgentID  *string       `json:"agent_id" validate:"omitempty,uuid"`
	Metadata jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type SessionsFilterDTO struct {
	Take    *int     `form:"take" validate:"omitempty,min=1"`
	Skip    *int     `form:"skip" validate:"omitempty,min=0"`
	Include []string `form:"include" validate:"omitempty,dive,oneof=agent"`
}

func (dto *SessionsFilterDTO) ToFilter() *filter.Options[model.Session] {
	opts := []filter.Option{
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	}
	opts = append(opts, sessionIncludeOptions(dto.Include)...)
	return filter.New[model.Session](opts...)
}

type SessionIncludeDTO struct {
	Include []string `form:"include" validate:"omitempty,dive,oneof=agent"`
}

func (dto *SessionIncludeDTO) ToFilter() *filter.Options[model.Session] {
	return filter.New[model.Session](sessionIncludeOptions(dto.Include)...)
}

func sessionIncludeOptions(includes []string) []filter.Option {
	opts := make([]filter.Option, 0, len(includes))
	for _, include := range includes {
		switch include {
		case "agent":
			opts = append(opts, filter.Include("Agent"))
		}
	}
	return opts
}

type MessagesFilterDTO struct {
	Sort *filter.OrderDirection `form:"sort" validate:"omitempty,oneof=asc desc"`
	Take *int                   `form:"take" validate:"omitempty,min=1"`
	Skip *int                   `form:"skip" validate:"omitempty,min=0"`
}

func (dto *MessagesFilterDTO) ToFilter() *filter.Options[model.ExecutionMessage] {
	return filter.New[model.ExecutionMessage](
		filter.PtrOrderBy("timestamp", dto.Sort),
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
