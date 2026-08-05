package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/runtime/execution"
	"github.com/usesnipet/snipet/internal/util"
)

type Execution struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	SessionID    *string          `gorm:"type:uuid;index" json:"session_id,omitempty"`
	AgentID      string           `gorm:"type:uuid;not null;index" json:"agent_id"`
	Status       execution.Status `gorm:"type:varchar(50);not null" json:"status"`
	ErrorMessage string           `gorm:"type:text;not null" json:"error_message,omitempty"`
	Turns        int              `gorm:"type:integer;not null;default:0" json:"turns"`
	Metadata     util.JSONMap     `gorm:"type:jsonb;not null" json:"metadata"`
	CreatedAt    time.Time        `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Session  *Session           `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Agent    *Agent             `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Messages []ExecutionMessage `gorm:"foreignKey:ExecutionID;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *Execution) ToRuntimeExecution(options ...execution.ExecutionOption) (*execution.Execution, error) {
	if e.Agent != nil {
		options = append(options, execution.WithAgent(e.Agent.ToRuntimeAgent()))
	}
	return execution.NewExecution(options...)
}
