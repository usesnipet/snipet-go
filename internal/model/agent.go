package model

import (
	"sort"

	"github.com/usesnipet/snipet/internal/runtime/execution"
)

type Agent struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Description  string `gorm:"type:text;not null" json:"description"`
	Instructions string `gorm:"type:text;not null" json:"instructions"`

	AgentToLLMs      []AgentToLLM       `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"llms"`
	AgentToKnowledge []AgentToKnowledge `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"knowledge"`
}

func (a Agent) ToRuntimeAgent() *execution.Agent {
	llms := make([]execution.LLMConfig, 0, len(a.AgentToLLMs))
	rels := append([]AgentToLLM(nil), a.AgentToLLMs...)
	sort.SliceStable(rels, func(i, j int) bool {
		return rels[i].Priority < rels[j].Priority
	})
	for _, rel := range rels {
		llms = append(llms, execution.LLMConfig{
			Key:    rel.LLM.Provider,
			Config: rel.LLM.Configuration,
		})
	}
	return execution.NewAgent(
		a.Name,
		a.Description,
		a.Instructions,
		llms[0],
	)
}

type AgentToLLM struct {
	AgentID  string `gorm:"type:uuid;primaryKey;index:idx_agent_to_llms_agent_priority,priority:1" json:"-"`
	LLMID    string `gorm:"type:uuid;primaryKey;column:llm_id" json:"llm_id"`
	Priority int    `gorm:"type:integer;not null;index:idx_agent_to_llms_agent_priority,priority:2" json:"priority"`

	Agent Agent `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	LLM   LLM   `gorm:"foreignKey:LLMID;references:ID;constraint:OnDelete:CASCADE" json:"llm"`
}

func (AgentToLLM) TableName() string {
	return "agent_to_llms"
}

type AgentToKnowledge struct {
	AgentID     string `gorm:"primaryKey" json:"agent_id"`
	KnowledgeID string `gorm:"primaryKey" json:"knowledge_id"`
	Active      bool   `gorm:"type:boolean;not null;default:true" json:"active"`

	Agent     Agent     `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Knowledge Knowledge `gorm:"foreignKey:KnowledgeID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}
