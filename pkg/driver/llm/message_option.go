package llm

import (
	"time"
)

type MessageOption func(message *Message)

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
