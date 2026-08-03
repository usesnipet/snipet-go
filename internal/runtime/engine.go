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

type Engine struct {
	LLMs   *DriverManager[llm.Driver]
	Tools  *ToolManager
	logger *logger.Logger
}

func NewEngine(
	llms *DriverManager[llm.Driver],
	tools *ToolManager,
	logger *logger.Logger,
) *Engine {
	return &Engine{
		LLMs:   llms,
		Tools:  tools,
		logger: logger,
	}
}

func (e *Engine) validateAgent(agent *Agent) error {
	if agent == nil {
		return fmt.Errorf("agent is required")
	}
	if len(agent.LLMs) == 0 {
		return ErrNoLLMConfigured
	}
	llmConfigs := make([]Configuration, 0, len(agent.LLMs))
	for _, cfg := range agent.LLMs {
		llmConfigs = append(llmConfigs, Configuration(cfg))
	}
	return e.LLMs.ValidateMultipleConfigurationsByKey(llmConfigs...)
}

func (e *Engine) Validate(execution *Execution) error {
	if err := e.validateAgent(execution.agent); err != nil {
		return err
	}
	return nil
}

func (e *Engine) Start(ctx context.Context, execution *Execution) error {
	if err := e.Validate(execution); err != nil {
		return err
	}
	return e.loop(ctx, execution)
}

func (e *Engine) loop(ctx context.Context, execution *Execution) error {
	err := execution.SetStatus(ctx, ExecutionStatusRunning)
	if err != nil {
		return err
	}

	for {
		result, err := e.step(ctx, execution)
		if err != nil {
			return execution.SetError(ctx, err.Error())
		}
		switch result {
		case StepContinue:
			if err := execution.CompleteTurn(ctx); err != nil {
				return err
			}
		case StepFinish:
			if err := execution.CompleteTurn(ctx); err != nil {
				return err
			}
			return execution.Finish(ctx)
		case StepCancel:
			return execution.Cancel(ctx)
		case StepMaxTurnsReached:
			return execution.SetMaxTurnsReachedError(ctx)
		}
	}
}

func (e *Engine) step(ctx context.Context, execution *Execution) (StepResult, error) {
	if err := ctx.Err(); err != nil {
		return StepCancel, nil
	}

	if execution.Config.MaxTurns > 0 && execution.Turns >= execution.Config.MaxTurns {
		return StepMaxTurnsReached, nil
	}

	message, err := e.runLLM(ctx, execution)
	if err != nil {
		if ctx.Err() != nil {
			return StepCancel, nil
		}
		return StepContinue, err
	}

	if err = execution.AddMessage(ctx, message); err != nil {
		return StepContinue, err
	}

	if len(message.ToolCalls) > 0 {
		if err := e.runTools(ctx, execution, message.ToolCalls); err != nil {
			return StepContinue, err
		}
		return StepContinue, nil
	}

	if message.IsFinal() {
		return StepFinish, nil
	}

	return StepContinue, nil
}

func (e *Engine) runLLM(ctx context.Context, execution *Execution) (message msg.Message, err error) {
	agent := execution.agent
	if len(agent.LLMs) == 0 {
		return message, ErrNoLLMConfigured
	}

	toolset, err := e.Tools.Toolset()
	if err != nil {
		return message, err
	}

	prompt := llm.NewPrompt(llm.WithSystem(agent.Instructions), llm.WithMessages(execution.Messages))
	generateOptions := llm.GenerateOptions{
		Prompt: prompt,
		Tools:  toolset,
	}

	var lastErr error
	for _, llmConfig := range agent.LLMs {
		llmInstance, getErr := e.LLMs.GetDriver(llmConfig.Key)
		if getErr != nil {
			lastErr = getErr
			continue
		}

		events, streamErr := llmInstance.Stream(ctx, llmConfig.Config, generateOptions)
		if streamErr != nil {
			e.logger.Error(streamErr)
			lastErr = streamErr
			continue
		}

		message, consumeErr := e.consumeStream(ctx, execution, events)
		if consumeErr != nil {
			e.logger.Error(consumeErr)
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
// runtime events, and assembles the resulting assistant message.
func (e *Engine) consumeStream(ctx context.Context, execution *Execution, events <-chan llm.StreamEvent) (msg.Message, error) {
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
				if err := execution.publish(ctx, ExecutionMessageDeltaEvent{MessageID: messageID, Content: ev.Text}); err != nil {
					return msg.Message{}, err
				}

			case llm.ToolCallStartedEvent:
				calls[ev.ID] = &pendingToolCall{id: ev.ID, name: ev.Name}
				order = append(order, ev.ID)
				if err := execution.publish(ctx, ExecutionToolCallStartedEvent{MessageID: messageID, ID: ev.ID, Name: ev.Name}); err != nil {
					return msg.Message{}, err
				}

			case llm.ToolCallArgumentsDeltaEvent:
				if call, ok := calls[ev.ID]; ok {
					call.arguments.WriteString(ev.Delta)
				}
				if err := execution.publish(ctx, ExecutionToolCallDeltaEvent{ID: ev.ID, Delta: ev.Delta}); err != nil {
					return msg.Message{}, err
				}

			case llm.ToolCallFinishedEvent:
				call, ok := calls[ev.ID]
				if !ok {
					continue
				}
				arguments, err := parseArguments(call.arguments.String())
				if err != nil {
					e.logger.Error(fmt.Errorf("parse arguments for tool call %s: %w", call.name, err))
					call.err = err
					continue
				}
				call.parsedArguments = arguments
				if err := execution.publish(ctx, ExecutionToolCallCompletedEvent{ID: call.id, Name: call.name, Arguments: arguments}); err != nil {
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

	options := []msg.MessageOption{msg.WithID(messageID), msg.WithToolCalls(toolCalls)}
	if len(toolCalls) == 0 {
		options = append(options, msg.WithFinal())
	}

	return msg.NewMessage(msg.RoleAssistant, content.String(), options...), nil
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

// runTools executes every tool call requested by the last assistant message
// and appends a RoleTool message with each result, so the next turn can feed
// them back to the LLM. Tool failures don't abort the execution — they're
// surfaced to the LLM as an error result so it can react.
func (e *Engine) runTools(ctx context.Context, execution *Execution, calls []tool.Call) error {
	for _, call := range calls {
		result, callErr := e.Tools.Call(ctx, tool.Call{Tool: call.Tool, Arguments: call.Arguments})

		content := result.Result
		errMessage := ""
		if callErr != nil {
			e.logger.Error(callErr)
			errMessage = callErr.Error()
			content = fmt.Sprintf("error: %s", errMessage)
		}

		if err := execution.publish(ctx, ExecutionToolResultEvent{
			ToolCallID: call.ID,
			Tool:       call.Tool,
			Result:     result.Result,
			Error:      errMessage,
		}); err != nil {
			return err
		}

		toolMessage := msg.NewMessage(msg.RoleTool, content, msg.WithToolCallID(call.ID))
		if err := execution.AddMessage(ctx, toolMessage); err != nil {
			return err
		}
	}
	return nil
}
