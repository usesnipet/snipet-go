package model

import "github.com/google/uuid"

type ClientBot struct {
	ClientID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_client_bots_client_bot" json:"client_id"`
	BotID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_client_bots_client_bot" json:"bot_id"`
	Client   Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Bot      Bot       `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
