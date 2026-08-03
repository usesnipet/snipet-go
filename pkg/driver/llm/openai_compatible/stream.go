package openaicompatible

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

// stream opens a chat completions SSE stream and translates it into
// llm.StreamEvent values on the returned channel. The channel is closed once
// the stream ends, either successfully (llm.CompletedEvent) or with an error
// (llm.ErrorEvent).
func stream(ctx context.Context, baseURL string, config util.JSONMap, options llm.GenerateOptions) (<-chan llm.StreamEvent, error) {
	cfg, err := NewConfig(config)
	if err != nil {
		return nil, err
	}

	body := buildChatRequest(cfg, options, true)
	res, err := openChatCompletionStream(ctx, baseURL, cfg, body)
	if err != nil {
		return nil, err
	}

	out := make(chan llm.StreamEvent)
	go func() {
		defer close(out)
		defer res.Body.Close()
		if err := readSSE(ctx, res.Body, out); err != nil {
			select {
			case <-ctx.Done():
			case out <- llm.ErrorEvent{Error: err.Error()}:
			}
			return
		}
	}()
	return out, nil
}

// readSSE parses the OpenAI-compatible SSE response line by line, emitting
// TextDeltaEvent for content chunks and reassembling tool call deltas
// (indexed by their position in the response, tracked via streamToolCall)
// into ToolCallStartedEvent/ToolCallArgumentsDeltaEvent/ToolCallFinishedEvent.
// It returns nil once a "[DONE]" marker or end of stream is reached after
// emitting a final llm.CompletedEvent.
func readSSE(ctx context.Context, body io.Reader, out chan<- llm.StreamEvent) error {
	scanner := bufio.NewScanner(body)
	// Tool-call argument chunks can be large; raise the default token limit.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	toolCalls := map[int]*streamToolCall{}
	finished := map[string]bool{}

	flushFinished := func() {
		for _, tc := range toolCalls {
			if tc.id == "" || finished[tc.id] {
				continue
			}
			finished[tc.id] = true
			select {
			case <-ctx.Done():
				return
			case out <- llm.ToolCallFinishedEvent{ID: tc.id}:
			}
		}
	}

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			flushFinished()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- llm.CompletedEvent{}:
			}
			return nil
		}

		var chunk chatResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("decode stream chunk: %w", err)
		}
		if chunk.Error != nil {
			return chunk.Error
		}
		if len(chunk.Choices) == 0 || chunk.Choices[0].Delta == nil {
			continue
		}

		delta := chunk.Choices[0].Delta
		if delta.Content != "" {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- llm.TextDeltaEvent{Text: delta.Content}:
			}
		}

		for _, tc := range delta.ToolCalls {
			idx := 0
			if tc.Index != nil {
				idx = *tc.Index
			}
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
			if state.id != "" && state.name != "" && !state.started {
				state.started = true
				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- llm.ToolCallStartedEvent{ID: state.id, Name: state.name}:
				}
			}
			if tc.Function.Arguments != "" && state.started {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- llm.ToolCallArgumentsDeltaEvent{ID: state.id, Delta: tc.Function.Arguments}:
				}
			}
		}

		if chunk.Choices[0].FinishReason != nil && *chunk.Choices[0].FinishReason != "" {
			flushFinished()
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}

	flushFinished()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- llm.CompletedEvent{}:
	}
	return nil
}

// streamToolCall accumulates the id/name for one in-progress tool call
// across SSE delta chunks until enough is known to emit
// ToolCallStartedEvent (tracked via started).
type streamToolCall struct {
	id      string
	name    string
	started bool
}
