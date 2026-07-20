package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	agent "github.com/usesnipet/snipet/internal/module/agent"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
)

type stubLLM struct {
	info driver.Info
}

func (s *stubLLM) Info() driver.Info { return s.info }
func (s *stubLLM) TestConnection(context.Context, util.JSONMap) error {
	return nil
}
func (s *stubLLM) Generate(
	context.Context,
	util.JSONMap,
	string,
	[]message.Message,
) (message.Message, error) {
	return message.Message{}, nil
}

func newTestService(
	agentRepo repository.IAgentRepository,
	llms map[string]driver.ILLM,
) *agent.Service {
	llmReg := registry.New[driver.ILLM]()
	for name, llm := range llms {
		llmReg.MustRegister(name, llm)
	}

	return agent.NewService(
		agentRepo,
		nil,
		driver.NewManager(llmReg),
		driver.NewManager(registry.New[driver.ITool]()),
		nil,
		nil,
		logger.NewLogger(logger.LevelError),
	)
}

func TestFilterDelegatesToRepository(t *testing.T) {
	t.Parallel()

	expected := page.NewPaginated([]model.Agent{{Name: "Agent A"}}, 1, 0, 10)
	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Run(func(_ context.Context, opts *filter.Options[model.Agent]) {
			assert.Equal(t, filter.Default[model.Agent]().Take, opts.Take)
		}).
		Return(expected, nil)

	svc := newTestService(agentRepo, nil)

	result, err := svc.Filter(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	expected := &model.Agent{ID: id, Name: "Found"}
	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		FindByID(mock.Anything, id).
		Return(expected, nil)

	svc := newTestService(agentRepo, nil)

	result, err := svc.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresAgentAndReturnsIt(t *testing.T) {
	t.Parallel()

	llms := []agent.LLMConfigDTO{{Key: "gpt-4", Config: util.JSONMap{}}}
	var stored *model.Agent

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, a *model.Agent) {
			stored = a
			a.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(agentRepo, map[string]driver.ILLM{
		"gpt-4": &stubLLM{info: driver.Info{
			Name:                "gpt-4",
			Description:         "test",
			ConfigurationSchema: util.JSONMap{"type": "object"},
		}},
	})

	result, err := svc.Create(context.Background(), agent.CreateAgentDTO{
		Name:         "Support Agent",
		Description:  "Handles support tickets",
		Instructions: "Be helpful",
		LLMs:         llms,
		Tools:        agent.ToolConfigDTO{},
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Support Agent", result.Name)
	assert.Equal(t, "Handles support tickets", result.Description)
	assert.Equal(t, "Be helpful", result.Instructions)
	assert.Equal(t, "gpt-4", result.Configuration.LLMs[0].Key)

	require.NotNil(t, stored)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Description, stored.Description)
	assert.Equal(t, result.Instructions, stored.Instructions)
	assert.Equal(t, result.Configuration, stored.Configuration)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("create failed")
	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expectedErr)

	svc := newTestService(agentRepo, nil)

	_, err := svc.Create(context.Background(), agent.CreateAgentDTO{
		Name:  "Agent",
		LLMs:  nil,
		Tools: agent.ToolConfigDTO{},
	})
	require.ErrorIs(t, err, expectedErr)
}

func TestUpdateDelegatesPartialFieldsToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	newName := "Updated Name"
	newDescription := "Updated description"

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, gotID string, updates *model.Agent) {
			assert.Equal(t, id, gotID)
			assert.Equal(t, newName, updates.Name)
			assert.Equal(t, newDescription, updates.Description)
			assert.Empty(t, updates.Configuration.LLMs)
			assert.Empty(t, updates.Configuration.Tools)
		}).
		Return(nil)

	svc := newTestService(agentRepo, nil)

	err := svc.Update(context.Background(), id, agent.UpdateAgentDTO{
		Name:        &newName,
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestUpdateDelegatesLLMsToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	llms := []agent.LLMConfigDTO{{Key: "claude", Config: util.JSONMap{"model": "sonnet"}}}

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, _ string, updates *model.Agent) {
			require.Len(t, updates.Configuration.LLMs, 1)
			assert.Equal(t, "claude", updates.Configuration.LLMs[0].Key)
			assert.Equal(t, util.JSONMap{"model": "sonnet"}, updates.Configuration.LLMs[0].Config)
		}).
		Return(nil)

	svc := newTestService(agentRepo, map[string]driver.ILLM{
		"claude": &stubLLM{info: driver.Info{
			Name:                "claude",
			Description:         "test",
			ConfigurationSchema: util.JSONMap{"type": "object"},
		}},
	})

	err := svc.Update(context.Background(), id, agent.UpdateAgentDTO{
		LLMs: llms,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		DeleteByID(mock.Anything, id).
		Return(nil)

	svc := newTestService(agentRepo, nil)

	err := svc.DeleteByID(context.Background(), id)
	require.NoError(t, err)
}
