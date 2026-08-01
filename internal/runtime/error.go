package runtime

import "errors"

var (
	ErrNoLLMConfigured     = errors.New("no llm configured")
	ErrLLMGenerationFailed = errors.New("llm generation failed")

	ErrToolNotFound = errors.New("tool not found")

	ErrFinishExecution = errors.New("execution finished")
)
