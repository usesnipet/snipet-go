package runtime

import (
	"github.com/usesnipet/snipet/internal/runtime/transport"
)

type ExecutionOption func(execution *Execution)

func WithMaxTurns(maxTurns int) ExecutionOption {
	return func(execution *Execution) {
		execution.Config.MaxTurns = maxTurns
	}
}

func WithMetadata(key string, value any) ExecutionOption {
	return func(execution *Execution) {
		execution.Config.Metadata[key] = value
	}
}

func WithInitialMessages(message transport.Message, messages ...transport.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, message)
		execution.Messages = append(execution.Messages, messages...)
	}
}

func WithMessage(message transport.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, message)
	}
}

func WithMessageFromUser(content string, options ...transport.MessageOption) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(
			execution.Messages,
			transport.NewMessage(transport.MessageRoleUser, content, options...),
		)
	}
}
