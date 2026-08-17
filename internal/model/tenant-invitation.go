package model

import "time"

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
)

type TenantInvitation struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID  string           `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Email     string           `gorm:"type:varchar(255);not null" json:"email"`
	Token     string           `gorm:"type:varchar(255);not null;unique" json:"-"`
	Role      MemberRole       `gorm:"type:varchar(255);not null" json:"role"`
	Status    InvitationStatus `gorm:"type:varchar(255);not null;index" json:"status"`
	ExpiresAt time.Time        `gorm:"type:timestamp;not null" json:"expires_at"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
