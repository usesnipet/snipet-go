package generator

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime/execution"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
	"github.com/usesnipet/snipet/pkg/msg"
)

// Generator turns an execution's conversation state into the next assistant
// message by streaming from the LLM configured on its agent. It holds the
// engine-lifetime dependencies (the LLM driver manager and a logger); the
// per-turn inputs are passed to Generate.
type Generator struct {
	llms   *manager.DriverManager[llm.Driver]
	logger *logger.Logger
}

func NewGenerator(llms *manager.DriverManager[llm.Driver], logger *logger.Logger) *Generator {
	return &Generator{llms: llms, logger: logger}
}

// Generate produces the next assistant message for the execution using the
// single LLM configured on its agent.
func (g *Generator) Generate(
	ctx context.Context,
	exe *execution.Execution,
	toolset tool.Toolset,
) (msg.Message, error) {
	agent := exe.Agent
	g.logger.Debugf("starting generation agent=%q", agent.Name)

	llmInstance, err := g.llms.GetDriver(agent.LLM.Key)
	if err != nil {
		return msg.Message{}, err
	}

	model, err := llmInstance.Model(ctx, agent.LLM.Config)
	if err != nil {
		return msg.Message{}, err
	}

	if !model.Can(llm.ModelCapabilitiesToolCall) {
		g.logger.Errorf("model %q does not support tool calls", model.Name)
		return msg.Message{}, ErrModelNotSupportToolCall
	}

	return g.stream(
		ctx,
		exe,
		llmInstance,
		agent.LLM.Config,
		llm.GenerateOptions{
			Prompt: llm.NewPrompt(llm.WithSystem(agent.Instructions), llm.WithMessages(exe.Messages)),
			Tools:  toolset,
		},
	)
}

func (g *Generator) stream(
	ctx context.Context,
	exe *execution.Execution,
	llmInstance llm.Driver,
	config jsonx.JSONMap,
	options llm.GenerateOptions,
) (msg.Message, error) {
	it, streamErr := llmInstance.Stream(ctx, config, options)
	if streamErr != nil {
		g.logger.Errorf("stream failed: %v", streamErr)
		return msg.Message{}, streamErr
	}

	defer it.Close()

	messageID := uuid.NewString()

	var content strings.Builder
	var toolCalls []tool.Call

	for it.Next(ctx) {
		switch ev := it.Event().(type) {
		case llm.TextDeltaEvent:
			content.WriteString(ev.Text)
			if exe.StreamMessages {
				err := exe.Publish(ctx, execution.MessageDeltaEvent{MessageID: messageID, Content: ev.Text})
				if err != nil {
					return msg.Message{}, err
				}
			}

		case llm.ToolCallEvent:
			g.logger.Debugf(
				"message_id=%s tool call requested id=%s name=%s arguments=%v",
				messageID, ev.ToolCall.ID, ev.ToolCall.Tool, ev.ToolCall.Arguments,
			)
			toolCalls = append(toolCalls, ev.ToolCall)
		}
	}

	if err := it.Err(); err != nil {
		if ctx.Err() != nil {
			return msg.Message{}, err
		}
		g.logger.Warnf("message_id=%s stream error: %s", messageID, err)
		pubErr := exe.Publish(ctx, execution.MessageAttemptFailedEvent{MessageID: messageID, Error: err.Error()})
		if pubErr != nil {
			return msg.Message{}, pubErr
		}
		return msg.Message{}, err
	}

	text := content.String()

	messageOptions := []msg.MessageOption{msg.WithID(messageID), msg.WithToolCalls(toolCalls)}
	if len(toolCalls) == 0 {
		messageOptions = append(messageOptions, msg.WithFinal())
	}
	g.logger.Debugf(
		"message_id=%s stream consumed content_len=%d tool_calls=%d",
		messageID, len(text), len(toolCalls),
	)

	return msg.NewMessage(msg.RoleAssistant, text, messageOptions...), nil
}
