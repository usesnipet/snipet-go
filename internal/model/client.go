package model

import (
	"github.com/google/uuid"
)

type ClientConfig struct {
	OIDC struct {
		Issuer   *string `json:"issuer" validate:"omitempty,url"`
		Audience *string `json:"audience" validate:"omitempty,url"`
		Enabled  bool    `json:"enabled"`
	} `json:"oidc"`
	Webhook struct {
		URL     string `json:"url" validate:"omitempty,url"`
		Enabled bool   `json:"enabled"`
	} `json:"webhook"`
}

type Client struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Code   string       `gorm:"type:char(10);not null;unique" json:"code"`
	Name   string       `gorm:"type:varchar(255);not null" json:"name"`
	Config ClientConfig `gorm:"type:jsonb;not null;serializer:json" json:"config"`

	ClientToBots  []ClientToBot  `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Sessions      []Session      `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientToUsers []ClientToUser `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
