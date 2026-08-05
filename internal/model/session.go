package model

import "github.com/usesnipet/snipet/pkg/jsonx"

type Session struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ClientID string        `gorm:"type:uuid;not null;index" json:"client_id"`
	AgentID  string        `gorm:"type:uuid;not null;index" json:"agent_id"`
	Metadata jsonx.JSONMap `gorm:"type:jsonb;not null;serializer:json" json:"metadata"`

	Client         Client          `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	UserToSessions []UserToSession `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Executions     []Execution     `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Agent          *Agent          `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"agent,omitempty"`
}
