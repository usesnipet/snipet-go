package gollm

import (
	gollmlib "github.com/teilomillet/gollm"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

func buildPrompt(prompt llm.Prompt) *gollmlib.Prompt {
	opts := []gollmlib.PromptOption{}
	if prompt.System != "" {
		opts = append(opts, gollmlib.WithSystemPrompt(prompt.System, gollmlib.CacheTypeEphemeral))
	}

	promptMessages := make([]gollmlib.PromptMessage, 0, len(prompt.Messages))
	for _, m := range prompt.Messages {
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
	for i := len(prompt.Messages) - 1; i >= 0; i-- {
		if prompt.Messages[i].Role == msg.RoleUser && prompt.Messages[i].Content != "" {
			input = prompt.Messages[i].Content
			break
		}
	}

	return gollmlib.NewPrompt(input, opts...)
}

func toGollmRole(role msg.Role) (string, bool) {
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
