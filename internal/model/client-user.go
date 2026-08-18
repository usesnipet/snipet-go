package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type ClientUser struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name     string        `gorm:"type:varchar(255);not null" json:"name"`
	Picture  *string       `gorm:"type:text" json:"picture"`
	Email    *string       `gorm:"type:text" json:"email"`
	Metadata jsonx.JSONMap `gorm:"type:jsonb;not null" json:"metadata"`

	ClientUserToSessions []ClientUserToSession `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"client_user_to_sessions"`
	ClientToClientUsers  []ClientToClientUser  `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"client_to_client_users"`
}

type ClientToClientUser struct {
	ClientID     string  `gorm:"primaryKey" json:"client_id"`
	ClientUserID string  `gorm:"primaryKey" json:"client_user_id"`
	ExternalID   *string `gorm:"type:varchar(255);index" json:"external_id"`

	Client     *Client     `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"client"`
	ClientUser *ClientUser `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"client_user"`
}

type ClientUserToSession struct {
	ClientUserID string `gorm:"primaryKey" json:"user_id"`
	SessionID    string `gorm:"primaryKey" json:"session_id"`

	ClientUser *ClientUser `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"client_user"`
	Session    *Session    `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"session"`
}
