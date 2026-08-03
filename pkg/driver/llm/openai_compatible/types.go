package openaicompatible

// chatRequest is the OpenAI-compatible /v1/chat/completions request body.
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Tools       []chatTool    `json:"tools,omitempty"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float64      `json:"top_p,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// chatMessage represents both a request message and a response
// message/delta: Content/ToolCalls/ToolCallID are used when sending, while
// the same fields (via chatChoice.Message or .Delta) hold response data.
type chatMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content,omitempty"`
	ToolCalls  []chatToolCall `json:"tool_calls,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
}

// chatTool is a single function tool advertised in a chat request.
type chatTool struct {
	Type     string           `json:"type"`
	Function chatToolFunction `json:"function"`
}

type chatToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// chatToolCall is a tool call as sent in a request (replaying a prior
// assistant turn) or received in a response/delta. Index identifies which
// in-progress tool call a streamed delta belongs to.
type chatToolCall struct {
	Index    *int                 `json:"index,omitempty"`
	ID       string               `json:"id,omitempty"`
	Type     string               `json:"type,omitempty"`
	Function chatToolCallFunction `json:"function"`
}

type chatToolCallFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// chatResponse is the shape of both a full chat completions response and
// each streamed SSE chunk; Error is populated instead of Choices when the
// provider reports a failure inline (with a 2xx status).
type chatResponse struct {
	Choices []chatChoice `json:"choices"`
	Error   *apiError    `json:"error,omitempty"`
}

// chatChoice is one completion candidate: Message on a non-streaming
// response, Delta on a streamed chunk.
type chatChoice struct {
	Index        int          `json:"index"`
	Message      *chatMessage `json:"message,omitempty"`
	Delta        *chatMessage `json:"delta,omitempty"`
	FinishReason *string      `json:"finish_reason"`
}

// apiError is an OpenAI-compatible inline API error payload.
type apiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code"`
}

// Error implements the error interface, formatting as "type: message" when
// Type is set.
func (e *apiError) Error() string {
	if e == nil {
		return "unknown api error"
	}
	if e.Type != "" {
		return e.Type + ": " + e.Message
	}
	return e.Message
}
