package runtime_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	llmmocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/msg"
)

// fakeStreamIterator is a fixed-slice llm.StreamIterator for MockDriver.Stream to return.
type fakeStreamIterator struct {
	events []llm.StreamEvent
	idx    int
}

func streamOf(events ...llm.StreamEvent) llm.StreamIterator {
	return &fakeStreamIterator{events: events}
}

func (it *fakeStreamIterator) Next(_ context.Context) bool {
	if it.idx >= len(it.events) {
		return false
	}
	it.idx++
	return true
}

func (it *fakeStreamIterator) Event() llm.StreamEvent { return it.events[it.idx-1] }
func (it *fakeStreamIterator) Err() error             { return nil }
func (it *fakeStreamIterator) Close() error           { return nil }

type recordingSubscriber struct {
	events []runtime.IEvent
}

func (r *recordingSubscriber) Handle(_ context.Context, event runtime.IEvent) error {
	r.events = append(r.events, event)
	return nil
}

func (r *recordingSubscriber) has(want runtime.IEvent) bool {
	for _, event := range r.events {
		if reflect.DeepEqual(event, want) {
			return true
		}
	}
	return false
}

func TestEngineExecutesToolCallsBeforeFinishing(t *testing.T) {
	t.Parallel()

	llmDriver := llmmocks.NewMockDriver(t)
	llmDriver.EXPECT().Info().Return(driver.Info{
		Name:                "primary",
		Description:         "primary",
		ConfigurationSchema: util.JSONMap{"type": "object"},
	})
	llmDriver.EXPECT().
		Model(mock.Anything, mock.Anything).
		Return(llm.NewModel("test-model", "test", []llm.ModelCapabilities{llm.ModelCapabilitiesToolCall}), nil)
	llmDriver.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(
			llm.ToolCallEvent{ToolCall: tool.Call{ID: "call-1", Tool: "echo__echo_tool", Arguments: map[string]any{"msg": "hi"}}},
		), nil).
		Once()
	llmDriver.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(llm.TextDeltaEvent{Text: "done"}), nil).
		Once()

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", llmDriver)

	var gotArguments map[string]any
	toolReg := registry.New[tool.Driver]()
	toolReg.MustRegister("echo", tool.CreateDriver(
		tool.WithName("Echo"),
		tool.WithDescription("Echo"),
		tool.WithToolSet(tool.NewToolset(tool.NewTool("echo_tool", "Echoes input", nil))),
		tool.WithAPI(tool.API{
			Call: func(_ context.Context, call tool.Call) (tool.Result, error) {
				gotArguments = call.Arguments
				return tool.Result{Tool: call.Tool, Arguments: call.Arguments, Result: "echo:hi"}, nil
			},
		}),
	))

	engine := runtime.NewEngine(
		runtime.NewDriverManager(llmReg),
		runtime.NewToolManager(runtime.NewDriverManager(toolReg)),
		logger.NewLogger(logger.LevelError),
	)

	agent := runtime.NewAgent("agent", "agent", "be helpful", []runtime.LLMConfig{
		{Key: "primary", Config: util.JSONMap{"model": "test-model"}},
	})

	execution, err := runtime.NewExecution(
		runtime.WithAgent(agent),
		runtime.WithMessageFromUser("hello"),
	)
	require.NoError(t, err)

	subscriber := &recordingSubscriber{}
	execution.Subscribe(subscriber)

	err = engine.Start(context.Background(), execution)
	require.NoError(t, err)

	require.Equal(t, runtime.ExecutionStatusCompleted, execution.Status)
	require.Equal(t, map[string]any{"msg": "hi"}, gotArguments)

	require.Len(t, execution.Messages, 4)
	assistantWithToolCall := execution.Messages[1]
	require.Equal(t, msg.RoleAssistant, assistantWithToolCall.Role)
	require.Len(t, assistantWithToolCall.ToolCalls, 1)
	require.Equal(t, "echo__echo_tool", assistantWithToolCall.ToolCalls[0].Tool)

	toolResult := execution.Messages[2]
	require.Equal(t, msg.RoleTool, toolResult.Role)
	require.Equal(t, "echo:hi", toolResult.Content)
	require.Equal(t, assistantWithToolCall.ToolCalls[0].ID, toolResult.ToolCallID)

	final := execution.Messages[3]
	require.Equal(t, msg.RoleAssistant, final.Role)
	require.Equal(t, "done", final.Content)
	require.True(t, final.IsFinal())

	require.True(t, subscriber.has(runtime.ExecutionToolCallStartedEvent{
		ToolCallID: "call-1",
		Tool:       "echo__echo_tool",
		Arguments:  map[string]any{"msg": "hi"},
	}))
	require.True(t, subscriber.has(runtime.ExecutionToolResultEvent{
		ToolCallID: "call-1",
		Tool:       "echo__echo_tool",
		Result:     "echo:hi",
	}))
	require.True(t, subscriber.has(runtime.ExecutionMessageDeltaEvent{
		MessageID: final.ID,
		Content:   "done",
	}))
}

func TestEngineSurfacesToolCallErrorAsToolMessage(t *testing.T) {
	t.Parallel()

	llmDriver := llmmocks.NewMockDriver(t)
	llmDriver.EXPECT().Info().Return(driver.Info{
		Name:                "primary",
		Description:         "primary",
		ConfigurationSchema: util.JSONMap{"type": "object"},
	})
	llmDriver.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(
			llm.ToolCallEvent{ToolCall: tool.Call{ID: "call-1", Tool: "unregistered__tool", Arguments: map[string]any{}}},
		), nil).
		Once()
	llmDriver.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(llm.TextDeltaEvent{Text: "done"}), nil).
		Once()

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", llmDriver)

	engine := runtime.NewEngine(
		runtime.NewDriverManager(llmReg),
		runtime.NewToolManager(runtime.NewDriverManager(registry.New[tool.Driver]())),
		logger.NewLogger(logger.LevelError),
	)

	agent := runtime.NewAgent("agent", "agent", "be helpful", []runtime.LLMConfig{
		{Key: "primary", Config: util.JSONMap{}},
	})
	execution, err := runtime.NewExecution(
		runtime.WithAgent(agent),
		runtime.WithMessageFromUser("hello"),
	)
	require.NoError(t, err)

	err = engine.Start(context.Background(), execution)
	require.NoError(t, err)

	require.Equal(t, runtime.ExecutionStatusCompleted, execution.Status)
	toolResult := execution.Messages[2]
	require.Equal(t, msg.RoleTool, toolResult.Role)
	require.Contains(t, toolResult.Content, "error:")
}

func TestEngineFallsBackToJSONToolCallWhenUnsupported(t *testing.T) {
	t.Parallel()

	llmDriver := llmmocks.NewMockDriver(t)
	llmDriver.EXPECT().Info().Return(driver.Info{
		Name:                "primary",
		Description:         "primary",
		ConfigurationSchema: util.JSONMap{"type": "object"},
	})
	llmDriver.EXPECT().
		Model(mock.Anything, mock.Anything).
		Return(llm.NewModel("test-model", "test", nil), nil)
	llmDriver.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.MatchedBy(func(options llm.GenerateOptions) bool {
			return len(options.Tools.Tools) == 0 && strings.Contains(options.Prompt.System, "echo__echo_tool")
		})).
		Return(streamOf(
			llm.TextDeltaEvent{Text: `{"tool_call": {"name": "echo__echo_tool", "arguments": {"msg":"hi"}}}`},
		), nil).
		Once()
	llmDriver.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.MatchedBy(func(options llm.GenerateOptions) bool {
			// The tool call/result turns from the previous attempt must survive
			// as plain text the fallback model can read - not as the structured
			// tool_calls/"tool"-role shape its chat template doesn't understand
			// (see rewriteMessagesForFallback) - or it never learns the tool
			// already ran and just calls it again forever.
			for _, m := range options.Prompt.Messages {
				if m.Role == msg.RoleTool {
					return false
				}
				if len(m.ToolCalls) > 0 {
					return false
				}
			}
			var sawToolRequest, sawToolResult bool
			for _, m := range options.Prompt.Messages {
				if strings.Contains(m.Content, "echo__echo_tool") && m.Role == msg.RoleAssistant {
					sawToolRequest = true
				}
				if strings.Contains(m.Content, "echo:hi") {
					sawToolResult = true
				}
			}
			return sawToolRequest && sawToolResult
		})).
		Return(streamOf(llm.TextDeltaEvent{Text: "done"}), nil).
		Once()

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", llmDriver)

	var gotArguments map[string]any
	toolReg := registry.New[tool.Driver]()
	toolReg.MustRegister("echo", tool.CreateDriver(
		tool.WithName("Echo"),
		tool.WithDescription("Echo"),
		tool.WithToolSet(tool.NewToolset(tool.NewTool("echo_tool", "Echoes input", nil))),
		tool.WithAPI(tool.API{
			Call: func(_ context.Context, call tool.Call) (tool.Result, error) {
				gotArguments = call.Arguments
				return tool.Result{Tool: call.Tool, Arguments: call.Arguments, Result: "echo:hi"}, nil
			},
		}),
	))

	engine := runtime.NewEngine(
		runtime.NewDriverManager(llmReg),
		runtime.NewToolManager(runtime.NewDriverManager(toolReg)),
		logger.NewLogger(logger.LevelError),
	)

	agent := runtime.NewAgent("agent", "agent", "be helpful", []runtime.LLMConfig{
		{Key: "primary", Config: util.JSONMap{"model": "test-model"}},
	})
	execution, err := runtime.NewExecution(
		runtime.WithAgent(agent),
		runtime.WithMessageFromUser("hello"),
	)
	require.NoError(t, err)

	err = engine.Start(context.Background(), execution)
	require.NoError(t, err)

	require.Equal(t, runtime.ExecutionStatusCompleted, execution.Status)
	require.Equal(t, map[string]any{"msg": "hi"}, gotArguments)

	assistantWithToolCall := execution.Messages[1]
	require.Equal(t, msg.RoleAssistant, assistantWithToolCall.Role)
	require.Len(t, assistantWithToolCall.ToolCalls, 1)
	require.Equal(t, "echo__echo_tool", assistantWithToolCall.ToolCalls[0].Tool)
	require.Equal(t, map[string]any{"msg": "hi"}, assistantWithToolCall.ToolCalls[0].Arguments)

	toolResult := execution.Messages[2]
	require.Equal(t, msg.RoleTool, toolResult.Role)
	require.Equal(t, "echo:hi", toolResult.Content)

	final := execution.Messages[3]
	require.Equal(t, msg.RoleAssistant, final.Role)
	require.Equal(t, "done", final.Content)
	require.True(t, final.IsFinal())
}
