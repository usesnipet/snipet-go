package runtime

import "github.com/usesnipet/snipet/pkg/driver/llm"

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

func WithInitialMessages(msg llm.Message, messages ...llm.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, msg)
		execution.Messages = append(execution.Messages, messages...)
	}
}

func WithMessage(msg llm.Message) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(execution.Messages, msg)
	}
}

func WithMessageFromUser(content string, options ...llm.MessageOption) ExecutionOption {
	return func(execution *Execution) {
		execution.Messages = append(
			execution.Messages,
			llm.NewMessage(llm.MessageRoleUser, content, options...),
		)
	}
}
