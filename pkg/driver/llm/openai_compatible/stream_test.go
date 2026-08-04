package openaicompatible

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

func sseLine(t *testing.T, chunk chatResponse) string {
	t.Helper()
	data, err := json.Marshal(chunk)
	require.NoError(t, err)
	return "data: " + string(data) + "\n\n"
}

func collectSSE(t *testing.T, sse string) []llm.StreamEvent {
	t.Helper()
	out := make(chan llm.StreamEvent, 16)
	err := readSSE(context.Background(), strings.NewReader(sse), out)
	require.NoError(t, err)
	close(out)

	events := make([]llm.StreamEvent, 0, len(out))
	for event := range out {
		events = append(events, event)
	}
	return events
}

func TestReadSSEAssemblesToolCallAcrossDeltas(t *testing.T) {
	var sse strings.Builder
	sse.WriteString(sseLine(t, chatResponse{Choices: []chatChoice{{
		Delta: &chatMessage{ToolCalls: []chatToolCall{{
			Index:    new(0),
			ID:       "call-1",
			Type:     "function",
			Function: chatToolCallFunction{Name: "echo_tool"},
		}}},
	}}}))
	sse.WriteString(sseLine(t, chatResponse{Choices: []chatChoice{{
		Delta: &chatMessage{ToolCalls: []chatToolCall{{
			Index:    new(0),
			Function: chatToolCallFunction{Arguments: `{"msg":`},
		}}},
		FinishReason: nil,
	}}}))
	sse.WriteString(sseLine(t, chatResponse{Choices: []chatChoice{{
		Delta: &chatMessage{ToolCalls: []chatToolCall{{
			Index:    new(0),
			Function: chatToolCallFunction{Arguments: `"hi"}`},
		}}},
		FinishReason: new("tool_calls"),
	}}}))
	sse.WriteString("data: [DONE]\n\n")

	events := collectSSE(t, sse.String())

	require.Len(t, events, 2)
	toolCall, ok := events[0].(llm.ToolCallEvent)
	require.True(t, ok, "expected llm.ToolCallEvent, got %T", events[0])
	require.Equal(t, "call-1", toolCall.ID)
	require.Equal(t, "echo_tool", toolCall.Name)
	require.Equal(t, map[string]any{"msg": "hi"}, toolCall.Arguments)
	require.IsType(t, llm.CompletedEvent{}, events[1])
}

func TestReadSSEEmitsToolCallErrorOnMalformedArguments(t *testing.T) {
	var sse strings.Builder
	sse.WriteString(sseLine(t, chatResponse{Choices: []chatChoice{{
		Delta: &chatMessage{ToolCalls: []chatToolCall{{
			Index:    new(0),
			ID:       "call-1",
			Type:     "function",
			Function: chatToolCallFunction{Name: "echo_tool", Arguments: `{not json`},
		}}},
		FinishReason: new("tool_calls"),
	}}}))
	sse.WriteString("data: [DONE]\n\n")

	events := collectSSE(t, sse.String())

	require.Len(t, events, 2)
	toolCallErr, ok := events[0].(llm.ToolCallErrorEvent)
	require.True(t, ok, "expected llm.ToolCallErrorEvent, got %T", events[0])
	require.Equal(t, "call-1", toolCallErr.ID)
	require.NotEmpty(t, toolCallErr.Error)
	require.IsType(t, llm.CompletedEvent{}, events[1])
}
