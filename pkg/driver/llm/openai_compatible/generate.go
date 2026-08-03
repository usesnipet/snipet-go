package openaicompatible

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

func generate(ctx context.Context, baseURL string, config util.JSONMap, options llm.GenerateOptions) (msg.Message, error) {
	cfg, err := NewConfig(config)
	if err != nil {
		return msg.Message{}, err
	}

	body := buildChatRequest(cfg, options, false)
	res, err := doChatCompletion(ctx, baseURL, cfg, body)
	if err != nil {
		return msg.Message{}, err
	}
	if len(res.Choices) == 0 || res.Choices[0].Message == nil {
		return msg.Message{}, fmt.Errorf("empty response from provider")
	}

	content := res.Choices[0].Message.Content
	return msg.NewMessage(msg.RoleAssistant, content, msg.WithFinal()), nil
}
