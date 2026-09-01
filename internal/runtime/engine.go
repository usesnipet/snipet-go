package runtime

import (
	"context"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime/execution"
	"github.com/usesnipet/snipet/internal/runtime/generator"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	toolexecutor "github.com/usesnipet/snipet/internal/runtime/tool_executor"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type Engine struct {
	llms   *manager.DriverManager[llm.Driver]
	tools  *manager.Toolbox
	logger *logger.Logger

	generator    *generator.Generator
	toolExecutor *toolexecutor.ToolExecutor
}

func NewEngine(
	llms *manager.DriverManager[llm.Driver],
	tools *manager.Toolbox,
	log *logger.Logger,
) *Engine {
	return &Engine{
		llms:   llms,
		tools:  tools,
		logger: log,

		generator:    generator.NewGenerator(llms, log.Child(logger.WithPrefix("generator: "))),
		toolExecutor: toolexecutor.NewToolExecutor(tools, log.Child(logger.WithPrefix("tool_executor: "))),
	}
}

func (e *Engine) validateAgent(agent *execution.Agent) error {
	if agent == nil {
		return ErrAgentNotConfigured
	}
	if agent.LLM.Key == "" || agent.LLM.Config == nil {
		return ErrNoLLMConfigured
	}
	return e.llms.ValidateConfigurationByKey(agent.LLM.Key, agent.LLM.Config)
}

// Validate reports whether exe is runnable by this engine (currently: its
// agent and the agent's LLM config are well formed).
func (e *Engine) Validate(exe *execution.Execution) error {
	return e.validateAgent(exe.Agent)
}

func (e *Engine) Start(ctx context.Context, execution *execution.Execution) error {
	if err := e.Validate(execution); err != nil {
		name := "unknown"
		if execution.Agent != nil {
			name = execution.Agent.Name
		}
		e.logger.Warnf("engine: validation failed agent=%q: %v", name, err)
		return err
	}
	e.logger.Debugf("engine: starting execution agent=%q max_turns=%d", execution.Agent.Name, execution.Config.MaxTurns)
	if err := execution.Start(ctx); err != nil {
		return err
	}
	return e.loop(ctx, execution)
}

func (e *Engine) loop(ctx context.Context, exe *execution.Execution) error {
	err := exe.SetStatus(ctx, execution.StatusRunning)
	if err != nil {
		return err
	}

	for {
		e.logger.Debugf(
			"starting turn %d/%d -------------------------------------------",
			exe.Turns, exe.Config.MaxTurns,
		)
		result, err := e.step(ctx, exe)
		if err != nil {
			e.logger.Errorf("step failed turn=%d: %v", exe.Turns, err)
			return exe.SetError(ctx, err.Error())
		}
		switch result {
		case StepContinue:
			e.logger.Debugf("turn=%d continuing", exe.Turns)
			if err := exe.CompleteTurn(ctx); err != nil {
				return err
			}
		case StepFinish:
			e.logger.Debugf("turn=%d final message reached, finishing execution", exe.Turns)
			if err := exe.CompleteTurn(ctx); err != nil {
				return err
			}
			return exe.Finish(ctx)
		case StepCancel:
			e.logger.Debugf("turn=%d cancelled via context", exe.Turns)
			return exe.Cancel(ctx)
		case StepMaxTurnsReached:
			e.logger.Warnf("max turns reached (%d/%d)", exe.Turns, exe.Config.MaxTurns)
			return exe.SetMaxTurnsReachedError(ctx)
		default:
			e.logger.Errorf("unexpected step result %q turn=%d", result, exe.Turns)
			return exe.SetError(ctx, "engine: unexpected step result")
		}
	}
}

func (e *Engine) step(ctx context.Context, exe *execution.Execution) (StepResult, error) {
	if err := ctx.Err(); err != nil {
		e.logger.Debugf("step cancelled via context: %v", err)
		return StepCancel, nil
	}

	if exe.Config.MaxTurns > 0 && exe.Turns >= exe.Config.MaxTurns {
		return StepMaxTurnsReached, nil
	}

	if err := exe.StartTurn(ctx); err != nil {
		return StepResultInvalid, err
	}

	toolset, err := e.tools.Toolset()
	if err != nil {
		e.logger.DebugErrorf("failed to resolve toolset, continuing with no tools: %v", err)
		// if the toolset cannot be resolved, we use an empty toolset to continue the execution with no tools
		toolset = tool.NewToolset()
	} else {
		e.logger.Debugf("resolved toolset with %d tool(s)", len(toolset.Tools))
	}

	message, err := e.generator.Generate(ctx, exe, toolset)
	if err != nil {
		if ctx.Err() != nil {
			return StepCancel, nil
		}
		return StepResultInvalid, err
	}

	if err = exe.AddMessage(ctx, message); err != nil {
		return StepResultInvalid, err
	}
	e.logger.Debugf("assistant message added id=%s content_len=%d tool_calls=%d final=%v",
		message.ID, len(message.Content), len(message.ToolCalls), message.IsFinal(),
	)

	if len(message.ToolCalls) > 0 {
		e.logger.Debugf("dispatching %d tool call(s)", len(message.ToolCalls))
		if err := e.toolExecutor.Run(ctx, exe, message.ToolCalls); err != nil {
			return StepResultInvalid, err
		}
		return StepContinue, nil
	}

	if message.IsFinal() {
		return StepFinish, nil
	}

	return StepContinue, nil
}
