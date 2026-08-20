package model

import (
	"time"

	"github.com/usesnipet/snipet/pkg/jsonx"
)

type AppUserRefreshToken struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Hash      string        `gorm:"type:text;not null;uniqueIndex" json:"-"`
	AppUserID string        `gorm:"type:uuid;not null;index" json:"app_user_id"`
	ExpiresAt time.Time     `gorm:"type:timestamp;not null" json:"expires_at"`
	CreatedAt time.Time     `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	RevokedAt *time.Time    `gorm:"type:timestamp" json:"revoked_at"`
	Metadata  jsonx.JSONMap `gorm:"type:jsonb;not null" json:"metadata"`

	AppUser *AppUser `gorm:"foreignKey:AppUserID;references:ID;constraint:OnDelete:CASCADE" json:"app_user"`
}
