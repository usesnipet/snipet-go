package model

import (
	"github.com/usesnipet/snipet/internal/util"
)

type Session struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ClientID string       `gorm:"type:uuid;not null;index" json:"client_id"`
	BotID    string       `gorm:"type:uuid;not null;index" json:"bot_id"`
	Metadata util.JSONMap `gorm:"type:jsonb;not null;serializer:json" json:"metadata"`

	Client          Client           `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	UserToSessions  []UserToSession  `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	SessionMessages []SessionMessage `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Bot             Bot              `gorm:"foreignKey:BotID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
