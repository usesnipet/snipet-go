package runtime

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

// withToolCallFallback rewrites options for a model with no native tool
// call support: the toolset is described in the system prompt instead of
// sent as structured Tools, and the model is instructed to answer with a
// JSON envelope instead of a native tool call.
func withToolCallFallback(options llm.GenerateOptions) llm.GenerateOptions {
	instructions := fallbackInstructions(options.Tools)
	if options.Prompt.System != "" {
		options.Prompt.System = options.Prompt.System + "\n\n" + instructions
	} else {
		options.Prompt.System = instructions
	}
	options.Tools = tool.Toolset{}
	return options
}

// fallbackInstructions describes toolset as plain-text instructions for a
// model that can't be sent a structured tool list.
func fallbackInstructions(toolset tool.Toolset) string {
	var b strings.Builder
	b.WriteString("You can call the following tools. To call one, respond with ONLY a single JSON object of the exact form {\"tool_call\": {\"name\": \"<tool name>\", \"arguments\": { ... }}} and nothing else - no prose, no code fences. If no tool call is needed, respond normally in plain text.\n\nAvailable tools:\n")
	for _, t := range toolset.Tools {
		params, _ := json.Marshal(t.Parameters)
		fmt.Fprintf(&b, "- %s: %s\n  parameters schema: %s\n", t.Name, t.Description, params)
	}
	return b.String()
}

// fallbackToolCall is the JSON envelope a model without native tool call
// support is instructed to respond with (see fallbackInstructions).
type fallbackToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// parseFallbackToolCall extracts a tool call from a fallback model's raw
// text response, tolerating a surrounding ```json code fence. It returns
// nil when raw isn't a recognizable tool-call envelope, meaning the model
// answered normally.
func parseFallbackToolCall(raw string) *fallbackToolCall {
	trimmed := strings.TrimSpace(raw)
	for _, fence := range []string{"```json", "```"} {
		if rest, ok := strings.CutPrefix(trimmed, fence); ok {
			trimmed = strings.TrimSpace(strings.TrimSuffix(rest, "```"))
			break
		}
	}

	var envelope struct {
		ToolCall *fallbackToolCall `json:"tool_call"`
	}
	if err := json.Unmarshal([]byte(trimmed), &envelope); err != nil || envelope.ToolCall == nil || envelope.ToolCall.Name == "" {
		return nil
	}
	return envelope.ToolCall
}
