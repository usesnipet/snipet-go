package agentloop

type Tool struct {
	Key          string // unique identifier for the tool
	Name         string // human-readable name for the tool
	Description  string // human-readable description for the tool
	ParamsSchema string // JSON schema for the parameters of the tool
	ResultSchema string // JSON schema for the result of the tool
}

type ToolCall struct {
	ID    string // unique id for this call, used to correlate with ToolResult
	Key   string // tool key to execute
	Input any
}

type ToolResult struct {
	ToolCallID string
	Key        string
	Success    bool
	Output     any
	Error      error
}
