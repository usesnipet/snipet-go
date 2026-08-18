package model

import (
	"time"

	"github.com/usesnipet/snipet/pkg/jsonx"
)

type ClientUserRefreshToken struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Hash         string        `gorm:"type:text;not null;uniqueIndex" json:"-"`
	ClientUserID string        `gorm:"type:uuid;not null;index" json:"client_user_id"`
	ExpiresAt    time.Time     `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt    time.Time     `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	RevokedAt    *time.Time    `gorm:"type:timestamp" json:"revoked_at"`
	Metadata     jsonx.JSONMap `gorm:"type:jsonb;not null" json:"metadata"`

	ClientUser *ClientUser `gorm:"foreignKey:ClientUserID;references:ID;constraint:OnDelete:CASCADE" json:"client_user"`
}
