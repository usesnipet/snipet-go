package agent

import (
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
)

// AgentResponse and AgentsPage exist so swagger annotations in this package can
// reference them without importing internal/model or internal/page directly.
type AgentResponse = model.Agent

type AgentsPage = page.Paginated[model.Agent]

type CreateAgentDTO struct {
	Name         string   `json:"name" validate:"required,max=255"`
	Description  string   `json:"description" validate:"max=1000"`
	Instructions string   `json:"instructions" validate:"omitempty,max=1000"`
	LLMIDs       []string `json:"llm_ids" validate:"required,min=1,dive,uuid"`
}

type UpdateAgentDTO struct {
	Name         *string  `json:"name" validate:"omitempty,max=255"`
	Description  *string  `json:"description" validate:"omitempty,max=1000"`
	Instructions *string  `json:"instructions" validate:"omitempty,max=1000"`
	LLMIDs       []string `json:"llm_ids" validate:"omitempty,min=1,dive,uuid"`
}

type LinkClientToAgentDTO struct {
	ClientCode string `json:"client_code" validate:"required,max=10"`
	AgentID    string `json:"agent_id" validate:"required,uuid"`
}

type RunAgentDTO struct {
	Message string `json:"message" validate:"required"`
}

// RunInput is the internal input for starting an agent execution.
// SessionID nil means playground run (no session history).
type RunInput struct {
	AgentID   string
	SessionID *string
	Message   string
}
