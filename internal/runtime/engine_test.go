package runtime_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type fakeLLM struct {
	key      string
	generate func(ctx context.Context, messages []message.Message) (message.Message, error)
}

func (f *fakeLLM) Info() driver.Info {
	return driver.Info{
		Name:                f.key,
		Description:         f.key,
		ConfigurationSchema: util.JSONMap{"type": "object"},
	}
}

func (f *fakeLLM) TestConnection(ctx context.Context, config util.JSONMap) error {
	return nil
}

func (f *fakeLLM) Generate(
	ctx context.Context,
	config util.JSONMap,
	instructions string,
	messages []message.Message,
) (message.Message, error) {
	return f.generate(ctx, messages)
}

type fakeTool struct {
	key string
}

func (f *fakeTool) Info() driver.Info {
	return driver.Info{
		Name:                f.key,
		Description:         f.key,
		ConfigurationSchema: util.JSONMap{"type": "object"},
	}
}

func (f *fakeTool) TestConnection(ctx context.Context, config util.JSONMap) error {
	return nil
}

func (f *fakeTool) Execute(ctx context.Context, config util.JSONMap, call tool.Call) tool.Result {
	return tool.Result{Key: call.Key, Success: true, Output: "ok"}
}

func newTestEngine(t *testing.T, llms map[string]*fakeLLM, tools map[string]*fakeTool) *runtime.Engine {
	t.Helper()

	llmReg := registry.New[llm.Driver]()
	for key, llm := range llms {
		llm.key = key
		llmReg.MustRegister(key, llm)
	}

	toolReg := registry.New[tool.Driver]()
	for key, tl := range tools {
		tl.key = key
		toolReg.MustRegister(key, tl)
	}

	return runtime.NewEngine(
		runtime.NewDriverManager(toolReg),
		runtime.NewDriverManager(llmReg),
		logger.NewLogger(logger.LevelError),
	)
}

func collectEvents(events *[]runtime.IEvent) func(runtime.IEvent) error {
	return func(event runtime.IEvent) error {
		*events = append(*events, event)
		return nil
	}
}

func TestEngine_HappyPathFinal(t *testing.T) {
	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				return message.Message{Role: message.MessageRoleFinal, Content: "done"}, nil
			},
		},
	}, nil)

	var events []runtime.IEvent
	err := engine.Start(context.Background(), runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs: []runtime.LLMConfig{{Key: "primary", Config: util.JSONMap{}}},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
		},
		OnEvent: collectEvents(&events),
	})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	var sawCompleted, sawMessage bool
	for _, event := range events {
		switch e := event.(type) {
		case runtime.ExecutionStatusChangedEvent:
			if e.Status == runtime.ExecutionStatusCompleted {
				sawCompleted = true
			}
		case runtime.ExecutionMessageAddedEvent:
			sawMessage = true
		}
	}
	if !sawCompleted {
		t.Fatal("expected completed status event")
	}
	if !sawMessage {
		t.Fatal("expected message added event")
	}
}

func TestEngine_ToolCallSequence(t *testing.T) {
	var callCount atomic.Int32
	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				if callCount.Add(1) == 1 {
					return message.Message{
						Role:    message.MessageRoleAssistant,
						Content: "calling",
						ToolCalls: []tool.Call{
							{Key: "echo", Input: "x"},
						},
					}, nil
				}
				return message.Message{Role: message.MessageRoleFinal, Content: "done"}, nil
			},
		},
	}, map[string]*fakeTool{"echo": {}})

	var events []runtime.IEvent
	err := engine.Start(context.Background(), runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs:  []runtime.LLMConfig{{Key: "primary", Config: util.JSONMap{}}},
			Tools: runtime.ToolConfig{"echo": util.JSONMap{}},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
			runtime.WithMaxTurns(5),
		},
		OnEvent: collectEvents(&events),
	})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	last := -1
	for _, event := range events {
		if e, ok := event.(runtime.ExecutionMessageAddedEvent); ok {
			for _, msg := range e.Messages {
				if msg.ID == "" {
					t.Fatal("message missing ID")
				}
				if msg.Sequence < last {
					t.Fatalf("sequence not monotonic: got %d after %d", msg.Sequence, last)
				}
				last = msg.Sequence
			}
		}
	}
}

func TestEngine_MaxTurns(t *testing.T) {
	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				return message.Message{
					Role:      message.MessageRoleAssistant,
					Content:   "loop",
					ToolCalls: []tool.Call{{Key: "echo", Input: "x"}},
				}, nil
			},
		},
	}, map[string]*fakeTool{"echo": {}})

	var events []runtime.IEvent
	err := engine.Start(context.Background(), runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs:  []runtime.LLMConfig{{Key: "primary", Config: util.JSONMap{}}},
			Tools: runtime.ToolConfig{"echo": util.JSONMap{}},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
			runtime.WithMaxTurns(1),
		},
		OnEvent: collectEvents(&events),
	})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	var sawMaxTurns bool
	for _, event := range events {
		if e, ok := event.(runtime.ExecutionStatusChangedEvent); ok {
			if e.Status == runtime.ExecutionStatusMaxTurns {
				sawMaxTurns = true
				if e.ErrorMessage == "" {
					t.Fatal("expected error_message on max_turns")
				}
			}
		}
	}
	if !sawMaxTurns {
		t.Fatal("expected max_turns status event")
	}
}

func TestEngine_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				t.Fatal("LLM should not be called after cancel")
				return message.Message{}, nil
			},
		},
	}, nil)

	var events []runtime.IEvent
	err := engine.Start(ctx, runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs: []runtime.LLMConfig{{Key: "primary", Config: util.JSONMap{}}},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
		},
		OnEvent: collectEvents(&events),
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	var sawCancelled bool
	for _, event := range events {
		if e, ok := event.(runtime.ExecutionStatusChangedEvent); ok {
			if e.Status == runtime.ExecutionStatusCancelled {
				sawCancelled = true
			}
		}
	}
	if !sawCancelled {
		t.Fatal("expected cancelled status event")
	}
}

func TestEngine_OnEventError(t *testing.T) {
	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				return message.Message{Role: message.MessageRoleFinal, Content: "done"}, nil
			},
		},
	}, nil)

	boom := errors.New("persist failed")
	err := engine.Start(context.Background(), runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs: []runtime.LLMConfig{{Key: "primary", Config: util.JSONMap{}}},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
		},
		OnEvent: func(event runtime.IEvent) error {
			if _, ok := event.(runtime.ExecutionStatusChangedEvent); ok {
				return boom
			}
			return nil
		},
	})
	if !errors.Is(err, boom) {
		t.Fatalf("expected persist error, got %v", err)
	}
}

func TestEngine_LLMFallback(t *testing.T) {
	engine := newTestEngine(t, map[string]*fakeLLM{
		"bad": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				return message.Message{}, errors.New("bad llm")
			},
		},
		"good": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				return message.Message{Role: message.MessageRoleFinal, Content: "ok"}, nil
			},
		},
	}, nil)

	var events []runtime.IEvent
	err := engine.Start(context.Background(), runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs: []runtime.LLMConfig{
				{Key: "bad", Config: util.JSONMap{}},
				{Key: "good", Config: util.JSONMap{}},
			},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
		},
		OnEvent: collectEvents(&events),
	})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	var completed bool
	for _, event := range events {
		if e, ok := event.(runtime.ExecutionStatusChangedEvent); ok {
			if e.Status == runtime.ExecutionStatusCompleted {
				completed = true
			}
		}
	}
	if !completed {
		t.Fatal("expected completed after fallback")
	}
}

func TestEngine_ValidateFail(t *testing.T) {
	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []message.Message) (message.Message, error) {
				t.Fatal("should not generate")
				return message.Message{}, nil
			},
		},
	}, nil)

	var events []runtime.IEvent
	err := engine.Start(context.Background(), runtime.StartOptions{
		Agent: runtime.Agent{
			LLMs: []runtime.LLMConfig{{Key: "missing", Config: util.JSONMap{}}},
		},
		ExecutionOptions: []runtime.ExecutionOption{
			runtime.WithMessageFromUser("hi"),
		},
		OnEvent: collectEvents(&events),
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	var sawFailed bool
	for _, event := range events {
		if e, ok := event.(runtime.ExecutionStatusChangedEvent); ok {
			if e.Status == runtime.ExecutionStatusFailed {
				sawFailed = true
				if e.ErrorMessage == "" {
					t.Fatal("expected error_message on failed")
				}
			}
		}
	}
	if !sawFailed {
		t.Fatal("expected failed status event")
	}
}
