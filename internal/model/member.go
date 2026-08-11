package model

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Member struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	UserID   string `gorm:"type:uuid;not null;uniqueIndex:idx_members_user_tenant,priority:1" json:"user_id"`
	TenantID string `gorm:"type:uuid;not null;index;uniqueIndex:idx_members_user_tenant,priority:2" json:"tenant_id"`
	Role     Role   `gorm:"type:varchar(255);not null" json:"role"`
	IsActive bool   `gorm:"type:boolean;not null;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
