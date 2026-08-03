package llm

// StreamEvent is a sealed interface implemented by every event Driver.Stream
// may emit on its channel. Consumers type-switch on the concrete event types
// below (ErrorEvent, TextDeltaEvent, ToolCall*Event, CompletedEvent); no
// other package can implement StreamEvent.
type StreamEvent interface {
	isStreamEvent()
}
type streamEvent struct{}

func (streamEvent) isStreamEvent() {}

// ErrorEvent signals that the stream failed. It is normally the last event
// received before the channel is closed.
type ErrorEvent struct {
	streamEvent
	Error string `json:"error"`
}

// TextDeltaEvent carries an incremental chunk of assistant text output.
type TextDeltaEvent struct {
	streamEvent
	Text string `json:"text"`
}

// ToolCallStartedEvent signals that the model began requesting a tool call,
// identified by ID, before its arguments are known.
type ToolCallStartedEvent struct {
	streamEvent
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToolCallArgumentsDeltaEvent carries an incremental chunk of a tool call's
// JSON arguments; concatenate Delta values by ID to reconstruct the full
// argument payload.
type ToolCallArgumentsDeltaEvent struct {
	streamEvent
	ID    string
	Delta string
}

// ToolCallFinishedEvent signals that a tool call's arguments are complete
// and it is ready to be invoked.
type ToolCallFinishedEvent struct {
	streamEvent
	ID string
}

// CompletedEvent signals that the stream finished successfully. It is
// normally the last event received before the channel is closed.
type CompletedEvent struct {
	streamEvent
}
