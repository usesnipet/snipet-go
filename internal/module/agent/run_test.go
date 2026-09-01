package agent_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	agent "github.com/usesnipet/snipet/internal/module/agent"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	llmmocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
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

func newTestEngine(llmReg *driver.Registry[llm.Driver]) *runtime.Engine {
	return runtime.NewEngine(
		manager.NewDriverManager(llmReg),
		manager.NewToolbox(manager.NewDriverManager(driver.NewRegistry[tool.Driver](logger.NewLogger(logger.LevelError)))),
		logger.NewLogger(logger.LevelError),
	)
}

func agentWithPrimaryLLM(agentID string) *model.Agent {
	llmID := uuid.New().String()
	return &model.Agent{
		ID:   agentID,
		Name: "Agent",
		AgentToLLMs: []model.AgentToLLM{{
			AgentID:  agentID,
			LLMID:    llmID,
			Priority: 0,
			LLM: model.LLM{
				ID:            llmID,
				Name:          "primary",
				Provider:      "primary",
				Configuration: jsonx.JSONMap{},
			},
		}},
	}
}

func newPrimaryLLM(t *testing.T) *llmmocks.MockDriver {
	t.Helper()

	d := llmmocks.NewMockDriver(t)
	d.EXPECT().Info().Return(driver.Info{
		Key:                 "primary",
		Name:                "primary",
		Description:         "primary",
		ConfigurationSchema: jsonx.JSONMap{"type": "object"},
	})
	d.EXPECT().Validate().Return(nil)
	return d
}

func TestRunPlaygroundCreatesExecutionWithoutSession(t *testing.T) {
	t.Parallel()

	agentID := uuid.New().String()
	var created *model.Execution
	var persisted []model.ExecutionMessage

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		FindByID(mock.Anything, agentID).
		Return(agentWithPrimaryLLM(agentID), nil)

	executionRepo := mocks.NewMockIExecutionRepository(t)
	executionRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, e *model.Execution) {
			created = e
			e.ID = uuid.New().String()
		}).
		Return(nil)
	executionRepo.EXPECT().
		UpdateByID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)
	messageRepo.EXPECT().
		CreateInExecution(mock.Anything, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ string, msgs ...model.ExecutionMessage) {
			persisted = append(persisted, msgs...)
		}).
		Return(nil)

	primary := newPrimaryLLM(t)
	primary.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		Return(streamOf(llm.TextDeltaEvent{Text: "done"}), nil)

	llmReg := driver.NewRegistry[llm.Driver](logger.NewLogger(logger.LevelError))
	llmReg.MustRegister(primary, nil)

	svc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		newTestEngine(llmReg),
		executionRepo,
		messageRepo,
		logger.NewLogger(logger.LevelError),
	)

	err := svc.Run(context.Background(), agent.RunInput{
		AgentID: agentID,
		Message: "hello",
	})
	require.NoError(t, err)

	require.NotNil(t, created)
	assert.Nil(t, created.SessionID)
	assert.Equal(t, agentID, created.AgentID)

	roles := make([]msg.Role, 0, len(persisted))
	for _, m := range persisted {
		roles = append(roles, m.Role)
	}
	assert.Contains(t, roles, msg.RoleUser)
	assert.Contains(t, roles, msg.RoleAssistant)
}

func TestRunWithSessionLoadsHistoryAndSkipsRePersistingIt(t *testing.T) {
	t.Parallel()

	agentID := uuid.New().String()
	sessionID := uuid.New().String()
	historyMsgID := uuid.New().String()

	var created *model.Execution
	var persisted []model.ExecutionMessage
	var llmSaw []msg.Message

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		FindByID(mock.Anything, agentID).
		Return(agentWithPrimaryLLM(agentID), nil)

	executionRepo := mocks.NewMockIExecutionRepository(t)
	executionRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, e *model.Execution) {
			created = e
			e.ID = uuid.New().String()
		}).
		Return(nil)
	executionRepo.EXPECT().
		UpdateByID(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	messageRepo := mocks.NewMockIExecutionMessageRepository(t)
	messageRepo.EXPECT().
		ListBySessionID(mock.Anything, sessionID).
		Return([]model.ExecutionMessage{{
			Message: msg.NewMessage(msg.RoleUser, "earlier", msg.WithID(historyMsgID)),
		}}, nil)
	messageRepo.EXPECT().
		CreateInExecution(mock.Anything, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ string, msgs ...model.ExecutionMessage) {
			persisted = append(persisted, msgs...)
		}).
		Return(nil)

	primary := newPrimaryLLM(t)
	primary.EXPECT().
		Stream(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _ jsonx.JSONMap, options llm.GenerateOptions) (llm.StreamIterator, error) {
			llmSaw = append([]msg.Message{}, options.Prompt.Messages...)
			return streamOf(llm.TextDeltaEvent{Text: "done"}), nil
		})

	llmReg := driver.NewRegistry[llm.Driver](logger.NewLogger(logger.LevelError))
	llmReg.MustRegister(primary, nil)

	svc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		newTestEngine(llmReg),
		executionRepo,
		messageRepo,
		logger.NewLogger(logger.LevelError),
	)

	err := svc.Run(context.Background(), agent.RunInput{
		AgentID:   agentID,
		SessionID: &sessionID,
		Message:   "follow up",
	})
	require.NoError(t, err)

	require.NotNil(t, created)
	require.NotNil(t, created.SessionID)
	assert.Equal(t, sessionID, *created.SessionID)

	require.GreaterOrEqual(t, len(llmSaw), 2)
	assert.Equal(t, "earlier", llmSaw[0].Content)
	assert.Equal(t, "follow up", llmSaw[1].Content)

	for _, m := range persisted {
		assert.NotEqual(t, historyMsgID, m.ID, "history must not be re-persisted")
	}
	contents := make([]string, 0, len(persisted))
	for _, m := range persisted {
		contents = append(contents, m.Content)
	}
	assert.Contains(t, contents, "follow up")
	assert.Contains(t, contents, "done")
	assert.NotContains(t, contents, "earlier")
}
