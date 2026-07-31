package runtime

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/msg"
)

var validate = validator.New()

type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"   // the execution is pending to start
	ExecutionStatusRunning   ExecutionStatus = "running"   // the execution is running
	ExecutionStatusCompleted ExecutionStatus = "completed" // the execution is completed
	ExecutionStatusFailed    ExecutionStatus = "failed"    // the execution is failed
	ExecutionStatusMaxTurns  ExecutionStatus = "max_turns" // the execution has reached the maximum number of turns
	ExecutionStatusCancelled ExecutionStatus = "cancelled" // the execution was cancelled via context
)

type ExecutionConfig struct {
	MaxTurns int          `json:"max_turns" validate:"required,min=1"`
	Metadata util.JSONMap `json:"metadata" validate:"omitempty"`
}

type Execution struct {
	ErrorMessage string          `json:"error_message" validate:"omitempty"`
	Config       ExecutionConfig `json:"config" validate:"required"`
	Status       ExecutionStatus `json:"status" validate:"required"`
	Messages     []msg.Message   `json:"messages" validate:"required,min=1"`
	Turns        int             `json:"turns" validate:"min=0"`
}

func NewExecution(options ...ExecutionOption) (Execution, error) {
	execution := Execution{
		Config: ExecutionConfig{
			MaxTurns: 10,
			Metadata: util.JSONMap{},
		},
		Status:   ExecutionStatusPending,
		Messages: []msg.Message{},
		Turns:    0,
	}

	for _, option := range options {
		option(&execution)
	}

	for i := range execution.Messages {
		if execution.Messages[i].ID == "" {
			execution.Messages[i].ID = uuid.NewString()
		}
		execution.Messages[i].Sequence = i
	}

	if err := validate.Struct(&execution); err != nil {
		return execution, err
	}

	return execution, nil
}

func (e *Execution) AddMessage(msg msg.Message) msg.Message {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	msg.Sequence = len(e.Messages)
	e.Messages = append(e.Messages, msg)
	return msg
}
