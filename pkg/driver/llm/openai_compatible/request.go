package openaicompatible

import (
	"encoding/json"

	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// buildChatRequest translates a GenerateOptions/Config pair into the
// OpenAI-compatible request body, omitting optional numeric parameters that
// were left at their zero value.
func buildChatRequest(cfg Config, options llm.GenerateOptions, stream bool) chatRequest {
	req := chatRequest{
		Model:    cfg.Model,
		Messages: buildMessages(options.Prompt),
		Tools:    buildTools(options.Tools),
		Stream:   stream,
	}

	if cfg.MaxTokens != 0 {
		m := cfg.MaxTokens
		req.MaxTokens = &m
	}

	req.Temperature = &cfg.Temperature
	req.TopP = &cfg.TopP

	return req
}

// buildMessages converts a llm.Prompt into the OpenAI-compatible message
// list, prepending a system message when Prompt.System is set and dropping
// any message whose Role has no OpenAI equivalent (see toOpenAIRole).
func buildMessages(prompt llm.Prompt) []chatMessage {
	messages := make([]chatMessage, 0, len(prompt.Messages)+1)
	if prompt.System != "" {
		messages = append(messages, chatMessage{
			Role:    "system",
			Content: prompt.System,
		})
	}
	for _, m := range prompt.Messages {
		role, ok := toOpenAIRole(m.Role)
		if !ok {
			continue
		}
		message := chatMessage{
			Role:    role,
			Content: m.Content,
		}
		if m.Role == msg.RoleAssistant && len(m.ToolCalls) > 0 {
			message.ToolCalls = buildToolCalls(m.ToolCalls)
		}
		if m.Role == msg.RoleTool {
			message.ToolCallID = m.ToolCallID
		}
		messages = append(messages, message)
	}
	return messages
}

// buildToolCalls converts assistant tool.Call records into the
// OpenAI-compatible tool_calls representation, JSON-encoding each call's
// arguments (falling back to "{}" if they don't marshal).
func buildToolCalls(calls []tool.Call) []chatToolCall {
	toolCalls := make([]chatToolCall, 0, len(calls))
	for _, c := range calls {
		arguments, err := json.Marshal(c.Arguments)
		if err != nil {
			arguments = []byte("{}")
		}
		toolCalls = append(toolCalls, chatToolCall{
			ID:   c.ID,
			Type: "function",
			Function: chatToolCallFunction{
				Name:      c.Tool,
				Arguments: string(arguments),
			},
		})
	}
	return toolCalls
}

// buildTools converts a tool.Toolset into the OpenAI-compatible tools list,
// returning nil when the toolset is empty and defaulting each tool's
// parameters to an empty object schema when not specified.
func buildTools(toolset tool.Toolset) []chatTool {
	if len(toolset.Tools) == 0 {
		return nil
	}
	tools := make([]chatTool, 0, len(toolset.Tools))
	for _, t := range toolset.Tools {
		params := t.Parameters
		if params == nil {
			params = map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			}
		}
		tools = append(tools, chatTool{
			Type: "function",
			Function: chatToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}
	return tools
}

// toOpenAIRole maps a msg.Role to its OpenAI-compatible string value. The
// second return value is false for roles with no OpenAI equivalent.
func toOpenAIRole(role msg.Role) (string, bool) {
	switch role {
	case msg.RoleSystem:
		return "system", true
	case msg.RoleUser:
		return "user", true
	case msg.RoleAssistant:
		return "assistant", true
	case msg.RoleTool:
		return "tool", true
	default:
		return "", false
	}
}
