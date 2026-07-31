package gollm

import (
	"context"
	"fmt"

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
		Generate: func(ctx context.Context, config util.JSONMap, prompt llm.Prompt) (msg.Message, error) {
			return generate(ctx, provider, config, prompt)
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

	prompt := gollmlib.NewPrompt(`Respond with "ok"`)

	_, err = client.Generate(ctx, prompt, gollmlib.WithJSONSchemaValidation())
	if err != nil {
		return fmt.Errorf("failed to generate test response: %w", err)
	}
	return nil
}
