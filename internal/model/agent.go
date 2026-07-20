package model

import (
	"github.com/usesnipet/snipet/internal/runtime"
)

type AgentConfiguration struct {
	LLMs  []runtime.LLMConfig `json:"llms"`
	Tools runtime.ToolConfig  `json:"tools"`
}

type Agent struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name          string             `gorm:"type:varchar(255);not null" json:"name"`
	Description   string             `gorm:"type:text;not null" json:"description"`
	Instructions  string             `gorm:"type:text;not null" json:"instructions"`
	Configuration AgentConfiguration `gorm:"type:jsonb;not null;serializer:json" json:"configuration"`

	AgentToKnowledge []AgentToKnowledge `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (a Agent) ToRuntimeAgent() runtime.Agent {
	return runtime.NewAgent(
		a.Name,
		a.Description,
		a.Instructions,
		a.Configuration.Tools,
		a.Configuration.LLMs,
	)
}

type AgentToKnowledge struct {
	AgentID     string `gorm:"primaryKey" json:"agent_id"`
	KnowledgeID string `gorm:"primaryKey" json:"knowledge_id"`
	Active      bool   `gorm:"type:boolean;not null;default:true" json:"active"`

	Agent     Agent     `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Knowledge Knowledge `gorm:"foreignKey:KnowledgeID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
