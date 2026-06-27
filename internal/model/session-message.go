package model

import (
	"github.com/google/uuid"
)

type SessionMessagePartType string

const (
	SessionMessagePartTypeText SessionMessagePartType = "text"
)

type SessionMessagePart struct {
	Type    SessionMessagePartType `json:"type"`
	Content any                    `json:"content"`
}

type SessionMessage struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	UserID    uuid.UUID            `gorm:"type:uuid;not null;index" json:"c_user_id"`
	SessionID uuid.UUID            `gorm:"type:uuid;not null;index" json:"session_id"`
	Role      string               `gorm:"type:varchar(255);not null" json:"role"`
	Parts     []SessionMessagePart `gorm:"type:jsonb;not null" json:"parts"`

	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Session Session `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
