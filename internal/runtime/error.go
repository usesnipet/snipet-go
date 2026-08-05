package runtime

import "errors"

var (
	ErrAgentNotConfigured = errors.New("agent not configured")
	ErrNoLLMConfigured    = errors.New("llm not configured")

	ErrLLMGenerationFailed = errors.New("llm generation failed")

	ErrToolNotFound = errors.New("tool not found")

	ErrFinishExecution = errors.New("execution finished")
)
