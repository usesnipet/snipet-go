package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type ExecutionMessage struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ExecutionID string       `gorm:"type:uuid;not null;index" json:"execution_id"`
	Sequence    int          `gorm:"type:integer;not null" json:"sequence"`
	Role        message.Role `gorm:"type:varchar(50);not null" json:"role"`
	Content     string       `gorm:"type:text" json:"content"`

	ToolCalls  []tool.Call  `gorm:"type:jsonb" json:"tool_calls,omitempty"`
	ToolResult *tool.Result `gorm:"type:jsonb" json:"tool_result,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`

	Execution Execution `gorm:"foreignKey:ExecutionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *ExecutionMessage) ToRuntimeExecutionMessage() *message.Message {
	return &message.Message{
		ID:         e.ID,
		Sequence:   e.Sequence,
		Timestamp:  e.CreatedAt,
		Role:       e.Role,
		Content:    e.Content,
		ToolCalls:  e.ToolCalls,
		ToolResult: e.ToolResult,
	}
}

func (e *ExecutionMessage) FromRuntimeExecutionMessage(msg message.Message) *ExecutionMessage {
	e.ID = msg.ID
	e.Sequence = msg.Sequence
	e.CreatedAt = msg.Timestamp
	e.Role = msg.Role
	e.Content = msg.Content
	e.ToolCalls = msg.ToolCalls
	e.ToolResult = msg.ToolResult
	return e
}
