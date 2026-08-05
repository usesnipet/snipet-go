package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type User struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name     string        `gorm:"type:varchar(255);not null" json:"name"`
	Picture  *string       `gorm:"type:text" json:"picture"`
	Email    *string       `gorm:"type:text" json:"email"`
	Metadata jsonx.JSONMap `gorm:"type:jsonb;not null" json:"metadata"`

	UserToSessions []UserToSession `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientToUsers  []ClientToUser  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type ClientToUser struct {
	ClientID   string  `gorm:"primaryKey" json:"client_id"`
	UserID     string  `gorm:"primaryKey" json:"user_id"`
	ExternalID *string `gorm:"type:varchar(255);index" json:"external_id"`

	Client Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserToSession struct {
	UserID    string `gorm:"primaryKey" json:"user_id"`
	SessionID string `gorm:"primaryKey" json:"session_id"`

	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Session Session `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
