package model

import "github.com/google/uuid"

type CUserSession struct {
	CUserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"client_user_id"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id"`

	CUser   CUser   `gorm:"foreignKey:CUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Session Session `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
