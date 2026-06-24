package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type MemoryType string

const (
	MemoryTypeBot          MemoryType = "bot"
	MemoryTypeConversation MemoryType = "conversation"
)

type Memory struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string          `gorm:"type:varchar(255);not null" json:"name"`
	Type          MemoryType      `gorm:"type:varchar(255);not null" json:"type"`
	IsDefault     bool            `gorm:"type:boolean;not null;default:false" json:"is_default"`
	Provider      string          `gorm:"type:varchar(255);not null" json:"provider"`
	Configuration json.RawMessage `gorm:"type:jsonb;not null" json:"configuration"`

	BotMemories []BotMemory `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
