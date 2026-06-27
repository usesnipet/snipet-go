package session

import (
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/util"
)

type CreateSessionDTO struct {
	MemoryID string       `json:"memory_id" validate:"required,uuid"`
	BotID    string       `json:"bot_id" validate:"required,uuid"`
	Metadata util.JSONMap `json:"metadata" validate:"omitempty"`
}

type FindMessagesFilterDTO struct {
	Sort filter.OrderDirection `form:"sort"`
	Take int                   `form:"take"`
	Skip int                   `form:"skip"`
}
