package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type Engine struct {
	Tools  *DriverManager[tool.Driver]
	LLMs   *DriverManager[llm.Driver]
	logger *logger.Logger
}

func NewEngine(
	tools *DriverManager[tool.Driver],
	llms *DriverManager[llm.Driver],
	logger *logger.Logger,
) *Engine {
	return &Engine{
		Tools:  tools,
		LLMs:   llms,
		logger: logger,
	}
}

func (e *Engine) validateAgent(agent Agent) error {
	llmConfigs := make([]Configuration, 0, len(agent.LLMs))
	for _, cfg := range agent.LLMs {
		llmConfigs = append(llmConfigs, Configuration(cfg))
	}
	err := e.LLMs.ValidateMultipleConfigurationsByKey(llmConfigs...)
	if err != nil {
		return err
	}

	toolConfigs := make([]Configuration, 0, len(agent.Tools))
	for key, config := range agent.Tools {
		toolConfigs = append(toolConfigs, Configuration{Key: key, Config: config})
	}
	return e.Tools.ValidateMultipleConfigurationsByKey(toolConfigs...)
}

type StartOptions struct {
	ExecutionOptions []ExecutionOption
	OnEvent          EventListener
	Agent            Agent
}

func (e *Engine) emit(onEvent EventListener, event IEvent) error {
	return onEvent(event)
}

func statusChanged(execution Execution) ExecutionStatusChangedEvent {
	return ExecutionStatusChangedEvent{
		Status:       execution.Status,
		ErrorMessage: execution.ErrorMessage,
		Turns:        execution.Turns,
	}
}

func (e *Engine) fail(execution *Execution, onEvent EventListener, err error) error {
	execution.ErrorMessage = err.Error()
	execution.Status = ExecutionStatusFailed
	if emitErr := e.emit(onEvent, statusChanged(*execution)); emitErr != nil {
		return emitErr
	}
	return err
}

func (e *Engine) Start(ctx context.Context, options StartOptions) error {
	if options.OnEvent == nil {
		options.OnEvent = func(event IEvent) error {
			return nil
		}
	}

	execution, err := NewExecution(options.ExecutionOptions...)
	if err != nil {
		return e.fail(&execution, options.OnEvent, err)
	}

	if err := e.validateAgent(options.Agent); err != nil {
		return e.fail(&execution, options.OnEvent, err)
	}

	return e.run(ctx, options, execution)
}

func (e *Engine) run(ctx context.Context, options StartOptions, execution Execution) error {
	execution.Status = ExecutionStatusRunning
	if err := e.emit(options.OnEvent, statusChanged(execution)); err != nil {
		return err
	}
	if err := e.emit(options.OnEvent, ExecutionMessageAddedEvent{Messages: execution.Messages}); err != nil {
		return err
	}

	for {
		if err := ctx.Err(); err != nil {
			execution.Status = ExecutionStatusCancelled
			execution.ErrorMessage = err.Error()
			if emitErr := e.emit(options.OnEvent, statusChanged(execution)); emitErr != nil {
				return emitErr
			}
			return err
		}

		if execution.Config.MaxTurns > 0 && execution.Turns >= execution.Config.MaxTurns {
			execution.Status = ExecutionStatusMaxTurns
			execution.ErrorMessage = "Max turns reached"
			if err := e.emit(options.OnEvent, statusChanged(execution)); err != nil {
				return err
			}
			return nil
		}

		msg, err := e.runLLM(ctx, options.Agent, execution.Messages)
		if err != nil {
			return e.fail(&execution, options.OnEvent, err)
		}

		msg = execution.AddMessage(msg)
		if err := e.emit(options.OnEvent, ExecutionMessageAddedEvent{
			Messages: []message.Message{msg},
		}); err != nil {
			return err
		}

		if msg.Role == message.MessageRoleFinal || len(msg.ToolCalls) == 0 {
			execution.Status = ExecutionStatusCompleted
			if err := e.emit(options.OnEvent, statusChanged(execution)); err != nil {
				return err
			}
			return nil
		}
		execution.Turns++

		for _, call := range msg.ToolCalls {
			result := e.executeTool(ctx, options.Agent, call)
			toolMsg := message.Message{
				Role:       message.MessageRoleTool,
				ToolResult: &result,
				Content:    toolMessageContent(result),
				Timestamp:  time.Now(),
			}
			toolMsg = execution.AddMessage(toolMsg)
			if err := e.emit(options.OnEvent, ExecutionMessageAddedEvent{
				Messages: []message.Message{toolMsg},
			}); err != nil {
				return err
			}
		}
	}
}

func toolMessageContent(result tool.Result) string {
	if result.Output != nil {
		return fmt.Sprint(result.Output)
	}
	if result.Error != nil {
		return result.Error.Error()
	}
	return ""
}

func (e *Engine) runLLM(
	ctx context.Context,
	agent Agent,
	messages []message.Message,
) (msg message.Message, err error) {
	if len(agent.LLMs) == 0 {
		return msg, ErrNoLLMConfigured
	}
	var lastErr error
	for _, llm := range agent.LLMs {
		llmInstance, getErr := e.LLMs.GetDriver(llm.Key)
		if getErr != nil {
			lastErr = getErr
			continue
		}
		msg, genErr := llmInstance.Generate(ctx, llm.Config, agent.Instructions, messages)
		if genErr != nil {
			e.logger.Error(genErr)
			lastErr = genErr
			continue
		}
		return msg, nil
	}
	if lastErr != nil {
		return msg, fmt.Errorf("%w: %v", ErrLLMGenerationFailed, lastErr)
	}
	return msg, ErrLLMGenerationFailed
}

func (e *Engine) executeTool(ctx context.Context, agent Agent, call tool.Call) tool.Result {
	toolInstance, err := e.Tools.GetDriver(call.Key)
	if err != nil {
		return tool.Result{
			Key:     call.Key,
			Success: false,
			Error:   err,
		}
	}

	toolConfig := agent.Tools[call.Key]
	if toolConfig == nil {
		return tool.Result{
			Key:     call.Key,
			Success: false,
			Error:   ErrToolNotFound,
		}
	}

	return toolInstance.Execute(ctx, toolConfig, call)
}
