package runtime

import "github.com/usesnipet/snipet/internal/runtime/message"

type EventListener func(event IEvent) error

type IEvent interface {
	isEvent()
}

// ExecutionStatusChangedEvent is emitted whenever execution status (and related fields) change.
type ExecutionStatusChangedEvent struct {
	Status       ExecutionStatus `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
	Turns        int             `json:"turns"`
}

func (e ExecutionStatusChangedEvent) isEvent() {}

// ExecutionMessageAddedEvent is emitted when one or more messages are appended.
type ExecutionMessageAddedEvent struct {
	Messages []message.Message `json:"messages"`
}

func (e ExecutionMessageAddedEvent) isEvent() {}
