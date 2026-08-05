package openaicompatible

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// stream opens a chat completions SSE stream via openai-go and returns an
// llm.StreamIterator that translates chunks into llm.StreamEvent values as
// the caller pulls them via Next.
func stream(ctx context.Context, defaultBaseURL string, config jsonx.JSONMap, options llm.GenerateOptions) (llm.StreamIterator, error) {
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

	return newStreamIterator(sdkStream, sdkStream.Close), nil
}

// sdkChunkStream is the slice of the openai-go SSE stream that
// streamIterator drives. It is narrowed to an interface so tests can fake it.
type sdkChunkStream interface {
	Next() bool
	Current() openai.ChatCompletionChunk
	Err() error
}

// streamIterator adapts an sdkChunkStream, which delivers raw chat
// completion chunks, into an llm.StreamIterator that yields llm.StreamEvent
// values one at a time: text deltas as they arrive, and tool-call deltas
// reassembled (by index) into a single ToolCallEvent once a chunk's
// FinishReason (or end of stream) confirms the call is complete. Malformed
// tool-call arguments are skipped (soft-fail).
type streamIterator struct {
	sdk     sdkChunkStream
	closeFn func() error

	pending []llm.StreamEvent
	event   llm.StreamEvent
	err     error
	done    bool

	toolCalls map[int]*streamToolCall
	flushed   map[string]bool
}

func newStreamIterator(sdk sdkChunkStream, closeFn func() error) *streamIterator {
	return &streamIterator{
		sdk:       sdk,
		closeFn:   closeFn,
		toolCalls: map[int]*streamToolCall{},
		flushed:   map[string]bool{},
	}
}

func (it *streamIterator) Next(ctx context.Context) bool {
	if it.done {
		return false
	}
	for len(it.pending) == 0 {
		if ctx.Err() != nil {
			it.err = ctx.Err()
			it.done = true
			return false
		}
		if !it.sdk.Next() {
			if err := it.sdk.Err(); err != nil {
				it.err = fmt.Errorf("read stream: %w", err)
				it.done = true
				return false
			}
			it.flushToolCalls()
			if len(it.pending) == 0 {
				it.done = true
				return false
			}
			break
		}
		it.consumeChunk(it.sdk.Current())
	}
	it.event, it.pending = it.pending[0], it.pending[1:]
	return true
}

func (it *streamIterator) Event() llm.StreamEvent { return it.event }
func (it *streamIterator) Err() error             { return it.err }

func (it *streamIterator) Close() error {
	if it.closeFn == nil {
		return nil
	}
	return it.closeFn()
}

func (it *streamIterator) consumeChunk(chunk openai.ChatCompletionChunk) {
	if len(chunk.Choices) == 0 {
		return
	}

	choice := chunk.Choices[0]
	delta := choice.Delta

	if delta.Content != "" {
		it.pending = append(it.pending, llm.TextDeltaEvent{Text: delta.Content})
	}

	for _, tc := range delta.ToolCalls {
		idx := int(tc.Index)
		state, ok := it.toolCalls[idx]
		if !ok {
			state = &streamToolCall{}
			it.toolCalls[idx] = state
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
		it.flushToolCalls()
	}
}

func (it *streamIterator) flushToolCalls() {
	for _, tc := range it.toolCalls {
		if tc.id == "" || it.flushed[tc.id] {
			continue
		}
		it.flushed[tc.id] = true
		if event, ok := tc.event(); ok {
			it.pending = append(it.pending, event)
		}
	}
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
