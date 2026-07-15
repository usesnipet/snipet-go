package agentloop

import "context"

// Persister is the port for persisting an Execution and its messages.
// Concrete implementations live outside this package (e.g. gorm repositories).
type Persister interface {
	SaveExecution(ctx context.Context, e *Execution) error
	AppendMessage(ctx context.Context, executionID string, m *Message) error
}

// NoopPersister is the default in-memory/no-op implementation for tests and early use.
type NoopPersister struct{}

func NewNoopPersister() *NoopPersister {
	return &NoopPersister{}
}

func (NoopPersister) SaveExecution(context.Context, *Execution) error { return nil }

func (NoopPersister) AppendMessage(context.Context, string, *Message) error { return nil }
