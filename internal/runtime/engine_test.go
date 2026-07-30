package runtime_test

import (
	"context"
	"errors"
	"testing"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

type fakeLLM struct {
	key      string
	generate func(ctx context.Context, messages []llm.Message) (llm.Message, error)
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
	messages []llm.Message,
) (llm.Message, error) {
	return f.generate(ctx, messages)
}

func newTestEngine(t *testing.T, llms map[string]*fakeLLM) *runtime.Engine {
	t.Helper()

	llmReg := registry.New[llm.Driver]()
	for key, instance := range llms {
		instance.key = key
		llmReg.MustRegister(key, instance)
	}

	return runtime.NewEngine(
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

func TestEngine_HappyPath(t *testing.T) {
	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []llm.Message) (llm.Message, error) {
				return llm.Message{Role: llm.MessageRoleAssistant, Content: "done"}, nil
			},
		},
	})

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
		switch event := event.(type) {
		case runtime.ExecutionStatusChangedEvent:
			if e := event; e.Status == runtime.ExecutionStatusCompleted {
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

func TestEngine_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	engine := newTestEngine(t, map[string]*fakeLLM{
		"primary": {
			generate: func(ctx context.Context, messages []llm.Message) (llm.Message, error) {
				t.Fatal("LLM should not be called after cancel")
				return llm.Message{}, nil
			},
		},
	})

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
			generate: func(ctx context.Context, messages []llm.Message) (llm.Message, error) {
				return llm.Message{Role: llm.MessageRoleAssistant, Content: "done"}, nil
			},
		},
	})

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
			generate: func(ctx context.Context, messages []llm.Message) (llm.Message, error) {
				return llm.Message{}, errors.New("bad llm")
			},
		},
		"good": {
			generate: func(ctx context.Context, messages []llm.Message) (llm.Message, error) {
				return llm.Message{Role: llm.MessageRoleAssistant, Content: "ok"}, nil
			},
		},
	})

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
			generate: func(ctx context.Context, messages []llm.Message) (llm.Message, error) {
				t.Fatal("should not generate")
				return llm.Message{}, nil
			},
		},
	})

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
