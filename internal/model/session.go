package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type Session struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	TenantID string        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AppID    string        `gorm:"type:uuid;not null;index" json:"app_id"`
	AgentID  string        `gorm:"type:uuid;not null;index" json:"agent_id"`
	Metadata jsonx.JSONMap `gorm:"type:jsonb;not null;serializer:json" json:"metadata"`

	Tenant            *Tenant            `gorm:"foreignKey:TenantID;references:ID;constraint:OnDelete:CASCADE" json:"tenant"`
	App               *App               `gorm:"foreignKey:AppID;references:ID;constraint:OnDelete:CASCADE" json:"app"`
	AppUserToSessions []AppUserToSession `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"app_user_to_sessions"`
	Executions        []Execution        `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"executions"`
	Agent             *Agent             `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"agent"`
}
