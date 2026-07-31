package msg

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Role      Role      `json:"role" gorm:"type:varchar(50);not null"`
	Sequence  int       `json:"sequence" gorm:"type:integer;not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Timestamp time.Time `json:"timestamp" gorm:"type:timestamp;not null;default:now()"`
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
