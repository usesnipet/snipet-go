package model

import "time"

type Account struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	UserID     string `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider   string `gorm:"type:varchar(255);not null;uniqueIndex:idx_accounts_provider_external_id,priority:1" json:"provider"`
	ExternalID string `gorm:"type:varchar(255);not null;uniqueIndex:idx_accounts_provider_external_id,priority:2" json:"external_id"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
}
