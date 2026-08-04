package openaicompatible

import (
	"context"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

type fakeChunkStream struct {
	chunks []openai.ChatCompletionChunk
	idx    int
	err    error
}

func (s *fakeChunkStream) Next() bool {
	if s.idx >= len(s.chunks) {
		return false
	}
	s.idx++
	return true
}

func (s *fakeChunkStream) Current() openai.ChatCompletionChunk {
	return s.chunks[s.idx-1]
}

func (s *fakeChunkStream) Err() error {
	return s.err
}

func collectChunks(t *testing.T, chunks ...openai.ChatCompletionChunk) []llm.StreamEvent {
	t.Helper()
	it := newStreamIterator(&fakeChunkStream{chunks: chunks}, nil)

	events := make([]llm.StreamEvent, 0, len(chunks))
	for it.Next(context.Background()) {
		events = append(events, it.Event())
	}
	require.NoError(t, it.Err())
	return events
}

func TestConsumeStreamAssemblesToolCallAcrossDeltas(t *testing.T) {
	events := collectChunks(t,
		openai.ChatCompletionChunk{Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
					Index:    0,
					ID:       "call-1",
					Type:     "function",
					Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{Name: "echo_tool"},
				}},
			},
		}}},
		openai.ChatCompletionChunk{Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
					Index:    0,
					Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: `{"msg":`},
				}},
			},
		}}},
		openai.ChatCompletionChunk{Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
					Index:    0,
					Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: `"hi"}`},
				}},
			},
			FinishReason: "tool_calls",
		}}},
	)

	require.Len(t, events, 1)
	toolCall, ok := events[0].(llm.ToolCallEvent)
	require.True(t, ok, "expected llm.ToolCallEvent, got %T", events[0])
	require.Equal(t, "call-1", toolCall.ToolCall.ID)
	require.Equal(t, "echo_tool", toolCall.ToolCall.Tool)
	require.Equal(t, map[string]any{"msg": "hi"}, toolCall.ToolCall.Arguments)
}

func TestConsumeStreamSkipsMalformedToolCallArguments(t *testing.T) {
	events := collectChunks(t,
		openai.ChatCompletionChunk{Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
					Index:    0,
					ID:       "call-1",
					Type:     "function",
					Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{Name: "echo_tool", Arguments: `{not json`},
				}},
			},
			FinishReason: "tool_calls",
		}}},
	)

	require.Len(t, events, 0)
}

func TestConsumeStreamEmitsTextDeltas(t *testing.T) {
	events := collectChunks(t,
		openai.ChatCompletionChunk{Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{Content: "Hello"},
		}}},
		openai.ChatCompletionChunk{Choices: []openai.ChatCompletionChunkChoice{{
			Delta:        openai.ChatCompletionChunkChoiceDelta{Content: " world"},
			FinishReason: "stop",
		}}},
	)

	require.Len(t, events, 2)
	require.Equal(t, llm.TextDeltaEvent{Text: "Hello"}, events[0])
	require.Equal(t, llm.TextDeltaEvent{Text: " world"}, events[1])
}

func TestParseToolCallArguments(t *testing.T) {
	args, err := parseToolCallArguments("")
	require.NoError(t, err)
	require.Equal(t, map[string]any{}, args)

	args, err = parseToolCallArguments(`{"a":1}`)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"a": float64(1)}, args)

	_, err = parseToolCallArguments(`{bad`)
	require.Error(t, err)
}
