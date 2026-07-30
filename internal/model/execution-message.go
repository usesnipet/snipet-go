package model

import (
	"time"

	"github.com/usesnipet/snipet/pkg/driver/llm"
)

type ExecutionMessage struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ExecutionID string          `gorm:"type:uuid;not null;index" json:"execution_id"`
	Sequence    int             `gorm:"type:integer;not null" json:"sequence"`
	Role        llm.MessageRole `gorm:"type:varchar(50);not null" json:"role"`
	Content     string          `gorm:"type:text" json:"content"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`

	Execution Execution `gorm:"foreignKey:ExecutionID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *ExecutionMessage) ToRuntimeMessage() *llm.Message {
	return &llm.Message{
		ID:        e.ID,
		Sequence:  e.Sequence,
		Timestamp: e.CreatedAt,
		Role:      e.Role,
		Content:   e.Content,
	}
}

func (e *ExecutionMessage) FromRuntimeMessage(msg llm.Message) *ExecutionMessage {
	e.ID = msg.ID
	e.Sequence = msg.Sequence
	e.CreatedAt = msg.Timestamp
	e.Role = msg.Role
	e.Content = msg.Content
	return e
}
