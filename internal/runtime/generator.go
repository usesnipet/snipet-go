package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// Generator talks to the LLM drivers configured on an agent: it builds the
// prompt, streams a completion (retrying across configured LLMs on
// failure), and assembles the resulting assistant message.
type Generator struct {
	llms   *DriverManager[llm.Driver]
	logger *logger.Logger
}

func NewGenerator(llms *DriverManager[llm.Driver], logger *logger.Logger) *Generator {
	return &Generator{llms: llms, logger: logger}
}

// Generate produces the next assistant message for the execution, trying
// every LLM configured on its agent in order until one succeeds.
func (g *Generator) Generate(ctx context.Context, execution *Execution, toolset tool.Toolset) (message msg.Message, err error) {
	agent := execution.agent
	if len(agent.LLMs) == 0 {
		return message, ErrNoLLMConfigured
	}

	prompt := llm.NewPrompt(llm.WithSystem(agent.Instructions), llm.WithMessages(execution.Messages))

	var lastErr error
	for _, llmConfig := range agent.LLMs {
		llmInstance, getErr := g.llms.GetDriver(llmConfig.Key)
		if getErr != nil {
			lastErr = getErr
			continue
		}

		generateOptions := llm.GenerateOptions{Prompt: prompt, Tools: toolset}
		fallback := false
		if len(toolset.Tools) > 0 {
			caps, capErr := llmInstance.Capabilities(ctx, llmConfig.Config)
			if capErr != nil {
				g.logger.Error(fmt.Errorf("check capabilities for %s: %w", llmConfig.Key, capErr))
			} else if !caps.ToolCall {
				generateOptions = withToolCallFallback(generateOptions)
				fallback = true
			}
		}

		events, streamErr := llmInstance.Stream(ctx, llmConfig.Config, generateOptions)
		if streamErr != nil {
			g.logger.Error(streamErr)
			lastErr = streamErr
			continue
		}

		message, consumeErr := g.consumeStream(ctx, execution, events, fallback)
		if consumeErr != nil {
			g.logger.Error(consumeErr)
			lastErr = consumeErr
			if ctx.Err() != nil {
				break
			}
			continue
		}
		return message, nil
	}
	if lastErr != nil {
		return message, fmt.Errorf("%w: %v", ErrLLMGenerationFailed, lastErr)
	}
	return message, ErrLLMGenerationFailed
}

// pendingToolCall accumulates a tool call's streamed arguments until
// ToolCallFinishedEvent is received.
type pendingToolCall struct {
	id        string
	name      string
	arguments strings.Builder
	err       error

	parsedArguments map[string]any
}

// consumeStream ranges over a driver's StreamEvents, re-emitting them as
// runtime events, and assembles the resulting assistant message. When
// fallback is true, the driver has no native tool call support: text isn't
// re-emitted as it arrives, and once the stream ends the buffered content is
// checked for a fallback JSON tool-call envelope (see parseFallbackToolCall)
// instead of relying on ToolCall*Events, which such a driver never emits.
func (g *Generator) consumeStream(ctx context.Context, execution *Execution, events <-chan llm.StreamEvent, fallback bool) (msg.Message, error) {
	messageID := uuid.NewString()

	var content strings.Builder
	var order []string
	calls := make(map[string]*pendingToolCall)

loop:
	for {
		select {
		case <-ctx.Done():
			return msg.Message{}, ctx.Err()
		case event, ok := <-events:
			if !ok {
				break loop
			}
			switch ev := event.(type) {
			case llm.TextDeltaEvent:
				content.WriteString(ev.Text)
				if fallback {
					continue
				}
				if err := execution.publish(ctx, ExecutionMessageDeltaEvent{MessageID: messageID, Content: ev.Text}); err != nil {
					return msg.Message{}, err
				}

			case llm.ToolCallStartedEvent:
				calls[ev.ID] = &pendingToolCall{id: ev.ID, name: ev.Name}
				order = append(order, ev.ID)

			case llm.ToolCallArgumentsDeltaEvent:
				if call, ok := calls[ev.ID]; ok {
					call.arguments.WriteString(ev.Delta)
				}

			case llm.ToolCallFinishedEvent:
				call, ok := calls[ev.ID]
				if !ok {
					continue
				}
				arguments, err := parseArguments(call.arguments.String())
				if err != nil {
					g.logger.Error(fmt.Errorf("parse arguments for tool call %s: %w", call.name, err))
					call.err = err
					continue
				}
				call.parsedArguments = arguments
				if err := execution.publish(ctx, ExecutionToolCallEvent{
					MessageID: messageID,
					ID:        call.id,
					Name:      call.name,
					Arguments: arguments,
				}); err != nil {
					return msg.Message{}, err
				}

			case llm.ErrorEvent:
				if pubErr := execution.publish(ctx, ExecutionAttemptFailedEvent{MessageID: messageID, Error: ev.Error}); pubErr != nil {
					return msg.Message{}, pubErr
				}
				return msg.Message{}, errors.New(ev.Error)

			case llm.CompletedEvent:
				// no-op: consumption ends when the channel closes.
			}
		}
	}

	toolCalls := make([]tool.Call, 0, len(order))
	for _, id := range order {
		call := calls[id]
		if call.err != nil {
			continue
		}
		toolCalls = append(toolCalls, tool.Call{ID: call.id, Tool: call.name, Arguments: call.parsedArguments})
	}

	text := content.String()
	if fallback && len(toolCalls) == 0 {
		if call := parseFallbackToolCall(text); call != nil {
			id := uuid.NewString()
			toolCalls = append(toolCalls, tool.Call{ID: id, Tool: call.Name, Arguments: call.Arguments})
			if err := execution.publish(ctx, ExecutionToolCallEvent{
				MessageID: messageID,
				ID:        id,
				Name:      call.Name,
				Arguments: call.Arguments,
			}); err != nil {
				return msg.Message{}, err
			}
			text = ""
		} else if err := execution.publish(ctx, ExecutionMessageDeltaEvent{MessageID: messageID, Content: text}); err != nil {
			return msg.Message{}, err
		}
	}

	options := []msg.MessageOption{msg.WithID(messageID), msg.WithToolCalls(toolCalls)}
	if len(toolCalls) == 0 {
		options = append(options, msg.WithFinal())
	}

	return msg.NewMessage(msg.RoleAssistant, text, options...), nil
}

func parseArguments(raw string) (map[string]any, error) {
	if raw == "" {
		return map[string]any{}, nil
	}
	arguments := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &arguments); err != nil {
		return map[string]any{}, err
	}
	return arguments, nil
}
