package llm

type StreamEvent interface {
	isStreamEvent()
}
type streamEvent struct{}

func (streamEvent) isStreamEvent() {}

type ErrorEvent struct {
	streamEvent
	Error string `json:"error"`
}

type TextDeltaEvent struct {
	streamEvent
	Text string `json:"text"`
}

type ToolCallStartedEvent struct {
	streamEvent
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ToolCallArgumentsDeltaEvent struct {
	streamEvent
	ID    string
	Delta string
}

type ToolCallFinishedEvent struct {
	streamEvent
	ID string
}

type CompletedEvent struct {
	streamEvent
}
