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
	Configuration BotConfiguration `gorm:"type:jsonb;not null;serializer:json" json:"configuration"`

	BotToMemories []BotToMemory `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientToBots  []ClientToBot `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type BotToMemory struct {
	BotID    uuid.UUID `gorm:"type:uuid;not null;index" json:"bot_id"`
	MemoryID uuid.UUID `gorm:"type:uuid;not null;index" json:"memory_id"`
	Active   bool      `gorm:"type:boolean;not null;default:true" json:"active"`

	Bot    Bot    `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Memory Memory `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type ClientToBot struct {
	ClientID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_client_bots_client_bot" json:"client_id"`
	BotID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_client_bots_client_bot" json:"bot_id"`
	Client   Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Bot      Bot       `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
