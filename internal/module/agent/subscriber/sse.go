package subscriber

import (
	"context"
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/runtime/execution"
)

type SSEEvent string

const (
	// Execution Events
	SSEEventExecutionStarted       SSEEvent = "execution.started"
	SSEEventExecutionStatusChanged SSEEvent = "execution.status_changed"
	SSEEventExecutionFinished      SSEEvent = "execution.finished"

	// Turn Events
	SSEEventTurnStarted   SSEEvent = "turn.started"
	SSEEventTurnCompleted SSEEvent = "turn.completed"

	// Message Events
	SSEEventMessageAdded         SSEEvent = "message.added"
	SSEEventMessageDelta         SSEEvent = "message.delta"
	SSEEventMessageAttemptFailed SSEEvent = "message.attempt_failed"

	// Tool Events
	SSEEventToolCallBegin   SSEEvent = "tool_call.begin"
	SSEEventToolCallStarted SSEEvent = "tool_call.started"
	SSEEventToolCallResult  SSEEvent = "tool_call.result"
)

type SSE struct {
	w   http.ResponseWriter
	sse *api.SSEWriter
}

func NewSSE(w http.ResponseWriter) *SSE {
	return &SSE{w: w}
}

func (s *SSE) Handle(ctx context.Context, event execution.IEvent) error {
	if err := s.ensureSSE(); err != nil {
		return err
	}

	switch event := event.(type) {
	case execution.StartedEvent:
		return s.sse.Write(string(SSEEventExecutionStarted), event)
	case execution.StatusChangedEvent:
		return s.sse.Write(string(SSEEventExecutionStatusChanged), event)
	case execution.FinishedEvent:
		return s.sse.Write(string(SSEEventExecutionFinished), map[string]string{"status": "done"})

	case execution.TurnStartedEvent:
		return s.sse.Write(string(SSEEventTurnStarted), event)
	case execution.TurnCompletedEvent:
		return s.sse.Write(string(SSEEventTurnCompleted), event)

	case execution.MessageAddedEvent:
		return s.sse.Write(string(SSEEventMessageAdded), event)
	case execution.MessageDeltaEvent:
		return s.sse.Write(string(SSEEventMessageDelta), event)
	case execution.MessageAttemptFailedEvent:
		return s.sse.Write(string(SSEEventMessageAttemptFailed), event)

	case execution.ToolCallBeginEvent:
		return s.sse.Write(string(SSEEventToolCallBegin), event)
	case execution.ToolCallStartedEvent:
		return s.sse.Write(string(SSEEventToolCallStarted), event)
	case execution.ToolCallResultEvent:
		return s.sse.Write(string(SSEEventToolCallResult), event)
	}
	return nil
}

func (s *SSE) Write(event string, data any) error {
	if err := s.ensureSSE(); err != nil {
		return err
	}
	return s.sse.Write(event, data)
}

func (s *SSE) HandleError(err error) error {
	if err := s.ensureSSE(); err != nil {
		return err
	}
	return s.sse.Write("error", map[string]string{"message": err.Error()})
}

func (s *SSE) ensureSSE() error {
	if s.sse != nil {
		return nil
	}
	var err error
	s.sse, err = api.NewSSEWriter(s.w)
	return err
}
