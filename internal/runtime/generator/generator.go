package generator

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime/execution"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// Generate produces the next assistant message for the execution, trying
// every LLM configured on its agent in order until one succeeds.
func Generate(
	ctx context.Context,
	exe *execution.Execution,
	llmDriverManager *manager.Driver[llm.Driver],
	toolset tool.Toolset,
	logger *logger.Logger,
) (message msg.Message, err error) {
	agent := exe.Agent
	logger.Debugf("starting generation agent=%q", agent.Name)

	llmInstance, err := llmDriverManager.GetDriver(agent.LLM.Key)
	if err != nil {
		return message, err
	}

	model, err := llmInstance.Model(ctx, agent.LLM.Config)
	if err != nil {
		return message, err
	}

	if !model.Can(llm.ModelCapabilitiesToolCall) {
		logger.Errorf("model %q does not support tool calls", model.Name)
		return message, ErrModelNotSupportToolCall
	}

	return stream(
		ctx,
		exe,
		llmInstance,
		agent.LLM.Config,
		llm.GenerateOptions{
			Prompt: llm.NewPrompt(llm.WithSystem(agent.Instructions), llm.WithMessages(exe.Messages)),
			Tools:  toolset,
		},
		logger,
	)
}

func stream(
	ctx context.Context,
	exe *execution.Execution,
	llmInstance llm.Driver,
	config util.JSONMap,
	options llm.GenerateOptions,
	logger *logger.Logger,
) (msg.Message, error) {
	it, streamErr := llmInstance.Stream(ctx, config, options)
	if streamErr != nil {
		logger.Errorf("stream failed: %v", streamErr)
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
			if err := exe.Publish(ctx, execution.MessageDeltaEvent{MessageID: messageID, Content: ev.Text}); err != nil {
				return msg.Message{}, err
			}

		case llm.ToolCallEvent:
			logger.Debugf(
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
		logger.Warnf("message_id=%s stream error: %s", messageID, err)
		pubErr := exe.Publish(ctx, execution.AttemptFailedEvent{MessageID: messageID, Error: err.Error()})
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
	logger.Debugf(
		"message_id=%s stream consumed content_len=%d tool_calls=%d",
		messageID, len(text), len(toolCalls),
	)

	return msg.NewMessage(msg.RoleAssistant, text, messageOptions...), nil
}
