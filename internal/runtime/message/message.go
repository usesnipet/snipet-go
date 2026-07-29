package message

import (
	"time"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type Role string

const (
	MessageRoleSystem    Role = "system"
	MessageRoleUser      Role = "user"
	MessageRoleAssistant Role = "assistant"
	MessageRoleTool      Role = "tool"
	MessageRoleFinal     Role = "final"
)

type Message struct {
	ID         string       `json:"id"`
	Role       Role         `json:"role"`
	Sequence   int          `json:"sequence"`
	Content    string       `json:"content"`
	ToolCalls  []tool.Call  `json:"toolCalls"`
	ToolResult *tool.Result `json:"toolResult"`
	Timestamp  time.Time    `json:"timestamp"`
}

func NewMessage(role Role, content string, options ...MessageOption) Message {
	message := Message{
		ID:         uuid.NewString(),
		Role:       role,
		Sequence:   0,
		Content:    content,
		ToolCalls:  []tool.Call{},
		ToolResult: nil,
		Timestamp:  time.Now(),
	}
	for _, option := range options {
		option(&message)
	}
	return message
}
