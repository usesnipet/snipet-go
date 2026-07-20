package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/runtime/tool"
	"github.com/usesnipet/snipet/internal/runtime/transport"
)

type Engine struct {
	Tools  *driver.Manager[driver.ITool]
	LLMs   *driver.Manager[driver.ILLM]
	logger *logger.Logger
}

func NewEngine(
	tools *driver.Manager[driver.ITool],
	llms *driver.Manager[driver.ILLM],
	logger *logger.Logger,
) *Engine {
	return &Engine{
		Tools:  tools,
		LLMs:   llms,
		logger: logger,
	}
}

func (e *Engine) validateAgent(agent Agent) error {
	llmConfigs := make([]driver.Configuration, 0, len(agent.LLMs))
	for _, cfg := range agent.LLMs {
		llmConfigs = append(llmConfigs, driver.Configuration(cfg))
	}
	err := e.LLMs.ValidateMultipleConfigurationsByKey(llmConfigs...)
	if err != nil {
		return err
	}

	toolConfigs := make([]driver.Configuration, 0, len(agent.Tools))
	for key, config := range agent.Tools {
		toolConfigs = append(toolConfigs, driver.Configuration{Key: key, Config: config})
	}
	return e.Tools.ValidateMultipleConfigurationsByKey(toolConfigs...)
}

type StartOptions struct {
	ExecutionOptions []ExecutionOption
	OnEvent          func(event IEvent) error
	Agent            Agent
}

func (e *Engine) Start(ctx context.Context, options StartOptions) {
	if options.OnEvent == nil {
		options.OnEvent = func(event IEvent) error {
			return nil
		}
	}

	execution, err := NewExecution(options.ExecutionOptions...)
	if err != nil {
		execution.ErrorMessage = err.Error()
		execution.Status = ExecutionStatusFailed
		options.OnEvent(ExecutionErrorEvent{Execution: execution, ErrorMessage: err.Error()})
		return
	}

	if err := e.validateAgent(options.Agent); err != nil {
		execution.ErrorMessage = err.Error()
		execution.Status = ExecutionStatusFailed
		options.OnEvent(ExecutionErrorEvent{Execution: execution, ErrorMessage: err.Error()})
		return
	}

	e.run(ctx, options, execution)
}

func (e *Engine) run(ctx context.Context, options StartOptions, execution Execution) {
	execution.Status = ExecutionStatusRunning
	options.OnEvent(ExecutionUpdatedEvent{Execution: execution})
	options.OnEvent(ExecutionMessageAddedEvent{Execution: execution, Messages: execution.Messages})

	for {
		if execution.Config.MaxTurns > 0 && execution.Turns >= execution.Config.MaxTurns {
			execution.Status = ExecutionStatusMaxTurns
			execution.ErrorMessage = "Max turns reached"
			options.OnEvent(ExecutionUpdatedEvent{Execution: execution})
			return
		}

		message, err := e.runLLM(ctx, options.Agent, execution.Messages)
		if err != nil {
			execution.Status = ExecutionStatusFailed
			execution.ErrorMessage = err.Error()
			options.OnEvent(ExecutionUpdatedEvent{Execution: execution})
			return
		}

		execution.AddMessage(message)
		options.OnEvent(ExecutionMessageAddedEvent{Execution: execution, Messages: []transport.Message{message}})

		if message.Role == transport.MessageRoleFinal || len(message.ToolCalls) == 0 {
			execution.Status = ExecutionStatusCompleted
			options.OnEvent(ExecutionUpdatedEvent{Execution: execution})
			return
		}
		execution.Turns++

		for _, call := range message.ToolCalls {
			result := e.executeTool(ctx, options.Agent, call)
			toolMsg := transport.Message{
				Role:       transport.MessageRoleTool,
				ToolResult: &result,
				Timestamp:  time.Now(),
			}
			if result.Output != nil {
				toolMsg.Content = fmt.Sprint(result.Output)
			} else if result.Error != nil {
				toolMsg.Content = result.Error.Error()
			}
			execution.AddMessage(toolMsg)
			options.OnEvent(ExecutionMessageAddedEvent{Execution: execution, Messages: []transport.Message{toolMsg}})
		}
	}
}

func (e *Engine) runLLM(
	ctx context.Context,
	agent Agent,
	messages []transport.Message,
) (message transport.Message, err error) {
	if len(agent.LLMs) == 0 {
		return message, ErrNoLLMConfigured
	}
	for _, llm := range agent.LLMs {
		llmInstance, err := e.LLMs.GetDriver(llm.Key)
		if err != nil {
			continue
		}
		message, err = llmInstance.Generate(ctx, llm.Config, agent.Instructions, messages)
		if err != nil {
			e.logger.Error(err)
			continue
		}
		return message, nil
	}
	return message, ErrLLMGenerationFailed
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
