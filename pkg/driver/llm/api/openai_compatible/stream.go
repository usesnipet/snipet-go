package openaicompatible

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

// stream opens a chat completions SSE stream via openai-go and translates it
// into llm.StreamEvent values on the returned channel. The channel is closed
// once the stream ends, either successfully (llm.CompletedEvent) or with an
// error (llm.ErrorEvent).
func stream(ctx context.Context, defaultBaseURL string, config util.JSONMap, options llm.GenerateOptions) (<-chan llm.StreamEvent, error) {
	cfg, err := NewConfig(config)
	if err != nil {
		return nil, err
	}

	baseURL, err := resolveBaseURL(defaultBaseURL, cfg)
	if err != nil {
		return nil, err
	}

	client := newClient(baseURL, cfg)
	params := buildChatParams(cfg, options)
	sdkStream := client.Chat.Completions.NewStreaming(ctx, params)

	out := make(chan llm.StreamEvent)
	go func() {
		defer close(out)
		defer sdkStream.Close()
		if err := consumeStream(ctx, sdkStream, out); err != nil {
			select {
			case <-ctx.Done():
			case out <- llm.ErrorEvent{Error: err.Error()}:
			}
		}
	}()
	return out, nil
}

// consumeStream reads openai-go chat completion chunks, emitting text deltas
// and reassembling tool-call deltas (by index) into ToolCallEvents. Malformed
// tool-call arguments are skipped (soft-fail). Ends with CompletedEvent.
func consumeStream(ctx context.Context, sdkStream interface {
	Next() bool
	Current() openai.ChatCompletionChunk
	Err() error
}, out chan<- llm.StreamEvent) error {
	toolCalls := map[int]*streamToolCall{}
	flushed := map[string]bool{}

	flushToolCalls := func() error {
		for _, tc := range toolCalls {
			if tc.id == "" || flushed[tc.id] {
				continue
			}
			flushed[tc.id] = true
			event, ok := tc.event()
			if !ok {
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- event:
			}
		}
		return nil
	}

	for sdkStream.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		chunk := sdkStream.Current()
		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]
		delta := choice.Delta

		if delta.Content != "" {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- llm.TextDeltaEvent{Text: delta.Content}:
			}
		}

		for _, tc := range delta.ToolCalls {
			idx := int(tc.Index)
			state, ok := toolCalls[idx]
			if !ok {
				state = &streamToolCall{}
				toolCalls[idx] = state
			}
			if tc.ID != "" {
				state.id = tc.ID
			}
			if tc.Function.Name != "" {
				state.name = tc.Function.Name
			}
			if tc.Function.Arguments != "" {
				state.arguments.WriteString(tc.Function.Arguments)
			}
		}

		if choice.FinishReason != "" {
			if err := flushToolCalls(); err != nil {
				return err
			}
		}
	}
	if err := sdkStream.Err(); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}

	if err := flushToolCalls(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- llm.CompletedEvent{}:
	}
	return nil
}

// streamToolCall accumulates one in-progress tool call's id/name/arguments
// across stream delta chunks until finish_reason (or end of stream) is seen.
type streamToolCall struct {
	id        string
	name      string
	arguments strings.Builder
}

// event assembles the accumulated state into a ToolCallEvent. ok is false
// when arguments aren't valid JSON (call is skipped).
func (tc *streamToolCall) event() (llm.StreamEvent, bool) {
	arguments, err := parseToolCallArguments(tc.arguments.String())
	if err != nil {
		return nil, false
	}
	return llm.ToolCallEvent{
		ToolCall: tool.Call{
			ID:        tc.id,
			Tool:      tc.name,
			Arguments: arguments,
		},
	}, true
}

// parseToolCallArguments decodes a tool call's accumulated JSON arguments,
// treating an empty payload as an empty object.
func parseToolCallArguments(raw string) (map[string]any, error) {
	if raw == "" {
		return map[string]any{}, nil
	}
	arguments := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &arguments); err != nil {
		return nil, err
	}
	return arguments, nil
}
