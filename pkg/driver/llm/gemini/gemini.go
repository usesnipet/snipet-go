package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"google.golang.org/genai"
)

func New() llm.API {
	return llm.API{
		TestConnection: testConnection,
		Generate:       generate,
	}
}

func testConnection(ctx context.Context, config util.JSONMap) error {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	_, err = client.Models.List(ctx, &genai.ListModelsConfig{
		PageSize: 1,
	})
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	return nil
}

func generate(
	ctx context.Context,
	config util.JSONMap,
	instructions string,
	messages []message.Message,
) (msg message.Message, err error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return msg, fmt.Errorf("failed to parse config: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return msg, fmt.Errorf("failed to create client: %w", err)
	}

	contents, err := buildContents(messages)
	if err != nil {
		return msg, fmt.Errorf("gemini: build contents: %w", err)
	}

	genCfg := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: messageSchema,
	}
	if instructions != "" {
		// Empty text is omitted by json omitempty, which leaves an invalid Part with no data.
		genCfg.SystemInstruction = genai.NewContentFromText(instructions, genai.RoleUser)
	}
	if cfg.Temperature != 0 {
		genCfg.Temperature = genai.Ptr(float32(cfg.Temperature))
	}
	if cfg.TopP != 0 {
		genCfg.TopP = genai.Ptr(float32(cfg.TopP))
	}
	if cfg.MaxTokens != 0 {
		genCfg.MaxOutputTokens = int32(cfg.MaxTokens)
	}
	result, err := client.Models.GenerateContent(ctx, cfg.Model, contents, genCfg)
	if err != nil {
		return msg, fmt.Errorf("gemini: generate response: %w", err)
	}

	text := result.Text()
	if text == "" {
		return msg, fmt.Errorf("gemini: empty response")
	}

	var out Message
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		return msg, fmt.Errorf("gemini: parse structured output: %w", err)
	}

	return toAgentMessage(out), nil
}

func buildContents(messages []message.Message) ([]*genai.Content, error) {
	out := make([]*genai.Content, 0, len(messages))
	for _, msg := range messages {
		item, ok, err := toContent(msg)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, item)
		}
	}
	return out, nil
}

func toContent(msg message.Message) (*genai.Content, bool, error) {
	switch msg.Role {
	case message.MessageRoleUser:
		return genai.NewContentFromText(msg.Content, genai.RoleUser), true, nil
	case message.MessageRoleSystem:
		return genai.NewContentFromText(msg.Content, genai.RoleUser), true, nil
	case message.MessageRoleAssistant, message.MessageRoleFinal:
		content := msg.Content
		if len(msg.ToolCalls) > 0 {
			raw, err := json.Marshal(Message{
				Content:   msg.Content,
				ToolCalls: util.Map(msg.ToolCalls, fromAgentToolCall),
			})
			if err != nil {
				return nil, false, fmt.Errorf("marshal assistant message: %w", err)
			}
			content = string(raw)
		}
		return genai.NewContentFromText(content, genai.RoleModel), true, nil
	case message.MessageRoleTool:
		content := msg.Content
		if msg.ToolResult != nil {
			content = fmt.Sprintf(
				"Tool %q result: %s",
				msg.ToolResult.Key,
				msg.Content,
			)
		}
		return genai.NewContentFromText(content, genai.RoleUser), true, nil
	default:
		return nil, false, nil
	}
}

func toAgentMessage(out Message) message.Message {
	toolCalls := make([]tool.Call, 0, len(out.ToolCalls))
	for _, call := range out.ToolCalls {
		toolCalls = append(toolCalls, toAgentToolCall(call))
	}

	role := message.MessageRoleFinal
	if len(toolCalls) > 0 {
		role = message.MessageRoleAssistant
	}

	return message.Message{
		Role:      role,
		Content:   out.Content,
		ToolCalls: toolCalls,
		Timestamp: time.Now(),
	}
}

func toAgentToolCall(call ToolCall) tool.Call {
	var input any
	if call.Input != "" {
		if err := json.Unmarshal([]byte(call.Input), &input); err != nil {
			input = call.Input
		}
	}

	return tool.Call{
		Key:   call.Key,
		Input: input,
	}
}

func fromAgentToolCall(call tool.Call) ToolCall {
	input := ""
	if call.Input != nil {
		if raw, err := json.Marshal(call.Input); err == nil {
			input = string(raw)
		} else {
			input = fmt.Sprint(call.Input)
		}
	}
	return ToolCall{
		Key:   call.Key,
		Input: input,
	}
}
