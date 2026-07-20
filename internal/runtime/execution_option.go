package runtime

import "github.com/usesnipet/snipet/internal/runtime/message"

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

func WithInitialMessages(msg message.Message, messages ...message.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, msg)
		execution.Messages = append(execution.Messages, messages...)
	}
}

func WithMessage(msg message.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, msg)
	}
}

func WithMessageFromUser(content string, options ...message.MessageOption) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(
			execution.Messages,
			message.NewMessage(message.MessageRoleUser, content, options...),
		)
	}
}
