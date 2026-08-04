package runtime

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// withToolCallFallback rewrites options for a model with no native tool
// call support: the toolset is described in the system prompt instead of
// sent as structured Tools, and the model is instructed to answer with a
// JSON envelope instead of a native tool call. Past tool call/result turns
// in the history are also re-rendered as plain assistant/user text (see
// rewriteMessagesForFallback) — a chat template for a non-tool-calling
// model typically has no branch for the "tool" role or an assistant
// message's tool_calls field, so left as-is those turns would be silently
// dropped and the model would never see that a tool already ran, calling it
// again instead of answering from its result.
func withToolCallFallback(options llm.GenerateOptions) llm.GenerateOptions {
	instructions := fallbackInstructions(options.Tools)
	if options.Prompt.System != "" {
		options.Prompt.System = options.Prompt.System + "\n\n" + instructions
	} else {
		options.Prompt.System = instructions
	}
	options.Prompt.Messages = rewriteMessagesForFallback(options.Prompt.Messages)
	options.Tools = tool.Toolset{}
	return options
}

// rewriteMessagesForFallback re-renders tool call requests and their
// results as plain-text assistant/user turns instead of the structured
// tool_calls/"tool"-role shape, so a fallback model's chat template (which
// generally only knows system/user/assistant) can still follow along.
func rewriteMessagesForFallback(messages []msg.Message) []msg.Message {
	toolNameByCallID := map[string]string{}
	rewritten := make([]msg.Message, 0, len(messages))

	for _, m := range messages {
		switch {
		case m.Role == msg.RoleAssistant && len(m.ToolCalls) > 0:
			for _, call := range m.ToolCalls {
				toolNameByCallID[call.ID] = call.Tool
			}
			rewritten = append(rewritten, msg.NewMessage(msg.RoleAssistant, renderFallbackToolCalls(m.ToolCalls)))

		case m.Role == msg.RoleTool:
			name := toolNameByCallID[m.ToolCallID]
			if name == "" {
				name = "tool"
			}
			rewritten = append(rewritten, msg.NewMessage(msg.RoleUser, fmt.Sprintf("Result of calling %s:\n%s", name, m.Content)))

		default:
			rewritten = append(rewritten, m)
		}
	}
	return rewritten
}

// renderFallbackToolCalls reconstructs the JSON envelope text (see
// fallbackInstructions) a fallback model would have produced to request
// calls, so its own prior turn reads back exactly as instructed - keeping
// its "memory" of what it said consistent with the convention it was taught.
func renderFallbackToolCalls(calls []tool.Call) string {
	var b strings.Builder
	for i, call := range calls {
		if i > 0 {
			b.WriteString("\n")
		}
		envelope, err := json.Marshal(struct {
			ToolCall fallbackToolCall `json:"tool_call"`
		}{ToolCall: fallbackToolCall{Name: call.Tool, Arguments: call.Arguments}})
		if err != nil {
			continue
		}
		b.Write(envelope)
	}
	return b.String()
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
// text response. Despite being instructed to respond with ONLY the JSON
// envelope, models without native tool call support often don't follow that
// strictly (surrounding prose, a trailing comment, a ```json code fence), so
// this first tries the whole trimmed text and, failing that, scans it for a
// balanced JSON object that decodes as the envelope. It returns nil when no
// such object is found, meaning the model answered normally.
func parseFallbackToolCall(raw string) *fallbackToolCall {
	trimmed := strings.TrimSpace(raw)
	for _, fence := range []string{"```json", "```"} {
		if rest, ok := strings.CutPrefix(trimmed, fence); ok {
			trimmed = strings.TrimSpace(strings.TrimSuffix(rest, "```"))
			break
		}
	}

	if call := decodeFallbackToolCall(trimmed); call != nil {
		return call
	}
	for _, candidate := range balancedJSONObjects(trimmed) {
		if call := decodeFallbackToolCall(candidate); call != nil {
			return call
		}
	}
	return nil
}

// decodeFallbackToolCall decodes s as the {"tool_call": {...}} envelope,
// returning nil if it doesn't match (invalid JSON, missing tool_call, or a
// tool_call with no name).
func decodeFallbackToolCall(s string) *fallbackToolCall {
	var envelope struct {
		ToolCall *fallbackToolCall `json:"tool_call"`
	}
	if err := json.Unmarshal([]byte(s), &envelope); err != nil || envelope.ToolCall == nil || envelope.ToolCall.Name == "" {
		return nil
	}
	return envelope.ToolCall
}

// balancedJSONObjects returns every top-level "{...}" substring of s, in
// order of appearance, tracking brace depth while ignoring braces inside
// JSON string literals so a value like {"note": "a { b }"} isn't split up.
func balancedJSONObjects(s string) []string {
	var objects []string
	depth := 0
	start := -1
	inString := false
	escaped := false

	for i, r := range s {
		if inString {
			switch {
			case escaped:
				escaped = false
			case r == '\\':
				escaped = true
			case r == '"':
				inString = false
			}
			continue
		}
		switch r {
		case '"':
			inString = true
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			if depth > 0 {
				depth--
				if depth == 0 && start != -1 {
					objects = append(objects, s[start:i+1])
					start = -1
				}
			}
		}
	}
	return objects
}
