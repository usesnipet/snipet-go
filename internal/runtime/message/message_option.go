package message

import (
	"time"

	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type MessageOption func(message *Message)

func WithToolCalls(toolCalls []tool.Call) MessageOption {
	return func(message *Message) {
		message.ToolCalls = toolCalls
	}
}

func WithToolResult(toolResult *tool.Result) MessageOption {
	return func(message *Message) {
		message.ToolResult = toolResult
	}
}

func WithTimestamp(timestamp time.Time) MessageOption {
	return func(message *Message) {
		message.Timestamp = timestamp
	}
}

func WithID(id string) MessageOption {
	return func(message *Message) {
		message.ID = id
	}
}

func WithSequence(sequence int) MessageOption {
	return func(message *Message) {
		message.Sequence = sequence
	}
}
