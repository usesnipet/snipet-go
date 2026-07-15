package agentloop

import (
	"context"
	"fmt"
	"time"
)

type AgentLoop struct {
	Agent        *Agent
	ToolProvider ToolProvider
	LLMProvider  LLMProvider
	Persister    Persister
}

func NewAgentLoop(agent *Agent, llm LLMProvider, tools ToolProvider, persister Persister) *AgentLoop {
	if persister == nil {
		persister = NewNoopPersister()
	}
	return &AgentLoop{
		Agent:        agent,
		LLMProvider:  llm,
		ToolProvider: tools,
		Persister:    persister,
	}
}

func (a *AgentLoop) Run(ctx context.Context, execution *Execution) error {
	if execution.Status == "" {
		execution.Status = ExecutionStatusRunning
	}
	if err := a.Persister.SaveExecution(ctx, execution); err != nil {
		return fmt.Errorf("save execution: %w", err)
	}

	for {
		if execution.MaxTurns > 0 && execution.Turns >= execution.MaxTurns {
			execution.Status = ExecutionStatusMaxTurns
			return a.Persister.SaveExecution(ctx, execution)
		}

		message, err := a.LLMProvider.GenerateResponse(ctx, a.Agent, execution)
		if err != nil {
			execution.Status = ExecutionStatusFailed
			_ = a.Persister.SaveExecution(ctx, execution)
			return fmt.Errorf("generate response: %w", err)
		}

		if err := a.appendMessage(ctx, execution, message); err != nil {
			return err
		}

		if message.Role == MessageRoleFinal {
			execution.Status = ExecutionStatusCompleted
			return a.Persister.SaveExecution(ctx, execution)
		}

		if len(message.ToolCalls) == 0 {
			// Assistant with no tool calls is treated as a completed turn.
			execution.Status = ExecutionStatusCompleted
			return a.Persister.SaveExecution(ctx, execution)
		}

		execution.Turns++
		if err := a.Persister.SaveExecution(ctx, execution); err != nil {
			return fmt.Errorf("save execution turns: %w", err)
		}

		for _, call := range message.ToolCalls {
			result, err := a.ToolProvider.ExecuteTool(ctx, call)
			if err != nil {
				execution.Status = ExecutionStatusFailed
				_ = a.Persister.SaveExecution(ctx, execution)
				return fmt.Errorf("execute tool %q: %w", call.Key, err)
			}

			toolMsg := &Message{
				Role:       MessageRoleTool,
				ToolResult: result,
				Timestamp:  time.Now(),
			}
			if result != nil && result.Output != nil {
				toolMsg.Content = fmt.Sprint(result.Output)
			} else if result != nil && result.Error != nil {
				toolMsg.Content = result.Error.Error()
			}

			if err := a.appendMessage(ctx, execution, toolMsg); err != nil {
				return err
			}
		}
	}
}

func (a *AgentLoop) appendMessage(ctx context.Context, execution *Execution, message *Message) error {
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}
	execution.AddMessage(message)
	if err := a.Persister.AppendMessage(ctx, execution.ID, message); err != nil {
		return fmt.Errorf("append message: %w", err)
	}
	return nil
}
