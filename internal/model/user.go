package model

import (
	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/util"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name     string       `gorm:"type:varchar(255);not null" json:"name"`
	Picture  *string      `gorm:"type:text" json:"picture"`
	Email    *string      `gorm:"type:text" json:"email"`
	Metadata util.JSONMap `gorm:"type:jsonb;not null" json:"metadata"`

	SessionMessages []SessionMessage `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	UserToSessions  []UserToSession  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientToUsers   []ClientToUser   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type ClientToUser struct {
	ClientID   uuid.UUID `gorm:"primaryKey" json:"client_id"`
	UserID     uuid.UUID `gorm:"primaryKey" json:"user_id"`
	ExternalID *string   `gorm:"type:varchar(255);index" json:"external_id"`

	Client Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserToSession struct {
	UserID    uuid.UUID `gorm:"primaryKey" json:"user_id"`
	SessionID uuid.UUID `gorm:"primaryKey" json:"session_id"`

	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Session Session `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
