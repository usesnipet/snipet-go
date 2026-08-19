package model

type ClientStatus string

const (
	ClientStatusPending          ClientStatus = "pending"
	ClientStatusActivating       ClientStatus = "activating"
	ClientStatusActive           ClientStatus = "active"
	ClientStatusActivationFailed ClientStatus = "activation_failed"
	ClientStatusDeactivated      ClientStatus = "deactivated"
)

type Client struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID     string       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name         string       `gorm:"type:varchar(255);not null" json:"name"`
	Description  string       `gorm:"type:text;not null" json:"description"`
	Status       ClientStatus `gorm:"type:varchar(255);not null" json:"status"`
	ErrorMessage string       `gorm:"type:text;not null" json:"error_message"`
	Provider     string       `gorm:"type:varchar(255);not null" json:"provider"`

	KeyID   string `gorm:"type:varchar(255);not null;unique" json:"key_id"`
	KeyHash string `gorm:"type:text;not null;unique" json:"-"`

	Tenant        *Tenant              `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"tenant"`
	Sessions      []Session            `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"sessions"`
	ClientToUsers []ClientToClientUser `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"client_to_users"`
}
