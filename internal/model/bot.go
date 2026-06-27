package model

import (
	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/util"
)

type BotConfiguration struct {
	util.JSONMap
	LLMs []any `json:"llms"`
}

type Bot struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string           `gorm:"type:varchar(255);not null" json:"name"`
	Description   string           `gorm:"type:text;not null" json:"description"`
	Configuration BotConfiguration `gorm:"type:jsonb;not null" json:"configuration"`

	BotMemories []BotMemory `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientBots  []ClientBot `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
