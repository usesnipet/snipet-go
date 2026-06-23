package model

import (
	"github.com/google/uuid"
)

type ConversationMessagePartType string

const (
	ConversationMessagePartTypeText ConversationMessagePartType = "text"
)

type ConversationMessagePart struct {
	Type    ConversationMessagePartType `json:"type"`
	Content any                         `json:"content"`
}

type ConversationMessage struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ClientUserID   uuid.UUID                 `gorm:"type:uuid;not null;index" json:"client_user_id"`
	ConversationID uuid.UUID                 `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Role           string                    `gorm:"type:varchar(255);not null" json:"role"`
	Parts          []ConversationMessagePart `gorm:"type:jsonb;not null" json:"parts"`

	ClientUser   ClientUser   `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Conversation Conversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
