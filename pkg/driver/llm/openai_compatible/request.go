package openaicompatible

import (
	"encoding/json"

	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

func buildChatRequest(cfg Config, options llm.GenerateOptions, stream bool) chatRequest {
	req := chatRequest{
		Model:    cfg.Model,
		Messages: buildMessages(options.Prompt),
		Tools:    buildTools(options.Tools),
		Stream:   stream,
	}
	if cfg.Temperature != 0 {
		t := cfg.Temperature
		req.Temperature = &t
	}
	if cfg.MaxTokens != 0 {
		m := cfg.MaxTokens
		req.MaxTokens = &m
	}
	if cfg.TopP != 0 {
		p := cfg.TopP
		req.TopP = &p
	}
	return req
}

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
