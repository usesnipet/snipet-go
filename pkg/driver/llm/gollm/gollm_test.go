package gollm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

func TestBuildPrompt(t *testing.T) {
	prompt := buildPrompt(
		llm.NewPrompt(
			llm.WithSystem("You are helpful."),
			llm.WithMessages([]msg.Message{
				msg.NewMessage(msg.RoleUser, "Hello"),
				msg.NewMessage(msg.RoleAssistant, "Hi there"),
				msg.NewMessage(msg.RoleUser, "How are you?"),
			}),
		),
	)
	require.Equal(t, "How are you?", prompt.Input)
	require.Equal(t, "You are helpful.", prompt.SystemPrompt)
	require.Len(t, prompt.Messages, 3)
	require.Equal(t, "user", prompt.Messages[0].Role)
	require.Equal(t, "assistant", prompt.Messages[1].Role)
	require.Equal(t, "user", prompt.Messages[2].Role)
}

func TestNewClientRequiresModel(t *testing.T) {
	api := New("openai")
	_, err := api.Generate(t.Context(), util.JSONMap{"api_key": "test-key"}, llm.Prompt{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "model is required")
}
