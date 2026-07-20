package runtime

import "github.com/usesnipet/snipet/internal/runtime/transport"

type IEvent interface {
	isEvent()
}

type ExecutionUpdatedEvent struct {
	Execution Execution
}

func (e ExecutionUpdatedEvent) isEvent() {}

type ExecutionMessageAddedEvent struct {
	Execution Execution
	Messages  []transport.Message
}

func (e ExecutionMessageAddedEvent) isEvent() {}

type ExecutionErrorEvent struct {
	Execution    Execution
	ErrorMessage string
}

func (e ExecutionErrorEvent) isEvent() {}
