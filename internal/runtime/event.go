package runtime

import "github.com/usesnipet/snipet/pkg/msg"

type EventListener func(event IEvent) error

type IEvent interface {
	isEvent()
}

// ExecutionStatusChangedEvent is emitted whenever execution status (and related fields) change.
type ExecutionStatusChangedEvent struct {
	Status       ExecutionStatus `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

func (e ExecutionStatusChangedEvent) isEvent() {}

// ExecutionMessageAddedEvent is emitted when one or more messages are appended.
type ExecutionMessageAddedEvent struct {
	Message msg.Message `json:"message"`
}

func (e ExecutionMessageAddedEvent) isEvent() {}

// ExecutionTurnCompletedEvent is emitted when a turn is completed.
type ExecutionTurnCompletedEvent struct {
	Turn int `json:"turn"`
}

func (e ExecutionTurnCompletedEvent) isEvent() {}

// ExecutionFinishedEvent is emitted when an execution is finished.
type ExecutionFinishedEvent struct {
}

func (e ExecutionFinishedEvent) isEvent() {}

// ExecutionMessageDeltaEvent is emitted for each chunk of assistant text
// streamed from the LLM, before the full message is added via
// ExecutionMessageAddedEvent.
type ExecutionMessageDeltaEvent struct {
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
}

func (e ExecutionMessageDeltaEvent) isEvent() {}

// ExecutionToolCallStartedEvent is emitted when the LLM begins requesting a tool call.
type ExecutionToolCallStartedEvent struct {
	MessageID string `json:"message_id"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

func (e ExecutionToolCallStartedEvent) isEvent() {}

// ExecutionToolCallDeltaEvent is emitted for each chunk of a tool call's
// streamed arguments.
type ExecutionToolCallDeltaEvent struct {
	ID    string `json:"id"`
	Delta string `json:"delta"`
}

func (e ExecutionToolCallDeltaEvent) isEvent() {}

// ExecutionToolCallCompletedEvent is emitted once a tool call's arguments
// have finished streaming and been parsed, right before it is executed.
type ExecutionToolCallCompletedEvent struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

func (e ExecutionToolCallCompletedEvent) isEvent() {}

// ExecutionToolResultEvent is emitted after a tool call has been executed,
// carrying its result or the error that occurred.
type ExecutionToolResultEvent struct {
	ToolCallID string `json:"tool_call_id"`
	Tool       string `json:"tool"`
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
}

func (e ExecutionToolResultEvent) isEvent() {}
