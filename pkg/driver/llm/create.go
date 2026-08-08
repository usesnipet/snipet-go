package llm

// CreateDriver builds a Driver from the given Options. Key, TestConnection,
// Stream, and Generate (the latter three set via WithAPI) are required;
// CreateDriver returns an error instead of a Driver if any of them is
// missing, so a misconfigured driver never gets registered. WithModelLoader
// is optional.
func CreateDriver(opts ...Option) (Driver, error) {
	d := &llmDriver{}
	for _, opt := range opts {
		opt(d)
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}
