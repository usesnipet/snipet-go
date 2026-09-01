package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/runtime/execution"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type Execution struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	SessionID      *string          `gorm:"type:uuid;index" json:"session_id,omitempty"`
	AgentID        string           `gorm:"type:uuid;not null;index" json:"agent_id"`
	Status         execution.Status `gorm:"type:varchar(50);not null" json:"status"`
	StreamMessages bool             `gorm:"type:boolean;not null;default:true" json:"stream_messages"`
	ErrorMessage   string           `gorm:"type:text;not null" json:"error_message,omitempty"`
	Turns          int              `gorm:"type:integer;not null;default:0" json:"turns"`
	Metadata       jsonx.JSONMap    `gorm:"type:jsonb;not null" json:"metadata"`
	CreatedAt      time.Time        `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Session  *Session           `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"session"`
	Agent    *Agent             `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"agent"`
	Messages []ExecutionMessage `gorm:"foreignKey:ExecutionID;constraint:OnDelete:CASCADE" json:"messages"`
}

func (e *Execution) ToRuntimeExecution(options ...execution.ExecutionOption) (*execution.Execution, error) {
	if e.Agent != nil {
		options = append(options, execution.WithAgent(e.Agent.ToRuntimeAgent()))
	}
	if e.StreamMessages {
		options = append(options, execution.WithStreamMessages(e.StreamMessages))
	}
	return execution.NewExecution(options...)
}
