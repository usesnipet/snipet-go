package runtime

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

type Engine struct {
	LLMs   *DriverManager[llm.Driver]
	Tools  *ToolManager
	logger *logger.Logger

	generator    *Generator
	toolExecutor *ToolExecutor
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

		generator:    NewGenerator(llms, logger),
		toolExecutor: NewToolExecutor(tools, logger),
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

	toolset, err := e.Tools.Toolset()
	if err != nil {
		return StepContinue, err
	}

	message, err := e.generator.Generate(ctx, execution, toolset)
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
		if err := e.toolExecutor.Run(ctx, execution, message.ToolCalls); err != nil {
			return StepContinue, err
		}
		return StepContinue, nil
	}

	if message.IsFinal() {
		return StepFinish, nil
	}

	return StepContinue, nil
}
