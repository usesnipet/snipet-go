package openaicompatible

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

// testConnection verifies that config can reach and authenticate against
// baseURL by issuing a minimal real completion request.
func testConnection(ctx context.Context, baseURL string, config util.JSONMap) error {
	cfg, err := NewConfig(config)
	if err != nil {
		return err
	}

	// Keep the probe cheap: one short completion proves auth, endpoint, and model.
	probe := cfg
	if probe.MaxTokens == 0 {
		probe.MaxTokens = 5
	}
	probeConfig, err := util.ToJSONMap(probe)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	_, err = stream(ctx, baseURL, probeConfig, llm.GenerateOptions{
		Prompt: llm.NewPrompt(
			llm.WithMessages([]msg.Message{
				msg.NewMessage(msg.RoleUser, `Respond with "ok"`),
			}),
		),
	})
	if err != nil {
		return fmt.Errorf("failed to generate test response: %w", err)
	}
	return nil
}
