package llm

// CreateDriver builds a Driver from the given Options. The actual behavior
// (TestConnection/Generate/Stream) comes from WithAPI; any method whose
// underlying func is nil returns an error when called.
func CreateDriver(opts ...Option) Driver {
	d := &llmDriver{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}
