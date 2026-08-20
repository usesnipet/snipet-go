package model

import "time"

type AppStatus string

const (
	AppStatusPending     AppStatus = "pending"
	AppStatusActive      AppStatus = "active"
	AppStatusDeactivated AppStatus = "deactivated"
)

type App struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID       string     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Code           string     `gorm:"type:varchar(10);not null;unique" json:"code"`
	Name           string     `gorm:"type:varchar(255);not null" json:"name"`
	Description    string     `gorm:"type:text;not null" json:"description"`
	Status         AppStatus  `gorm:"type:varchar(255);not null" json:"status"`
	LastVerifiedAt *time.Time `gorm:"type:timestamp" json:"last_verified_at"`

	KeyID   string `gorm:"type:varchar(255);not null;unique" json:"key_id"`
	KeyHash string `gorm:"type:text;not null;unique" json:"-"`

	Tenant     *Tenant        `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"tenant"`
	Sessions   []Session      `gorm:"foreignKey:AppID;references:ID;constraint:OnDelete:CASCADE" json:"sessions"`
	AppToUsers []AppToAppUser `gorm:"foreignKey:AppID;references:ID;constraint:OnDelete:CASCADE" json:"app_to_users"`
}
