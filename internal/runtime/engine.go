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

func (e *Engine) validateAgent(agent *Agent) error {
	if agent == nil {
		return fmt.Errorf("agent is required")
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
			return execution.SetError(ctx, "Max turns reached")
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

	message, err := e.runLLM(ctx, execution.agent, execution.Messages)
	if err != nil {
		return StepContinue, err
	}

	if err = execution.AddMessage(ctx, message); err != nil {
		return StepContinue, err
	}

	if message.IsFinal() {
		return StepFinish, nil
	}

	return StepContinue, nil
}

func (e *Engine) runLLM(ctx context.Context, agent *Agent, messages []msg.Message) (msg msg.Message, err error) {
	if len(agent.LLMs) == 0 {
		return msg, ErrNoLLMConfigured
	}
	var lastErr error

	prompt := llm.NewPrompt(llm.WithSystem(agent.Instructions), llm.WithMessages(messages))
	for _, llm := range agent.LLMs {
		llmInstance, getErr := e.LLMs.GetDriver(llm.Key)
		if getErr != nil {
			lastErr = getErr
			continue
		}

		msg, genErr := llmInstance.Generate(ctx, llm.Config, prompt)
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
