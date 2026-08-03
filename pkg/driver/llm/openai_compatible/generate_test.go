package openaicompatible

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

func TestGenerate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/chat/completions", r.URL.Path)
		require.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		var body chatRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "gpt-4o-mini", body.Model)
		require.Equal(t, "system", body.Messages[0].Role)
		require.Len(t, body.Tools, 1)
		require.False(t, body.Stream)

		_ = json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{{
				Message: &chatMessage{Role: "assistant", Content: "hello world"},
			}},
		})
	}))
	defer server.Close()

	api := New(server.URL + "/v1")
	message, err := api.Generate(t.Context(), util.JSONMap{
		"api_key": "test-key",
		"model":   "gpt-4o-mini",
	}, llm.GenerateOptions{
		Prompt: llm.NewPrompt(
			llm.WithSystem("Be brief."),
			llm.WithMessages([]msg.Message{msg.NewMessage(msg.RoleUser, "hi")}),
		),
		Tools: tool.NewToolset(tool.NewTool("ping", "Ping", nil)),
	})
	require.NoError(t, err)
	require.Equal(t, msg.RoleAssistant, message.Role)
	require.Equal(t, "hello world", message.Content)
	require.True(t, message.IsFinal())
}

func TestStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "text/event-stream", r.Header.Get("Accept"))

		var body chatRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.True(t, body.Stream)

		flusher, ok := w.(http.Flusher)
		require.True(t, ok)
		w.Header().Set("Content-Type", "text/event-stream")

		chunks := []string{
			`{"choices":[{"delta":{"content":"Hel"}}]}`,
			`{"choices":[{"delta":{"content":"lo"}}]}`,
			`{"choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"ping","arguments":""}}]}}]}`,
			`{"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"x\":1}"}}]},"finish_reason":"tool_calls"}]}`,
		}
		for _, chunk := range chunks {
			_, _ = io.WriteString(w, "data: "+chunk+"\n\n")
			flusher.Flush()
		}
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	api := New(server.URL + "/v1")
	events, err := api.Stream(t.Context(), util.JSONMap{
		"api_key": "test-key",
		"model":   "gpt-4o-mini",
	}, llm.GenerateOptions{
		Prompt: llm.NewPrompt(llm.WithMessages([]msg.Message{
			msg.NewMessage(msg.RoleUser, "hi"),
		})),
	})
	require.NoError(t, err)

	var got []llm.StreamEvent
	for event := range events {
		got = append(got, event)
	}

	require.Len(t, got, 6)
	require.Equal(t, "Hel", got[0].(llm.TextDeltaEvent).Text)
	require.Equal(t, "lo", got[1].(llm.TextDeltaEvent).Text)
	require.Equal(t, "call_1", got[2].(llm.ToolCallStartedEvent).ID)
	require.Equal(t, "ping", got[2].(llm.ToolCallStartedEvent).Name)
	require.Equal(t, "call_1", got[3].(llm.ToolCallArgumentsDeltaEvent).ID)
	require.Equal(t, `{"x":1}`, got[3].(llm.ToolCallArgumentsDeltaEvent).Delta)
	require.Equal(t, "call_1", got[4].(llm.ToolCallFinishedEvent).ID)
	_, ok := got[5].(llm.CompletedEvent)
	require.True(t, ok)
}

func TestStreamCompleted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher := w.(http.Flusher)
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\n")
		flusher.Flush()
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	api := New(server.URL + "/v1")
	events, err := api.Stream(context.Background(), util.JSONMap{
		"api_key": "k",
		"model":   "m",
	}, llm.GenerateOptions{
		Prompt: llm.NewPrompt(llm.WithMessages([]msg.Message{
			msg.NewMessage(msg.RoleUser, "hi"),
		})),
	})
	require.NoError(t, err)

	var texts []string
	var completed bool
	for event := range events {
		switch e := event.(type) {
		case llm.TextDeltaEvent:
			texts = append(texts, e.Text)
		case llm.CompletedEvent:
			completed = true
		}
	}
	require.Equal(t, []string{"ok"}, texts)
	require.True(t, completed)
}

func TestGenerateAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(chatResponse{
			Error: &apiError{Message: "invalid api key", Type: "invalid_request_error"},
		})
	}))
	defer server.Close()

	api := New(server.URL + "/v1")
	_, err := api.Generate(t.Context(), util.JSONMap{
		"api_key": "bad",
		"model":   "m",
	}, llm.GenerateOptions{
		Prompt: llm.NewPrompt(llm.WithMessages([]msg.Message{
			msg.NewMessage(msg.RoleUser, "hi"),
		})),
	})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "invalid api key"))
}

func TestGenerateRequiresModel(t *testing.T) {
	api := New("https://api.openai.com/v1")
	_, err := api.Generate(t.Context(), util.JSONMap{"api_key": "k"}, llm.GenerateOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "model is required")
}
