package model

import "github.com/google/uuid"

type Conversation struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	MemoryID uuid.UUID `gorm:"type:uuid;not null;index" json:"memory_id"`
	BotID    uuid.UUID `gorm:"type:uuid;not null;index" json:"bot_id"`

	ClientUserConversations []ClientUserConversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ConversationMessages    []ConversationMessage    `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Memory                  Memory                   `gorm:"foreignKey:MemoryID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Bot                     Bot                      `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
