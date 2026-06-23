package model

import "github.com/google/uuid"

type ClientUser struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name       string `gorm:"type:varchar(255);not null" json:"name"`
	Anonymous  bool   `gorm:"type:boolean;not null;default:false" json:"anonymous"`
	SessionID  string `gorm:"type:varchar(255);not null" json:"session_id"`
	ExternalID string `gorm:"type:varchar(255);not null" json:"external_id"`

	ConversationMessages    []ConversationMessage    `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientUserConversations []ClientUserConversation `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
