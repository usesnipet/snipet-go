package llm

// StreamEvent is a sealed interface implemented by every event Driver.Stream
// may emit on its channel. Consumers type-switch on the concrete event types
// below (ErrorEvent, TextDeltaEvent, ToolCallEvent, ToolCallErrorEvent,
// CompletedEvent); no other package can implement StreamEvent.
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

// ToolCallEvent signals that the model finished requesting a tool call: its
// ID, Name and Arguments are fully known. The driver is responsible for
// assembling this internally from whatever incremental shape its wire
// protocol uses.
type ToolCallEvent struct {
	streamEvent
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolCallErrorEvent signals that the model requested a tool call whose
// arguments could not be assembled (e.g. malformed JSON). It does not end
// the stream; consumers should drop this one call and keep processing.
type ToolCallErrorEvent struct {
	streamEvent
	ID    string `json:"id"`
	Error string `json:"error"`
}

// CompletedEvent signals that the stream finished successfully. It is
// normally the last event received before the channel is closed.
type CompletedEvent struct {
	streamEvent
}
