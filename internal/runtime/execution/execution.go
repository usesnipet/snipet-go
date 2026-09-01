package execution

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/usesnipet/snipet/pkg/jsonx"
	"github.com/usesnipet/snipet/pkg/msg"
)

var validate = validator.New()

type Config struct {
	MaxTurns int           `json:"max_turns" validate:"omitempty"`
	Metadata jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type Execution struct {
	Agent          *Agent
	publisher      IPublisher
	Config         Config
	StreamMessages bool

	ErrorMessage string        `json:"error_message,omitempty" validate:"omitempty,max=255"`
	Status       Status        `json:"status" validate:"required,oneof=pending running completed failed max_turns cancelled"`
	Messages     []msg.Message `json:"messages" validate:"required"`
	Turns        int           `json:"turns" validate:"min=0"`
}

func NewExecution(options ...ExecutionOption) (*Execution, error) {
	execution := &Execution{
		Agent:     nil,
		publisher: NewLocalPublisher(),
		Config:    Config{MaxTurns: 10, Metadata: jsonx.JSONMap{}},
		Status:    StatusPending,
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

func (e *Execution) Start(ctx context.Context) error {
	return e.Publish(ctx, StartedEvent{})
}

func (e *Execution) AddMessage(ctx context.Context, msg msg.Message) error {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	msg.Sequence = len(e.Messages)
	e.Messages = append(e.Messages, msg)
	return e.Publish(ctx, MessageAddedEvent{Message: msg})
}

func (e *Execution) StartTurn(ctx context.Context) error {
	e.Turns++
	return e.Publish(ctx, TurnStartedEvent{Turn: e.Turns})
}

func (e *Execution) CompleteTurn(ctx context.Context) error {
	return e.Publish(ctx, TurnCompletedEvent{Turn: e.Turns})
}

func (e *Execution) Finish(ctx context.Context) error {
	e.Status = StatusCompleted
	return e.Publish(ctx, FinishedEvent{})
}

func (e *Execution) SetStatus(ctx context.Context, status Status) error {
	e.Status = status
	return e.Publish(ctx, StatusChangedEvent{
		Status: status,
	})
}

func (e *Execution) SetError(ctx context.Context, errorMessage string) error {
	e.ErrorMessage = errorMessage
	e.Status = StatusFailed
	return e.Publish(ctx, StatusChangedEvent{Status: e.Status, ErrorMessage: e.ErrorMessage})
}

func (e *Execution) SetMaxTurnsReachedError(ctx context.Context) error {
	e.ErrorMessage = "Maximum number of turns reached"
	e.Status = StatusMaxTurns
	return e.Publish(ctx, StatusChangedEvent{Status: e.Status, ErrorMessage: e.ErrorMessage})
}

func (e *Execution) Cancel(ctx context.Context) error {
	return e.SetStatus(ctx, StatusCancelled)
}

func (e *Execution) Publish(ctx context.Context, event IEvent) error {
	return e.publisher.Publish(ctx, event)
}

func (e *Execution) Subscribe(subscribers ...Subscriber) {
	e.publisher.Subscribe(subscribers...)
}
