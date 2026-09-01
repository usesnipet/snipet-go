package runtime

import "errors"

var (
	ErrAgentNotConfigured = errors.New("agent not configured")
	ErrNoLLMConfigured    = errors.New("llm not configured")
)
