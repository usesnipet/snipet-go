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
