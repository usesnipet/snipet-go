package model

import (
	"github.com/usesnipet/snipet/internal/util"
)

type MemoryType string

const (
	MemoryTypeBot     MemoryType = "bot"
	MemoryTypeSession MemoryType = "session"
)

type Memory struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string       `gorm:"type:varchar(255);not null" json:"name"`
	Type          MemoryType   `gorm:"type:varchar(255);not null" json:"type"`
	IsDefault     bool         `gorm:"type:boolean;not null;default:false" json:"is_default"`
	Provider      string       `gorm:"type:varchar(255);not null" json:"provider"`
	Configuration util.JSONMap `gorm:"type:jsonb;not null" json:"configuration"`

	BotMemories []BotToMemory `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
