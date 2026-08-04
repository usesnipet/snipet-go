package runtime

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// ToolExecutor executes tool calls requested by an assistant message on
// behalf of an execution.
type ToolExecutor struct {
	tools  *ToolManager
	logger *logger.Logger
}

func NewToolExecutor(tools *ToolManager, logger *logger.Logger) *ToolExecutor {
	return &ToolExecutor{tools: tools, logger: logger}
}

// Run executes every tool call requested by the last assistant message and
// appends a RoleTool message with each result, so the next turn can feed
// them back to the LLM. Tool failures don't abort the execution — they're
// surfaced to the LLM as an error result so it can react.
func (e *ToolExecutor) Run(ctx context.Context, execution *Execution, calls []tool.Call) error {
	e.logger.Debugf("tool_executor: running %d tool call(s)", len(calls))

	for i, call := range calls {
		e.logger.Debugf("tool_executor: [%d/%d] invoking tool=%q id=%s arguments=%v", i+1, len(calls), call.Tool, call.ID, call.Arguments)

		if err := execution.publish(ctx, ExecutionToolCallStartedEvent{
			ToolCallID: call.ID,
			Tool:       call.Tool,
			Arguments:  call.Arguments,
		}); err != nil {
			return err
		}

		result, callErr := e.tools.Call(ctx, tool.Call{Tool: call.Tool, Arguments: call.Arguments})

		content := result.Result
		errMessage := ""
		if callErr != nil {
			e.logger.Errorf("tool_executor: [%d/%d] tool=%q id=%s failed: %v", i+1, len(calls), call.Tool, call.ID, callErr)
			errMessage = callErr.Error()
			content = fmt.Sprintf("error: %s", errMessage)
		} else {
			e.logger.Debugf("tool_executor: [%d/%d] tool=%q id=%s succeeded result_len=%d", i+1, len(calls), call.Tool, call.ID, len(result.Result))
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
