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
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	llmmocks "github.com/usesnipet/snipet/pkg/driver/llm/mocks"
	"github.com/usesnipet/snipet/pkg/msg"
)

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
				Configuration: util.JSONMap{},
			},
		}},
	}
}

func newPrimaryLLM(t *testing.T) *llmmocks.MockDriver {
	t.Helper()

	d := llmmocks.NewMockDriver(t)
	d.EXPECT().Info().Return(driver.Info{
		Name:                "primary",
		Description:         "primary",
		ConfigurationSchema: util.JSONMap{"type": "object"},
	})
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
		Generate(mock.Anything, mock.Anything, mock.Anything).
		Return(msg.Message{Role: msg.RoleAssistant, Content: "done"}, nil)

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", primary)

	svc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		runtime.NewEngine(
			runtime.NewDriverManager(llmReg),
			logger.NewLogger(logger.LevelError),
		),
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
		Generate(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, _ util.JSONMap, prompt llm.Prompt) (msg.Message, error) {
			llmSaw = append([]msg.Message{}, prompt.Messages...)
			return msg.NewMessage(msg.RoleAssistant, "done"), nil
		})

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", primary)

	svc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		runtime.NewEngine(
			runtime.NewDriverManager(llmReg),
			logger.NewLogger(logger.LevelError),
		),
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
