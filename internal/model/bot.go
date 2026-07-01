package model

import (
	"github.com/usesnipet/snipet/internal/util"
)

type BotConfiguration struct {
	util.JSONMap
	LLMs []any `json:"llms"`
}

type Bot struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string           `gorm:"type:varchar(255);not null" json:"name"`
	Description   string           `gorm:"type:text;not null" json:"description"`
	Configuration BotConfiguration `gorm:"type:jsonb;not null;serializer:json" json:"configuration"`

	BotToMemories []BotToMemory `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type BotToMemory struct {
	BotID    string `gorm:"primaryKey" json:"bot_id"`
	MemoryID string `gorm:"primaryKey" json:"memory_id"`
	Active   bool   `gorm:"type:boolean;not null;default:true" json:"active"`

	Bot    Bot    `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Memory Memory `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
