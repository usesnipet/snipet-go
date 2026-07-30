package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
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
		Generate: func(ctx context.Context, config util.JSONMap, instructions string, messages []llm.Message) (llm.Message, error) {
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
	messages []llm.Message,
	defaultBaseURL string,
) (msg llm.Message, err error) {
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

func buildMessages(instructions string, messages []llm.Message) ([]openai.ChatCompletionMessageParamUnion, error) {
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

func toChatMessage(msg llm.Message) (openai.ChatCompletionMessageParamUnion, bool, error) {
	switch msg.Role {
	case llm.MessageRoleUser:
		return openai.UserMessage(msg.Content), true, nil
	case llm.MessageRoleSystem:
		return openai.SystemMessage(msg.Content), true, nil
	case llm.MessageRoleAssistant:
		content := msg.Content
		return openai.AssistantMessage(content), true, nil
	default:
		return openai.ChatCompletionMessageParamUnion{}, false, nil
	}
}

func toAgentMessage(out Message) llm.Message {
	role := llm.MessageRoleAssistant

	return llm.Message{
		Role:      role,
		Content:   out.Content,
		Timestamp: time.Now(),
	}
}
