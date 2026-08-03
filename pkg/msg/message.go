package msg

import (
	"time"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Role      Role      `json:"role" gorm:"type:varchar(50);not null"`
	Sequence  int       `json:"sequence" gorm:"type:integer;not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Timestamp time.Time `json:"timestamp" gorm:"type:timestamp;not null;default:now()"`

	// ToolCalls is set on assistant messages that requested tool calls.
	ToolCalls []tool.Call `json:"tool_calls,omitempty" gorm:"type:jsonb;serializer:json"`
	// ToolCallID links a RoleTool message back to the ToolCall.ID it answers.
	ToolCallID string `json:"tool_call_id,omitempty" gorm:"type:varchar(255)"`

	Final bool `json:"final" gorm:"-"`
}

func (m *Message) IsFinal() bool {
	return m.Final
}

type MessageOption func(*Message)

func WithID(id string) MessageOption {
	return func(message *Message) {
		message.ID = id
	}
}

func WithSequence(sequence int) MessageOption {
	return func(message *Message) {
		message.Sequence = sequence
	}
}

func WithTimestamp(timestamp time.Time) MessageOption {
	return func(message *Message) {
		message.Timestamp = timestamp
	}
}

func WithFinal() MessageOption {
	return func(message *Message) {
		message.Final = true
	}
}

func WithToolCalls(toolCalls []tool.Call) MessageOption {
	return func(message *Message) {
		message.ToolCalls = toolCalls
	}
}

func WithToolCallID(id string) MessageOption {
	return func(message *Message) {
		message.ToolCallID = id
	}
}

func NewMessage(role Role, content string, options ...MessageOption) Message {
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
