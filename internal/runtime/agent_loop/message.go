package agentloop

import "time"

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleTool      MessageRole = "tool"
	MessageRoleFinal     MessageRole = "final"
)

type Message struct {
	ID         string
	Role       MessageRole
	Content    string
	ToolCalls  []ToolCall  // present on assistant messages that request tools
	ToolResult *ToolResult // present on role=tool messages
	Timestamp  time.Time
}
