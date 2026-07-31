package runtime_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	llmmocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/msg"
)

func driverInfo(key string) driver.Info {
	return driver.Info{
		Name:                key,
		Description:         key,
		ConfigurationSchema: util.JSONMap{"type": "object"},
	}
}

func newTestEngine(t *testing.T, llms map[string]*llmmocks.MockDriver) *runtime.Engine {
	t.Helper()

	llmReg := registry.New[llm.Driver]()
	for key, instance := range llms {
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
	primary := llmmocks.NewMockDriver(t)
	primary.EXPECT().Info().Return(driverInfo("primary"))
	primary.EXPECT().
		Generate(mock.Anything, mock.Anything, mock.Anything).
		Return(msg.Message{Role: msg.RoleAssistant, Content: "done"}, nil)

	engine := newTestEngine(t, map[string]*llmmocks.MockDriver{
		"primary": primary,
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

	primary := llmmocks.NewMockDriver(t)
	primary.EXPECT().Info().Return(driverInfo("primary"))

	engine := newTestEngine(t, map[string]*llmmocks.MockDriver{
		"primary": primary,
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
	primary := llmmocks.NewMockDriver(t)
	primary.EXPECT().Info().Return(driverInfo("primary"))
	primary.EXPECT().
		Generate(mock.Anything, mock.Anything, mock.Anything).
		Return(msg.Message{Role: msg.RoleAssistant, Content: "done"}, nil).
		Maybe()

	engine := newTestEngine(t, map[string]*llmmocks.MockDriver{
		"primary": primary,
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
	bad := llmmocks.NewMockDriver(t)
	bad.EXPECT().Info().Return(driverInfo("bad"))
	bad.EXPECT().
		Generate(mock.Anything, mock.Anything, mock.Anything).
		Return(msg.Message{}, errors.New("bad llm"))

	good := llmmocks.NewMockDriver(t)
	good.EXPECT().Info().Return(driverInfo("good"))
	good.EXPECT().
		Generate(mock.Anything, mock.Anything, mock.Anything).
		Return(msg.Message{Role: msg.RoleAssistant, Content: "ok"}, nil)

	engine := newTestEngine(t, map[string]*llmmocks.MockDriver{
		"bad":  bad,
		"good": good,
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
	primary := llmmocks.NewMockDriver(t)

	engine := newTestEngine(t, map[string]*llmmocks.MockDriver{
		"primary": primary,
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
