package runtime

import (
	"github.com/usesnipet/snipet/pkg/msg"
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

func WithInitialMessages(messages ...msg.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, messages...)
	}
}

func WithMessage(msg msg.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, msg)
	}
}

func WithMessageFromUser(content string, options ...msg.MessageOption) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(
			execution.Messages,
			msg.NewMessage(msg.RoleUser, content, options...),
		)
	}
}
