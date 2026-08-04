package execution

import (
	"context"

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
	MaxTurns int          `json:"max_turns" validate:"omitempty"`
	Metadata util.JSONMap `json:"metadata" validate:"omitempty"`
}

type Execution struct {
	agent     *Agent
	publisher IPublisher
	Config    ExecutionConfig

	ErrorMessage string          `json:"error_message,omitempty" validate:"omitempty,max=255"`
	Status       ExecutionStatus `json:"status" validate:"required,oneof=pending running completed failed max_turns cancelled"`
	Messages     []msg.Message   `json:"messages" validate:"required"`
	Turns        int             `json:"turns" validate:"min=0"`
}

func NewExecution(options ...ExecutionOption) (*Execution, error) {
	execution := &Execution{
		agent:     nil,
		publisher: NewLocalPublisher(),
		Config:    ExecutionConfig{MaxTurns: 10, Metadata: util.JSONMap{}},
		Status:    ExecutionStatusPending,
		Messages:  []msg.Message{},
		Turns:     0,
	}

	for _, option := range options {
		option(execution)
	}

	for i := range execution.Messages {
		if execution.Messages[i].ID == "" {
			execution.Messages[i].ID = uuid.NewString()
		}
		execution.Messages[i].Sequence = i
	}

	if err := validate.Struct(execution); err != nil {
		return execution, err
	}

	return execution, nil
}

func (e *Execution) AddMessage(ctx context.Context, msg msg.Message) error {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	msg.Sequence = len(e.Messages)
	e.Messages = append(e.Messages, msg)
	return e.publish(ctx, ExecutionMessageAddedEvent{Message: msg})
}

func (e *Execution) CompleteTurn(ctx context.Context) error {
	e.Turns++
	return e.publish(ctx, ExecutionTurnCompletedEvent{Turn: e.Turns})
}

func (e *Execution) Finish(ctx context.Context) error {
	e.Status = ExecutionStatusCompleted
	return e.publish(ctx, ExecutionFinishedEvent{})
}

func (e *Execution) SetStatus(ctx context.Context, status ExecutionStatus) error {
	e.Status = status
	return e.publish(ctx, ExecutionStatusChangedEvent{
		Status: status,
	})
}

func (e *Execution) SetError(ctx context.Context, errorMessage string) error {
	e.ErrorMessage = errorMessage
	e.Status = ExecutionStatusFailed
	return e.publish(ctx, ExecutionStatusChangedEvent{Status: e.Status, ErrorMessage: e.ErrorMessage})
}

func (e *Execution) SetMaxTurnsReachedError(ctx context.Context) error {
	e.ErrorMessage = "Maximum number of turns reached"
	e.Status = ExecutionStatusMaxTurns
	return e.publish(ctx, ExecutionStatusChangedEvent{Status: e.Status, ErrorMessage: e.ErrorMessage})
}

func (e *Execution) Cancel(ctx context.Context) error {
	return e.SetStatus(ctx, ExecutionStatusCancelled)
}

func (e *Execution) publish(ctx context.Context, event IEvent) error {
	return e.publisher.Publish(ctx, event)
}

func (e *Execution) Subscribe(subscribers ...Subscriber) {
	e.publisher.Subscribe(subscribers...)
}
