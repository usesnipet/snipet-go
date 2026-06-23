package model

import "github.com/google/uuid"

type BotMemory struct {
	BotID    uuid.UUID `gorm:"type:uuid;not null;index" json:"bot_id"`
	MemoryID uuid.UUID `gorm:"type:uuid;not null;index" json:"memory_id"`
	Active   bool      `gorm:"type:boolean;not null;default:true" json:"active"`

	Bot    Bot    `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Memory Memory `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
