package message

import (
	"time"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/runtime/tool"
)

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleTool      MessageRole = "tool"
	MessageRoleFinal     MessageRole = "final"
)

type Message struct {
	ID         string       `json:"id"`
	Role       MessageRole  `json:"role"`
	Sequence   int          `json:"sequence"`
	Content    string       `json:"content"`
	ToolCalls  []tool.Call  `json:"toolCalls"`
	ToolResult *tool.Result `json:"toolResult"`
	Timestamp  time.Time    `json:"timestamp"`
}

func NewMessage(role MessageRole, content string, options ...MessageOption) Message {
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
