package llm

import (
	"time"

	"github.com/google/uuid"
)

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type Message struct {
	ID        string      `json:"id"`
	Role      MessageRole `json:"role"`
	Sequence  int         `json:"sequence"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewMessage(role MessageRole, content string, options ...MessageOption) Message {
	message := Message{
		ID:        uuid.NewString(),
		Role:      role,
		Sequence:  0,
		Content:   content,
		Timestamp: time.Now(),
	}
	for _, option := range options {
		option(&message)
	}
	return message
}
