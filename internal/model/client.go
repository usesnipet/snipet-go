package model

import (
	"github.com/google/uuid"
)

type Client struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Code       string `gorm:"type:char(10);not null;unique" json:"code"`
	WebhookURL string `gorm:"type:text" json:"webhook_url"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`

	ClientBots    []ClientBot    `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Conversations []Conversation `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
