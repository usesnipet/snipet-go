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
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
	"github.com/usesnipet/snipet/pkg/driver/tool"
)

type capturingLLM struct {
	key      string
	generate func(ctx context.Context, messages []message.Message) (message.Message, error)
}

func (f *capturingLLM) Info() driver.Info {
	return driver.Info{
		Name:                f.key,
		Description:         f.key,
		ConfigurationSchema: util.JSONMap{"type": "object"},
	}
}

func (f *capturingLLM) TestConnection(context.Context, util.JSONMap) error { return nil }

func (f *capturingLLM) Generate(
	ctx context.Context,
	_ util.JSONMap,
	_ string,
	messages []message.Message,
) (message.Message, error) {
	return f.generate(ctx, messages)
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
				Configuration: util.JSONMap{},
			},
		}},
	}
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
		Run(func(_ context.Context, _ string, msgs []model.ExecutionMessage) {
			persisted = append(persisted, msgs...)
		}).
		Return(nil)

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", &capturingLLM{
		key: "primary",
		generate: func(_ context.Context, _ []message.Message) (message.Message, error) {
			return message.Message{Role: message.MessageRoleFinal, Content: "done"}, nil
		},
	})

	svc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		runtime.NewEngine(
			runtime.NewDriverManager(registry.New[tool.Driver]()),
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
	}, nil)
	require.NoError(t, err)

	require.NotNil(t, created)
	assert.Nil(t, created.SessionID)
	assert.Equal(t, agentID, created.AgentID)

	roles := make([]message.Role, 0, len(persisted))
	for _, m := range persisted {
		roles = append(roles, m.Role)
	}
	assert.Contains(t, roles, message.MessageRoleUser)
	assert.Contains(t, roles, message.MessageRoleFinal)
}

func TestRunWithSessionLoadsHistoryAndSkipsRePersistingIt(t *testing.T) {
	t.Parallel()

	agentID := uuid.New().String()
	sessionID := uuid.New().String()
	historyMsgID := uuid.New().String()

	var created *model.Execution
	var persisted []model.ExecutionMessage
	var llmSaw []message.Message

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
			ID:      historyMsgID,
			Role:    message.MessageRoleUser,
			Content: "earlier",
		}}, nil)
	messageRepo.EXPECT().
		CreateInExecution(mock.Anything, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ string, msgs []model.ExecutionMessage) {
			persisted = append(persisted, msgs...)
		}).
		Return(nil)

	llmReg := registry.New[llm.Driver]()
	llmReg.MustRegister("primary", &capturingLLM{
		key: "primary",
		generate: func(_ context.Context, messages []message.Message) (message.Message, error) {
			llmSaw = append([]message.Message{}, messages...)
			return message.Message{Role: message.MessageRoleFinal, Content: "done"}, nil
		},
	})

	svc := agent.NewService(
		agentRepo,
		mocks.NewMockILLMRepository(t),
		mocks.NewMockITxManager(t),
		runtime.NewEngine(
			runtime.NewDriverManager(registry.New[tool.Driver]()),
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
	}, nil)
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
