package model

import (
	"time"
)

type APIKey struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID string `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	KeyID    string `gorm:"type:varchar(255);not null;unique" json:"key_id"`
	Key      string `gorm:"type:text;not null;unique" json:"-"`
	Active   bool   `gorm:"type:boolean;not null;default:true" json:"active"`

	ExpiresAt *time.Time `gorm:"type:timestamp" json:"expires_at"`
	CreatedAt time.Time  `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"tenant"`
}
