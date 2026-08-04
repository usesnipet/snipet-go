package runtime

import "github.com/usesnipet/snipet/pkg/msg"

type EventListener func(event IEvent) error

type IEvent interface {
	isEvent()
}

// event is embedded in every event struct so they satisfy IEvent without
// each having to declare its own isEvent method.
type event struct{}

func (event) isEvent() {}

// ExecutionStatusChangedEvent is emitted whenever execution status (and related fields) change.
type ExecutionStatusChangedEvent struct {
	event
	Status       ExecutionStatus `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

// ExecutionMessageAddedEvent is emitted when one or more messages are appended.
type ExecutionMessageAddedEvent struct {
	event
	Message msg.Message `json:"message"`
}

// ExecutionTurnCompletedEvent is emitted when a turn is completed.
type ExecutionTurnCompletedEvent struct {
	event
	Turn int `json:"turn"`
}

// ExecutionFinishedEvent is emitted when an execution is finished.
type ExecutionFinishedEvent struct {
	event
}

// ExecutionMessageDeltaEvent is emitted for each chunk of assistant text
// streamed from the LLM, before the full message is added via
// ExecutionMessageAddedEvent.
type ExecutionMessageDeltaEvent struct {
	event
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
}

// ExecutionToolCallEvent is emitted once a tool call's arguments have
// finished streaming and been parsed, right before it is executed.
type ExecutionToolCallEvent struct {
	event
	MessageID string         `json:"message_id"`
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ExecutionAttemptFailedEvent is emitted when an LLM generation attempt
// fails after part of its response was already streamed to subscribers
// (text deltas and/or tool call events). It signals that everything
// published under MessageID for this attempt must be discarded — the engine
// will either retry with a different LLM configuration or fail the
// execution, but this attempt's partial content never becomes part of the
// conversation history.
type ExecutionAttemptFailedEvent struct {
	event
	MessageID string `json:"message_id"`
	Error     string `json:"error"`
}

// ExecutionToolResultEvent is emitted after a tool call has been executed,
// carrying its result or the error that occurred.
type ExecutionToolResultEvent struct {
	event
	ToolCallID string `json:"tool_call_id"`
	Tool       string `json:"tool"`
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
}
