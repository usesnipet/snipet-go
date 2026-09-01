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

// transition sets the execution's status and error message and publishes the
// StatusChangedEvent that every status change goes through. Callers that
// represent a terminal outcome publish their specific lifecycle event
// afterwards (see Finish, SetError, SetMaxTurnsReachedError, Cancel).
func (e *Execution) transition(ctx context.Context, status Status, errorMessage string) error {
	e.Status = status
	e.ErrorMessage = errorMessage
	return e.Publish(ctx, StatusChangedEvent{Status: status, ErrorMessage: errorMessage})
}

// SetStatus performs a non-terminal status change (e.g. to StatusRunning). It
// clears any error message; use SetError for failures.
func (e *Execution) SetStatus(ctx context.Context, status Status) error {
	return e.transition(ctx, status, "")
}

func (e *Execution) Finish(ctx context.Context) error {
	if err := e.transition(ctx, StatusCompleted, ""); err != nil {
		return err
	}
	return e.Publish(ctx, FinishedEvent{})
}

func (e *Execution) SetError(ctx context.Context, errorMessage string) error {
	if err := e.transition(ctx, StatusFailed, errorMessage); err != nil {
		return err
	}
	return e.Publish(ctx, FailedEvent{Error: errorMessage})
}

func (e *Execution) SetMaxTurnsReachedError(ctx context.Context) error {
	if err := e.transition(ctx, StatusMaxTurns, "Maximum number of turns reached"); err != nil {
		return err
	}
	return e.Publish(ctx, MaxTurnsReachedEvent{MaxTurns: e.Config.MaxTurns})
}

func (e *Execution) Cancel(ctx context.Context) error {
	if err := e.transition(ctx, StatusCancelled, ""); err != nil {
		return err
	}
	return e.Publish(ctx, CancelledEvent{})
}

func (e *Execution) Publish(ctx context.Context, event IEvent) error {
	return e.publisher.Publish(ctx, event)
}

func (e *Execution) Subscribe(subscribers ...Subscriber) {
	e.publisher.Subscribe(subscribers...)
}
