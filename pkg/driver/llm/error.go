package llm

import "errors"

var (
	ErrTestConnectionNotConfigured = errors.New("test connection not configured")
	ErrModelLoaderNotConfigured    = errors.New("model loader not configured")
	ErrGenerateNotConfigured       = errors.New("generate not configured")
	ErrStreamNotConfigured         = errors.New("stream not configured")

	ErrModelNotFound = errors.New("model not found")
)
