package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

func New(opts ...Option) llm.API {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return llm.API{
		TestConnection: func(ctx context.Context, config util.JSONMap) error {
			return testConnection(ctx, config, o.baseURL)
		},
		Generate: func(ctx context.Context, config util.JSONMap, instructions string, messages []message.Message) (message.Message, error) {
			return generate(ctx, config, instructions, messages, o.baseURL)
		},
	}
}

func newClient(cfg Config, defaultBaseURL string) openai.Client {
	opts := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey),
	}
	baseURL := defaultBaseURL
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return openai.NewClient(opts...)
}

func testConnection(ctx context.Context, config util.JSONMap, defaultBaseURL string) error {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	client := newClient(cfg, defaultBaseURL)
	_, err = client.Models.List(ctx)
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
	defaultBaseURL string,
) (msg message.Message, err error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return msg, fmt.Errorf("failed to parse config: %w", err)
	}

	client := newClient(cfg, defaultBaseURL)

	chatMessages, err := buildMessages(instructions, messages)
	if err != nil {
		return msg, fmt.Errorf("openai: build messages: %w", err)
	}

	params := openai.ChatCompletionNewParams{
		Model:    shared.ChatModel(cfg.Model),
		Messages: chatMessages,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "agent_message",
					Schema: messageSchema,
					Strict: openai.Bool(true),
				},
			},
		},
	}
	if cfg.Temperature != 0 {
		params.Temperature = openai.Float(cfg.Temperature)
	}
	if cfg.TopP != 0 {
		params.TopP = openai.Float(cfg.TopP)
	}
	if cfg.MaxTokens != 0 {
		params.MaxTokens = openai.Int(int64(cfg.MaxTokens))
	}

	result, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return msg, fmt.Errorf("openai: generate response: %w", err)
	}
	if len(result.Choices) == 0 {
		return msg, fmt.Errorf("openai: empty response")
	}

	text := result.Choices[0].Message.Content
	if text == "" {
		return msg, fmt.Errorf("openai: empty response content")
	}

	var out Message
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		return msg, fmt.Errorf("openai: parse structured output: %w", err)
	}

	return toAgentMessage(out), nil
}

func buildMessages(instructions string, messages []message.Message) ([]openai.ChatCompletionMessageParamUnion, error) {
	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)
	if instructions != "" {
		out = append(out, openai.SystemMessage(instructions))
	}
	for _, msg := range messages {
		item, ok, err := toChatMessage(msg)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, item)
		}
	}
	return out, nil
}

func toChatMessage(msg message.Message) (openai.ChatCompletionMessageParamUnion, bool, error) {
	switch msg.Role {
	case message.MessageRoleUser:
		return openai.UserMessage(msg.Content), true, nil
	case message.MessageRoleSystem:
		return openai.SystemMessage(msg.Content), true, nil
	case message.MessageRoleAssistant, message.MessageRoleFinal:
		content := msg.Content
		if len(msg.ToolCalls) > 0 {
			raw, err := json.Marshal(Message{
				Content:   msg.Content,
				ToolCalls: util.Map(msg.ToolCalls, fromAgentToolCall),
			})
			if err != nil {
				return openai.ChatCompletionMessageParamUnion{}, false, fmt.Errorf("marshal assistant message: %w", err)
			}
			content = string(raw)
		}
		return openai.AssistantMessage(content), true, nil
	case message.MessageRoleTool:
		content := msg.Content
		if msg.ToolResult != nil {
			content = fmt.Sprintf(
				"Tool %q result: %s",
				msg.ToolResult.Key,
				msg.Content,
			)
		}
		return openai.UserMessage(content), true, nil
	default:
		return openai.ChatCompletionMessageParamUnion{}, false, nil
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
