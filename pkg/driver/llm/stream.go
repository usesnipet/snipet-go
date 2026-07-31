package llm

type StreamDelta struct {
	Content string
	Done    bool
	Err     error
}
