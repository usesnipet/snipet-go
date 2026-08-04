package subscriber

import (
	"context"
	"net/http"

	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/runtime"
)

type SSEEvent string

const (
	SSEEventExecutionTurnCompleted SSEEvent = "execution_turn_completed"
	SSEEventExecutionMessageAdded  SSEEvent = "execution_message_added"
	SSEEventExecutionStatusChanged SSEEvent = "execution_status_changed"
	SSEEventExecutionFinished      SSEEvent = "execution_finished"
	SSEEventExecutionMessageDelta  SSEEvent = "execution_message_delta"
	SSEEventExecutionToolCall      SSEEvent = "execution_tool_call"
	SSEEventExecutionToolResult    SSEEvent = "execution_tool_result"
)

type SSE struct {
	w   http.ResponseWriter
	sse *api.SSEWriter
}

func NewSSE(w http.ResponseWriter) *SSE {
	return &SSE{w: w}
}

func (s *SSE) Handle(ctx context.Context, event runtime.IEvent) error {
	if err := s.ensureSSE(); err != nil {
		return err
	}

	switch event := event.(type) {
	case runtime.ExecutionMessageAddedEvent:
		return s.sse.Write(string(SSEEventExecutionMessageAdded), event)
	case runtime.ExecutionTurnCompletedEvent:
		return s.sse.Write(string(SSEEventExecutionTurnCompleted), event)
	case runtime.ExecutionStatusChangedEvent:
		return s.sse.Write(string(SSEEventExecutionStatusChanged), event)
	case runtime.ExecutionFinishedEvent:
		return s.sse.Write(string(SSEEventExecutionFinished), map[string]string{"status": "done"})
	case runtime.ExecutionMessageDeltaEvent:
		return s.sse.Write(string(SSEEventExecutionMessageDelta), event)
	case runtime.ExecutionToolCallEvent:
		return s.sse.Write(string(SSEEventExecutionToolCall), event)
	case runtime.ExecutionToolResultEvent:
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
