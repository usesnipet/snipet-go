package agent

import (
	"github.com/usesnipet/snipet/internal/util"
)

type LLMConfigDTO struct {
	Key    string       `json:"key" validate:"required"`
	Config util.JSONMap `json:"config" validate:"required"`
}

type ToolConfigDTO map[string]util.JSONMap

type CreateAgentDTO struct {
	Name         string         `json:"name" validate:"required,max=255"`
	Description  string         `json:"description" validate:"max=1000"`
	Instructions string         `json:"instructions" validate:"omitempty,max=1000"`
	LLMs         []LLMConfigDTO `json:"llms" validate:"required"`
	Tools        ToolConfigDTO  `json:"tools" validate:"required"`
}

type UpdateAgentDTO struct {
	Name         *string        `json:"name" validate:"omitempty,max=255"`
	Description  *string        `json:"description" validate:"omitempty,max=1000"`
	Instructions *string        `json:"instructions" validate:"omitempty,max=1000"`
	LLMs         []LLMConfigDTO `json:"llms" validate:"omitempty"`
	Tools        ToolConfigDTO  `json:"tools" validate:"omitempty"`
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
