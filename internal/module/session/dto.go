package session

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type CreateSessionDTO struct {
	MemoryID string       `json:"memory_id" validate:"required,uuid"`
	BotID    string       `json:"bot_id" validate:"required,uuid"`
	Metadata util.JSONMap `json:"metadata" validate:"omitempty"`
}

type FindSessionsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindSessionsFilterDTO) ToFilter() *filter.Options[model.Session] {
	return filter.New[model.Session](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type FindMessagesFilterDTO struct {
	Sort *filter.OrderDirection `form:"sort" validate:"omitempty,oneof=asc desc"`
	Take *int                   `form:"take" validate:"omitempty,min=1"`
	Skip *int                   `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindMessagesFilterDTO) ToFilter() *filter.Options[model.SessionMessage] {
	return filter.New[model.SessionMessage](
		filter.PtrOrderBy("created_at", dto.Sort),
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
