package agentloop

import "time"

type ExecutionStatus string

const (
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusMaxTurns  ExecutionStatus = "max_turns"
)

type Execution struct {
	ID        string
	AgentID   string
	SessionID string
	Status    ExecutionStatus
	Messages  []*Message
	Turns     int
	MaxTurns  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewExecution(agentID string) *Execution {
	now := time.Now()
	return &Execution{
		AgentID:   agentID,
		Status:    ExecutionStatusRunning,
		Messages:  []*Message{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (e *Execution) AddMessage(message *Message) {
	e.Messages = append(e.Messages, message)
	e.UpdatedAt = time.Now()
}
