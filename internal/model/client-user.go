package model

import (
	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/util"
)

type CUser struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name     string       `gorm:"type:varchar(255);not null" json:"name"`
	Metadata util.JSONMap `gorm:"type:jsonb;not null;serializer:json" json:"metadata"`

	ConversationMessages    []ConversationMessage    `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientUserConversations []ClientUserConversation `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientToUsers           []ClientToUser           `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type ClientToUser struct {
	ClientID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_client_to_users_client_user" json:"client_id"`
	ClientUserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_client_to_users_client_user" json:"client_user_id"`
	ExternalID   *string   `gorm:"type:varchar(255);index" json:"external_id"`

	Client     Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientUser CUser  `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
