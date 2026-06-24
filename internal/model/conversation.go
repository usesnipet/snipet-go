package model

import (
	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/util"
)

type Conversation struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ClientID uuid.UUID    `gorm:"type:uuid;not null;index" json:"client_id"`
	MemoryID uuid.UUID    `gorm:"type:uuid;not null;index" json:"memory_id"`
	BotID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"bot_id"`
	Metadata util.JSONMap `gorm:"type:jsonb;not null;serializer:json" json:"metadata"`

	Client                  Client                   `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientUserConversations []ClientUserConversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ConversationMessages    []ConversationMessage    `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Memory                  Memory                   `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Bot                     Bot                      `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
