package model

import "github.com/google/uuid"

type ClientUserConversation struct {
	ClientUserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"client_user_id"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index" json:"conversation_id"`

	ClientUser   ClientUser   `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Conversation Conversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
