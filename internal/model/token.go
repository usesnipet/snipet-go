package model

import (
	"time"

	"github.com/usesnipet/snipet/pkg/jsonx"
)

type TokenType string

const (
	TokenTypeRefresh         TokenType = "refresh"
	TokenTypeActivateAccount TokenType = "activate_account"
	TokenTypeResetPassword   TokenType = "reset_password"
)

type Token struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Type      TokenType     `gorm:"type:varchar(255);not null;index" json:"type"`
	Hash      string        `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	UserID    string        `gorm:"type:uuid;not null;index" json:"user_id"`
	ExpiresAt time.Time     `gorm:"type:timestamp;not null" json:"expires_at"`
	RevokedAt *time.Time    `gorm:"type:timestamp" json:"revoked_at"`
	Metadata  jsonx.JSONMap `gorm:"type:jsonb;not null" json:"metadata"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
