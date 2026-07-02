package session

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/util"
)

type SendMessageDTO struct {
	Message string `json:"message" validate:"required"`
}

type CreateSessionDTO struct {
	BotID    string       `json:"bot_id" validate:"required,uuid"`
	Metadata util.JSONMap `json:"metadata" validate:"omitempty"`
}

type SessionsFilterDTO struct {
	Take *int `form:"take" validate:"omitempty,min=1"`
	Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *SessionsFilterDTO) ToFilter() *filter.Options[model.Session] {
	return filter.New[model.Session](
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}

type MessagesFilterDTO struct {
	Sort *filter.OrderDirection `form:"sort" validate:"omitempty,oneof=asc desc"`
	Take *int                   `form:"take" validate:"omitempty,min=1"`
	Skip *int                   `form:"skip" validate:"omitempty,min=0"`
}

func (dto *MessagesFilterDTO) ToFilter() *filter.Options[model.SessionMessage] {
	return filter.New[model.SessionMessage](
		filter.PtrOrderBy("created_at", dto.Sort),
		filter.PtrTake(dto.Take),
		filter.PtrSkip(dto.Skip),
	)
}
