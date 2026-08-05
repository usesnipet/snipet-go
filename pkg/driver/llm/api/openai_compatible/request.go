package openaicompatible

import (
	"encoding/json"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// buildChatParams translates a GenerateOptions/Config pair into openai-go
// Chat Completions params. Optional numeric fields left at zero are omitted.
func buildChatParams(cfg Config, options llm.GenerateOptions) openai.ChatCompletionNewParams {
	params := openai.ChatCompletionNewParams{
		Model:    cfg.Model,
		Messages: buildMessages(options.Prompt),
		Tools:    buildTools(options.Tools),
	}

	if cfg.MaxTokens != 0 {
		params.MaxTokens = openai.Int(int64(cfg.MaxTokens))
	}
	if cfg.Temperature != 0 {
		params.Temperature = openai.Float(cfg.Temperature)
	}
	if cfg.TopP != 0 {
		params.TopP = openai.Float(cfg.TopP)
	}

	return params
}

// buildMessages converts a llm.Prompt into openai-go message params,
// prepending a system message when Prompt.System is set and dropping any
// message whose Role has no OpenAI equivalent.
func buildMessages(prompt llm.Prompt) []openai.ChatCompletionMessageParamUnion {
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(prompt.Messages)+1)
	if prompt.System != "" {
		messages = append(messages, openai.SystemMessage(prompt.System))
	}
	for _, m := range prompt.Messages {
		switch m.Role {
		case msg.RoleSystem:
			messages = append(messages, openai.SystemMessage(m.Content))
		case msg.RoleUser:
			messages = append(messages, openai.UserMessage(m.Content))
		case msg.RoleAssistant:
			messages = append(messages, buildAssistantMessage(m))
		case msg.RoleTool:
			messages = append(messages, openai.ToolMessage(m.Content, m.ToolCallID))
		}
	}
	return messages
}

// buildAssistantMessage builds an assistant turn, including tool_calls when
// present. Content is always set (even to "") so Ollama-compatible servers
// that reject missing content on tool-call-only messages keep working.
func buildAssistantMessage(m msg.Message) openai.ChatCompletionMessageParamUnion {
	assistant := openai.ChatCompletionAssistantMessageParam{
		Content: openai.ChatCompletionAssistantMessageParamContentUnion{
			OfString: openai.String(m.Content),
		},
	}
	if len(m.ToolCalls) > 0 {
		assistant.ToolCalls = buildToolCalls(m.ToolCalls)
	}
	return openai.ChatCompletionMessageParamUnion{OfAssistant: &assistant}
}

// buildToolCalls converts assistant tool.Call records into openai-go tool
// call params, JSON-encoding each call's arguments (falling back to "{}" if
// they don't marshal).
func buildToolCalls(calls []tool.Call) []openai.ChatCompletionMessageToolCallUnionParam {
	toolCalls := make([]openai.ChatCompletionMessageToolCallUnionParam, 0, len(calls))
	for _, c := range calls {
		arguments, err := json.Marshal(c.Arguments)
		if err != nil {
			arguments = []byte("{}")
		}
		toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallUnionParam{
			OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
				ID: c.ID,
				Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
					Name:      c.Tool,
					Arguments: string(arguments),
				},
			},
		})
	}
	return toolCalls
}

// buildTools converts a tool.Toolset into openai-go tools, returning nil when
// empty and defaulting each tool's parameters to an empty object schema when
// not specified.
func buildTools(toolset tool.Toolset) []openai.ChatCompletionToolUnionParam {
	if len(toolset.Tools) == 0 {
		return nil
	}
	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(toolset.Tools))
	for _, t := range toolset.Tools {
		params := shared.FunctionParameters(t.Parameters)
		if params == nil {
			params = shared.FunctionParameters{
				"type":       "object",
				"properties": map[string]any{},
			}
		}
		tools = append(tools, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        t.Name,
			Description: openai.String(t.Description),
			Parameters:  params,
		}))
	}
	return tools
}
