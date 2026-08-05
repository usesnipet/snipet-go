package subscriber

import (
	"context"
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/runtime/execution"
)

type SSEEvent string

const (
	SSEEventExecutionTurnCompleted   SSEEvent = "execution_turn_completed"
	SSEEventExecutionMessageAdded    SSEEvent = "execution_message_added"
	SSEEventExecutionStatusChanged   SSEEvent = "execution_status_changed"
	SSEEventExecutionFinished        SSEEvent = "execution_finished"
	SSEEventExecutionMessageDelta    SSEEvent = "execution_message_delta"
	SSEEventExecutionToolCallStarted SSEEvent = "execution_tool_call_started"
	SSEEventExecutionToolResult      SSEEvent = "execution_tool_result"
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
	case execution.MessageAddedEvent:
		return s.sse.Write(string(SSEEventExecutionMessageAdded), event)
	case execution.TurnCompletedEvent:
		return s.sse.Write(string(SSEEventExecutionTurnCompleted), event)
	case execution.StatusChangedEvent:
		return s.sse.Write(string(SSEEventExecutionStatusChanged), event)
	case execution.FinishedEvent:
		return s.sse.Write(string(SSEEventExecutionFinished), map[string]string{"status": "done"})
	case execution.MessageDeltaEvent:
		return s.sse.Write(string(SSEEventExecutionMessageDelta), event)
	case execution.ToolCallStartedEvent:
		return s.sse.Write(string(SSEEventExecutionToolCallStarted), event)
	case execution.ToolResultEvent:
		return s.sse.Write(string(SSEEventExecutionToolResult), event)
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
