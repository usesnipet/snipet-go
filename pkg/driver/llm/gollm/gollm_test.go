package gollm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/msg"
)

func TestBuildPrompt(t *testing.T) {
	prompt, err := buildPrompt("You are helpful.", []msg.Message{
		msg.NewMessage(msg.RoleUser, "Hello"),
		msg.NewMessage(msg.RoleAssistant, "Hi there"),
		msg.NewMessage(msg.RoleUser, "How are you?"),
	})
	require.NoError(t, err)
	require.Equal(t, "How are you?", prompt.Input)
	require.Equal(t, "You are helpful.", prompt.SystemPrompt)
	require.Equal(t, outputSpec, prompt.Output)
	require.Len(t, prompt.Messages, 3)
	require.Equal(t, "user", prompt.Messages[0].Role)
	require.Equal(t, "assistant", prompt.Messages[1].Role)
	require.Equal(t, "user", prompt.Messages[2].Role)
}

func TestExtractJSON(t *testing.T) {
	require.Equal(t, `{"content":"ok"}`, extractJSON(`{"content":"ok"}`))
	require.Equal(t, `{"content":"ok"}`, extractJSON("```json\n{\"content\":\"ok\"}\n```"))
	require.Equal(t, `{"content":"ok"}`, extractJSON("```\n{\"content\":\"ok\"}\n```"))
}

func TestNewClientRequiresModel(t *testing.T) {
	api := New("openai")
	_, err := api.Generate(t.Context(), util.JSONMap{
		"api_key": "test-key",
	}, "", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "model is required")
}
