package execution

type Status string

const (
	StatusPending   Status = "pending"   // the execution is pending to start
	StatusRunning   Status = "running"   // the execution is running
	StatusCompleted Status = "completed" // the execution is completed
	StatusFailed    Status = "failed"    // the execution is failed
	StatusMaxTurns  Status = "max_turns" // the execution has reached the maximum number of turns
	StatusCancelled Status = "cancelled" // the execution was cancelled via context
)
