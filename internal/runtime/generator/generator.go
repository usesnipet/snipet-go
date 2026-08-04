package generator

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// Generator talks to the LLM drivers configured on an agent: it builds the
// prompt, streams a completion (retrying across configured LLMs on
// failure), and assembles the resulting assistant message.
type Generator struct {
	llms   *manager.Driver[llm.Driver]
	logger *logger.Logger
}

func NewGenerator(llms *manager.Driver[llm.Driver], logger *logger.Logger) *Generator {
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
	g.logger.Debugf("generator: starting generation agent=%q llms=%d tools=%d messages=%d",
		agent.Name, len(agent.LLMs), len(toolset.Tools), len(execution.Messages))

	var lastErr error
	for i, llmConfig := range agent.LLMs {
		attempt := i + 1
		g.logger.Debugf("generator: attempt %d/%d llm=%q", attempt, len(agent.LLMs), llmConfig.Key)

		llmInstance, getErr := g.llms.GetDriver(llmConfig.Key)
		if getErr != nil {
			g.logger.Warnf("generator: attempt %d/%d llm=%q driver not found: %v", attempt, len(agent.LLMs), llmConfig.Key, getErr)
			lastErr = getErr
			continue
		}

		generateOptions := llm.GenerateOptions{Prompt: prompt, Tools: toolset}
		fallback := false
		if len(toolset.Tools) > 0 {
			model, modelErr := llmInstance.Model(ctx, llmConfig.Config)
			if modelErr != nil {
				g.logger.Error(fmt.Errorf("check model for %s: %w", llmConfig.Key, modelErr))
			} else {
				canToolCall := model.Can(llm.ModelCapabilitiesToolCall)
				g.logger.Debugf("generator: llm=%q model=%q tool_call=%v", llmConfig.Key, model.Name, canToolCall)
				if !canToolCall {
					g.logger.Infof("generator: llm=%q model=%q has no native tool call support, switching to JSON fallback prompt (%d tools)", llmConfig.Key, model.Name, len(toolset.Tools))
					generateOptions = withToolCallFallback(generateOptions)
					fallback = true
				}
			}
		}

		g.logger.Debugf("generator: llm=%q streaming request fallback=%v", llmConfig.Key, fallback)
		it, streamErr := llmInstance.Stream(ctx, llmConfig.Config, generateOptions)
		if streamErr != nil {
			g.logger.Errorf("generator: attempt %d/%d llm=%q stream failed: %v", attempt, len(agent.LLMs), llmConfig.Key, streamErr)
			lastErr = streamErr
			continue
		}

		message, consumeErr := g.consumeStream(ctx, execution, it, fallback)
		if consumeErr != nil {
			g.logger.Errorf("generator: attempt %d/%d llm=%q consume failed: %v", attempt, len(agent.LLMs), llmConfig.Key, consumeErr)
			lastErr = consumeErr
			if ctx.Err() != nil {
				g.logger.Debugf("generator: context cancelled, aborting remaining attempts")
				break
			}
			continue
		}

		g.logger.Debugf("generator: attempt %d/%d llm=%q succeeded content_len=%d tool_calls=%d final=%v",
			attempt, len(agent.LLMs), llmConfig.Key, len(message.Content), len(message.ToolCalls), message.IsFinal())
		return message, nil
	}
	if lastErr != nil {
		g.logger.Errorf("generator: all %d attempt(s) failed, last error: %v", len(agent.LLMs), lastErr)
		return message, fmt.Errorf("%w: %v", ErrLLMGenerationFailed, lastErr)
	}
	g.logger.Errorf("generator: all %d attempt(s) failed with no error recorded", len(agent.LLMs))
	return message, ErrLLMGenerationFailed
}

// consumeStream walks a driver's StreamIterator, re-emitting events as
// runtime events, and assembles the resulting assistant message. When
// fallback is true, the driver has no native tool call support: text isn't
// re-emitted as it arrives, and once the stream ends the buffered content is
// checked for a fallback JSON tool-call envelope (see parseFallbackToolCall)
// instead of relying on llm.ToolCallEvent, which such a driver never emits.
func (g *Generator) consumeStream(ctx context.Context, execution *Execution, it llm.StreamIterator, fallback bool) (msg.Message, error) {
	defer it.Close()

	messageID := uuid.NewString()
	g.logger.Debugf("generator: consuming stream message_id=%s fallback=%v", messageID, fallback)

	var content strings.Builder
	var toolCalls []tool.Call

	for it.Next(ctx) {
		switch ev := it.Event().(type) {
		case llm.TextDeltaEvent:
			content.WriteString(ev.Text)
			if fallback {
				continue
			}
			if err := execution.publish(ctx, ExecutionMessageDeltaEvent{MessageID: messageID, Content: ev.Text}); err != nil {
				return msg.Message{}, err
			}

		case llm.ToolCallEvent:
			g.logger.Debugf("generator: message_id=%s tool call requested id=%s name=%s arguments=%v", messageID, ev.ToolCall.ID, ev.ToolCall.Tool, ev.ToolCall.Arguments)
			toolCalls = append(toolCalls, ev.ToolCall)
		}
	}
	if err := it.Err(); err != nil {
		if ctx.Err() != nil {
			return msg.Message{}, err
		}
		g.logger.Warnf("generator: message_id=%s stream error: %s", messageID, err)
		if pubErr := execution.publish(ctx, ExecutionAttemptFailedEvent{MessageID: messageID, Error: err.Error()}); pubErr != nil {
			return msg.Message{}, pubErr
		}
		return msg.Message{}, err
	}

	text := content.String()
	if fallback && len(toolCalls) == 0 {
		if call := parseFallbackToolCall(text); call != nil {
			g.logger.Debugf("generator: message_id=%s fallback parsed tool call name=%s arguments=%v", messageID, call.Name, call.Arguments)
			toolCalls = append(toolCalls, tool.Call{ID: uuid.NewString(), Tool: call.Name, Arguments: call.Arguments})
			text = ""
		} else {
			g.logger.Debugf("generator: message_id=%s fallback produced plain text response len=%d", messageID, len(text))
			if err := execution.publish(ctx, ExecutionMessageDeltaEvent{MessageID: messageID, Content: text}); err != nil {
				return msg.Message{}, err
			}
		}
	}

	options := []msg.MessageOption{msg.WithID(messageID), msg.WithToolCalls(toolCalls)}
	if len(toolCalls) == 0 {
		options = append(options, msg.WithFinal())
	}

	g.logger.Debugf("generator: message_id=%s stream consumed content_len=%d tool_calls=%d", messageID, len(text), len(toolCalls))

	return msg.NewMessage(msg.RoleAssistant, text, options...), nil
}
