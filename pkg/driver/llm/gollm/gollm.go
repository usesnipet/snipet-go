package gollm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gollmlib "github.com/teilomillet/gollm"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

// New returns an llm.API backed by gollm for the given provider name
// (e.g. "openai", "anthropic", "ollama").
func New(provider string) llm.API {
	return llm.API{
		TestConnection: func(ctx context.Context, config util.JSONMap) error {
			return testConnection(ctx, provider, config)
		},
		Generate: func(
			ctx context.Context,
			config util.JSONMap,
			instructions string,
			messages []msg.Message,
		) (msg.Message, error) {
			return generate(ctx, provider, config, instructions, messages)
		},
	}
}

func newClient(provider string, cfg Config) (gollmlib.LLM, error) {
	opts := []gollmlib.ConfigOption{
		gollmlib.SetProvider(provider),
		gollmlib.SetModel(cfg.Model),
		gollmlib.SetLogLevel(gollmlib.LogLevelOff),
	}

	if cfg.APIKey != "" || provider == "ollama" {
		opts = append(opts, gollmlib.SetAPIKey(cfg.APIKey))
	}
	if cfg.Temperature != 0 {
		opts = append(opts, gollmlib.SetTemperature(cfg.Temperature))
	}
	if cfg.MaxTokens != 0 {
		opts = append(opts, gollmlib.SetMaxTokens(cfg.MaxTokens))
	}
	if cfg.TopP != 0 {
		opts = append(opts, gollmlib.SetTopP(cfg.TopP))
	}
	if cfg.Endpoint != "" {
		switch provider {
		case "ollama":
			opts = append(opts, gollmlib.SetOllamaEndpoint(cfg.Endpoint))
		case "vllm":
			opts = append(opts, gollmlib.SetVLLMEndpoint(cfg.Endpoint))
		}
	}

	return gollmlib.NewLLM(opts...)
}

func testConnection(ctx context.Context, provider string, config util.JSONMap) error {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	if cfg.Model == "" {
		return fmt.Errorf("model is required")
	}

	client, err := newClient(provider, cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	prompt := gollmlib.NewPrompt(
		`Respond with {"content":"ok"}`,
		gollmlib.WithOutput(outputSpec),
	)

	_, err = client.Generate(ctx, prompt, gollmlib.WithJSONSchemaValidation())
	if err != nil {
		return fmt.Errorf("failed to generate test response: %w", err)
	}
	return nil
}

func generate(
	ctx context.Context,
	provider string,
	config util.JSONMap,
	instructions string,
	messages []msg.Message,
) (message msg.Message, err error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return message, fmt.Errorf("failed to parse config: %w", err)
	}
	if cfg.Model == "" {
		return message, fmt.Errorf("model is required")
	}

	client, err := newClient(provider, cfg)
	if err != nil {
		return message, fmt.Errorf("failed to create client: %w", err)
	}

	prompt, err := buildPrompt(instructions, messages)
	if err != nil {
		return message, fmt.Errorf("%s: build prompt: %w", provider, err)
	}

	text, err := client.Generate(ctx, prompt, gollmlib.WithJSONSchemaValidation())
	if err != nil {
		return message, fmt.Errorf("%s: generate response: %w", provider, err)
	}
	text = extractJSON(text)
	if text == "" {
		return message, fmt.Errorf("%s: empty response", provider)
	}

	var out structuredMessage
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		return message, fmt.Errorf("%s: parse structured output: %w", provider, err)
	}

	return msg.NewMessage(msg.RoleAssistant, out.Content, msg.WithTimestamp(time.Now())), nil
}

func buildPrompt(instructions string, messages []msg.Message) (*gollmlib.Prompt, error) {
	opts := []gollmlib.PromptOption{
		gollmlib.WithOutput(outputSpec),
	}
	if instructions != "" {
		opts = append(opts, gollmlib.WithSystemPrompt(instructions, gollmlib.CacheTypeEphemeral))
	}

	promptMessages := make([]gollmlib.PromptMessage, 0, len(messages))
	for _, m := range messages {
		role, ok := toGollmRole(m.Role)
		if !ok {
			continue
		}
		promptMessages = append(promptMessages, gollmlib.PromptMessage{
			Role:    role,
			Content: m.Content,
		})
	}
	if len(promptMessages) > 0 {
		opts = append(opts, gollmlib.WithMessages(promptMessages))
	}

	input := "Respond to the conversation."
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == msg.RoleUser && messages[i].Content != "" {
			input = messages[i].Content
			break
		}
	}

	return gollmlib.NewPrompt(input, opts...), nil
}

func toGollmRole(role msg.MessageRole) (string, bool) {
	switch role {
	case msg.RoleUser:
		return "user", true
	case msg.RoleAssistant:
		return "assistant", true
	case msg.RoleSystem:
		return "system", true
	default:
		return "", false
	}
}

// extractJSON strips optional markdown fences around a JSON payload.
func extractJSON(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "```") {
		return text
	}
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSpace(text)
	if strings.HasPrefix(strings.ToLower(text), "json") {
		text = strings.TrimSpace(text[4:])
	}
	if i := strings.LastIndex(text, "```"); i >= 0 {
		text = text[:i]
	}
	return strings.TrimSpace(text)
}
