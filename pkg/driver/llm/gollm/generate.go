package gollm

import (
	"context"
	"fmt"

	gollmlib "github.com/teilomillet/gollm"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

func generate(ctx context.Context, provider string, config util.JSONMap, prompt llm.Prompt) (message msg.Message, err error) {
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

	text, err := client.Generate(
		ctx,
		buildPrompt(prompt),
		gollmlib.WithJSONSchemaValidation(),
	)
	if err != nil {
		return message, fmt.Errorf("%s: generate response: %w", provider, err)
	}
	return msg.NewMessage(msg.RoleAssistant, text, msg.WithFinal()), nil
}
