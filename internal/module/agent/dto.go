package agent

import "github.com/usesnipet/snipet/internal/model"

type CreateAgentDTO struct {
	Name          string                   `json:"name" validate:"required,max=255"`
	Description   string                   `json:"description" validate:"max=1000"`
	Configuration model.AgentConfiguration `json:"configuration" validate:"required"`
}

type UpdateAgentDTO struct {
	Name          *string                   `json:"name" validate:"omitempty,max=255"`
	Description   *string                   `json:"description" validate:"omitempty,max=1000"`
	Configuration *model.AgentConfiguration `json:"configuration"`
}

type LinkClientToAgentDTO struct {
	ClientCode string `json:"client_code" validate:"required,max=10"`
	AgentID    string `json:"agent_id" validate:"required,uuid"`
}
