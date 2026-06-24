package model

import "github.com/google/uuid"

type Client struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name string `gorm:"type:varchar(255);not null" json:"name"`

	ClientBots    []ClientBot    `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Conversations []Conversation `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
