package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type Session struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID string        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ClientID string        `gorm:"type:uuid;not null;index" json:"client_id"`
	AgentID  string        `gorm:"type:uuid;not null;index" json:"agent_id"`
	Metadata jsonx.JSONMap `gorm:"type:jsonb;not null;serializer:json" json:"metadata"`

	Tenant               Tenant                `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Client               Client                `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ClientUserToSessions []ClientUserToSession `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Executions           []Execution           `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Agent                *Agent                `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"agent,omitempty"`
}
