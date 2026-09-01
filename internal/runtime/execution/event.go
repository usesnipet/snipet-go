package execution

import (
	"github.com/usesnipet/snipet/pkg/msg"
)

type EventListener func(event IEvent) error

type IEvent interface {
	isEvent()
}

// event is embedded in every event struct so they satisfy IEvent without
// each having to declare its own isEvent method.
type event struct{}

func (event) isEvent() {}

// #region ExecutionEvents

// StartedEvent is emitted when an execution is started.
type StartedEvent struct {
	event
}

// StatusChangedEvent is emitted whenever execution status (and related fields) change.
type StatusChangedEvent struct {
	event
	Status       Status `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// FinishedEvent is emitted when an execution is finished.
type FinishedEvent struct {
	event
}

// #endregion

// #region TurnEvents

// TurnStartedEvent is emitted when a turn is started.
type TurnStartedEvent struct {
	event
	Turn int `json:"turn"`
}

// TurnCompletedEvent is emitted when a turn is completed.
type TurnCompletedEvent struct {
	event
	Turn int `json:"turn"`
}

// #endregion

// #region MessageEvents

// MessageAddedEvent is emitted when one or more messages are appended.
type MessageAddedEvent struct {
	event
	Message msg.Message `json:"message"`
}

// MessageDeltaEvent is emitted for each chunk of assistant text
// streamed from the LLM, before the full message is added via
// ExecutionMessageAddedEvent.
type MessageDeltaEvent struct {
	event
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
}

// MessageAttemptFailedEvent is emitted when an LLM generation attempt
// fails after part of its response was already streamed to subscribers
// (text deltas and/or tool call events). It signals that everything
// published under MessageID for this attempt must be discarded — the engine
// will either retry with a different LLM configuration or fail the
// execution, but this attempt's partial content never becomes part of the
// conversation history.
type MessageAttemptFailedEvent struct {
	event
	MessageID string `json:"message_id"`
	Error     string `json:"error"`
}

// #endregion

// #region ToolEvents

// ToolCallBeginEvent is emitted when the LLM start to request a tool call.
type ToolCallBeginEvent struct {
	event
}

// ToolCallStartedEvent is emitted right before a requested tool
// call is invoked.
type ToolCallStartedEvent struct {
	event
	ToolCallID string         `json:"tool_call_id"`
	Tool       string         `json:"tool"`
	Arguments  map[string]any `json:"arguments"`
}

// ToolCallResultEvent is emitted after a tool call has been executed,
// carrying its result or the error that occurred.
type ToolCallResultEvent struct {
	event
	ToolCallID string `json:"tool_call_id"`
	Tool       string `json:"tool"`
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
}

// #endregion
