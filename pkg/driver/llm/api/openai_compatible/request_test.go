package openaicompatible

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

func TestBuildChatRequest(t *testing.T) {
	cfg := Config{
		Model:       "gpt-4o-mini",
		Temperature: 0.5,
		MaxTokens:   100,
		TopP:        0.9,
	}
	options := llm.GenerateOptions{
		Prompt: llm.NewPrompt(
			llm.WithSystem("You are helpful."),
			llm.WithMessages([]msg.Message{
				msg.NewMessage(msg.RoleUser, "Hello"),
				msg.NewMessage(msg.RoleAssistant, "Hi"),
				msg.NewMessage(msg.RoleUser, "How are you?"),
			}),
		),
		Tools: tool.NewToolset(
			tool.NewTool("get_weather", "Get weather", map[string]any{
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{"type": "string"},
				},
			}),
		),
	}

	req := buildChatRequest(cfg, options, false)
	require.Equal(t, "gpt-4o-mini", req.Model)
	require.False(t, req.Stream)
	require.NotNil(t, req.Temperature)
	require.Equal(t, 0.5, *req.Temperature)
	require.NotNil(t, req.MaxTokens)
	require.Equal(t, 100, *req.MaxTokens)
	require.NotNil(t, req.TopP)
	require.Equal(t, 0.9, *req.TopP)

	require.Len(t, req.Messages, 4)
	require.Equal(t, "system", req.Messages[0].Role)
	require.Equal(t, "You are helpful.", req.Messages[0].Content)
	require.Equal(t, "user", req.Messages[1].Role)
	require.Equal(t, "assistant", req.Messages[2].Role)
	require.Equal(t, "user", req.Messages[3].Role)
	require.Equal(t, "How are you?", req.Messages[3].Content)

	require.Len(t, req.Tools, 1)
	require.Equal(t, "function", req.Tools[0].Type)
	require.Equal(t, "get_weather", req.Tools[0].Function.Name)
	require.Equal(t, "Get weather", req.Tools[0].Function.Description)

	payload, err := json.Marshal(req)
	require.NoError(t, err)
	require.Contains(t, string(payload), `"tools"`)
	require.NotContains(t, string(payload), `"stream"`)
}

func TestBuildChatRequestStreamAndEmptyTools(t *testing.T) {
	req := buildChatRequest(Config{Model: "m"}, llm.GenerateOptions{
		Prompt: llm.NewPrompt(llm.WithMessages([]msg.Message{
			msg.NewMessage(msg.RoleUser, "hi"),
		})),
	}, true)
	require.True(t, req.Stream)
	require.Nil(t, req.Tools)
	require.Nil(t, req.Temperature)
	require.Len(t, req.Messages, 1)
}

func TestBuildMessagesRoundTripsToolCalls(t *testing.T) {
	assistant := msg.NewMessage(msg.RoleAssistant, "", msg.WithID("assistant-1"), msg.WithToolCalls([]tool.Call{
		{ID: "call-1", Tool: "get_weather", Arguments: map[string]any{"city": "SP"}},
	}))
	toolResult := msg.NewMessage(msg.RoleTool, "sunny", msg.WithToolCallID("call-1"))

	req := buildChatRequest(Config{Model: "m"}, llm.GenerateOptions{
		Prompt: llm.NewPrompt(llm.WithMessages([]msg.Message{
			msg.NewMessage(msg.RoleUser, "weather?"),
			assistant,
			toolResult,
		})),
	}, false)

	require.Len(t, req.Messages, 3)

	assistantMsg := req.Messages[1]
	require.Equal(t, "assistant", assistantMsg.Role)
	require.Len(t, assistantMsg.ToolCalls, 1)
	require.Equal(t, "call-1", assistantMsg.ToolCalls[0].ID)
	require.Equal(t, "function", assistantMsg.ToolCalls[0].Type)
	require.Equal(t, "get_weather", assistantMsg.ToolCalls[0].Function.Name)
	require.JSONEq(t, `{"city":"SP"}`, assistantMsg.ToolCalls[0].Function.Arguments)

	toolMsg := req.Messages[2]
	require.Equal(t, "tool", toolMsg.Role)
	require.Equal(t, "call-1", toolMsg.ToolCallID)
	require.Equal(t, "sunny", toolMsg.Content)
}

func TestBuildToolsNilParameters(t *testing.T) {
	tools := buildTools(tool.NewToolset(tool.NewTool("noop", "No op", nil)))
	require.Len(t, tools, 1)
	require.Equal(t, "object", tools[0].Function.Parameters["type"])
}

func TestResolveBaseURL(t *testing.T) {
	url, err := resolveBaseURL("https://api.openai.com/v1", Config{})
	require.NoError(t, err)
	require.Equal(t, "https://api.openai.com/v1", url)

	url, err = resolveBaseURL("https://api.openai.com/v1", Config{Endpoint: "https://proxy.example/v1/"})
	require.NoError(t, err)
	require.Equal(t, "https://proxy.example/v1", url)

	url, err = resolveBaseURL("https://api.openai.com/v1/", Config{})
	require.NoError(t, err)
	require.Equal(t, "https://api.openai.com/v1", url)

	_, err = resolveBaseURL("", Config{})
	require.Error(t, err)
}

func TestConfigValidate(t *testing.T) {
	require.Error(t, Config{}.validate())
	require.NoError(t, Config{Model: "gpt"}.validate())
}
