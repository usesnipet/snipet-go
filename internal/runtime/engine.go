package runtime

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/msg"
)

type Engine struct {
	LLMs   *DriverManager[llm.Driver]
	logger *logger.Logger
}

func NewEngine(
	llms *DriverManager[llm.Driver],
	logger *logger.Logger,
) *Engine {
	return &Engine{
		LLMs:   llms,
		logger: logger,
	}
}

func (e *Engine) validateAgent(agent Agent) error {
	llmConfigs := make([]Configuration, 0, len(agent.LLMs))
	for _, cfg := range agent.LLMs {
		llmConfigs = append(llmConfigs, Configuration(cfg))
	}
	return e.LLMs.ValidateMultipleConfigurationsByKey(llmConfigs...)
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

	message, err := e.runLLM(ctx, options.Agent, execution.Messages)
	if err != nil {
		return e.fail(&execution, options.OnEvent, err)
	}

	message = execution.AddMessage(message)
	if err := e.emit(options.OnEvent, ExecutionMessageAddedEvent{
		Messages: []msg.Message{message},
	}); err != nil {
		return err
	}

	execution.Status = ExecutionStatusCompleted
	if err := e.emit(options.OnEvent, statusChanged(execution)); err != nil {
		return err
	}
	return nil
}

func (e *Engine) runLLM(
	ctx context.Context,
	agent Agent,
	messages []msg.Message,
) (msg msg.Message, err error) {
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
