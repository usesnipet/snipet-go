package model

import (
	"time"

	"github.com/usesnipet/snipet/internal/runtime/tool"
	"github.com/usesnipet/snipet/internal/runtime/transport"
)

type ExecutionMessage struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ExecutionID string                `gorm:"type:uuid;not null;index" json:"execution_id"`
	Sequence    int                   `gorm:"type:integer;not null" json:"sequence"`
	Role        transport.MessageRole `gorm:"type:varchar(50);not null" json:"role"`
	Content     string                `gorm:"type:text" json:"content"`

	ToolCalls  []tool.Call  `gorm:"type:jsonb" json:"tool_calls,omitempty"`
	ToolResult *tool.Result `gorm:"type:jsonb" json:"tool_result,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`

	Execution Execution `gorm:"foreignKey:ExecutionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *ExecutionMessage) ToRuntimeExecutionMessage() *transport.Message {
	return &transport.Message{
		ID:         e.ID,
		Timestamp:  e.CreatedAt,
		Role:       e.Role,
		Content:    e.Content,
		ToolCalls:  e.ToolCalls,
		ToolResult: e.ToolResult,
	}
}

func (e *ExecutionMessage) FromRuntimeExecutionMessage(message transport.Message) {
	e.ID = message.ID
	e.CreatedAt = message.Timestamp
	e.Role = message.Role
	e.Content = message.Content
	e.ToolCalls = message.ToolCalls
	e.ToolResult = message.ToolResult
}
