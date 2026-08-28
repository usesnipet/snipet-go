package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	agent "github.com/usesnipet/snipet/internal/module/agent"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
)

func newTestService(
	t *testing.T,
	agentRepo repository.IAgentRepository,
	llmRepo repository.ILLMRepository,
	txManager repository.ITxManager,
) *agent.Service {
	t.Helper()
	if llmRepo == nil {
		llmRepo = mocks.NewMockILLMRepository(t)
	}
	if txManager == nil {
		tx := mocks.NewMockITxManager(t)
		tx.EXPECT().
			WithTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			}).
			Maybe()
		txManager = tx
	}

	return agent.NewService(
		agentRepo,
		llmRepo,
		txManager,
		nil,
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
		Return(expected, nil)

	svc := newTestService(t, agentRepo, nil, nil)

	result, err := svc.Filter(context.Background(), nil)
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

	svc := newTestService(t, agentRepo, nil, nil)

	result, err := svc.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresAgentAndReplacesLLMs(t *testing.T) {
	t.Parallel()

	llmID := uuid.New().String()
	agentID := uuid.New().String()
	var stored *model.Agent

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, a *model.Agent) {
			stored = a
			a.ID = agentID
		}).
		Return(nil)
	agentRepo.EXPECT().
		ReplaceLLMs(mock.Anything, agentID, []string{llmID}).
		Return(nil)
	agentRepo.EXPECT().
		FindByID(mock.Anything, agentID).
		Return(&model.Agent{
			ID:   agentID,
			Name: "Support Agent",
			AgentToLLMs: []model.AgentToLLM{{
				AgentID:  agentID,
				LLMID:    llmID,
				Priority: 0,
				LLM:      model.LLM{ID: llmID, Name: "gpt prod", Provider: "gpt-4"},
			}},
		}, nil)

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.LLM{{ID: llmID}}, 1, 0, 1), nil)

	txManager := mocks.NewMockITxManager(t)
	txManager.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	svc := newTestService(t, agentRepo, llmRepo, txManager)

	result, err := svc.Create(context.Background(), agent.CreateAgentDTO{
		Name:         "Support Agent",
		Description:  "Handles support tickets",
		Instructions: "Be helpful",
		LLMIDs:       []string{llmID},
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Support Agent", result.Name)
	require.Len(t, result.AgentToLLMs, 1)
	assert.Equal(t, llmID, result.AgentToLLMs[0].LLMID)

	require.NotNil(t, stored)
	assert.Equal(t, "Support Agent", stored.Name)
	assert.Equal(t, "Handles support tickets", stored.Description)
	assert.Equal(t, "Be helpful", stored.Instructions)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	llmID := uuid.New().String()
	expectedErr := errors.New("create failed")

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expectedErr)

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.LLM{{ID: llmID}}, 1, 0, 1), nil)

	txManager := mocks.NewMockITxManager(t)
	txManager.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	svc := newTestService(t, agentRepo, llmRepo, txManager)

	_, err := svc.Create(context.Background(), agent.CreateAgentDTO{
		Name:   "Agent",
		LLMIDs: []string{llmID},
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
		}).
		Return(nil)

	txManager := mocks.NewMockITxManager(t)
	txManager.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	svc := newTestService(t, agentRepo, nil, txManager)

	err := svc.Update(context.Background(), id, agent.UpdateAgentDTO{
		Name:        &newName,
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestUpdateReplacesLLMs(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	llmID := uuid.New().String()

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		ReplaceLLMs(mock.Anything, id, []string{llmID}).
		Return(nil)

	llmRepo := mocks.NewMockILLMRepository(t)
	llmRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.LLM{{ID: llmID}}, 1, 0, 1), nil)

	txManager := mocks.NewMockITxManager(t)
	txManager.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	svc := newTestService(t, agentRepo, llmRepo, txManager)

	err := svc.Update(context.Background(), id, agent.UpdateAgentDTO{
		LLMIDs: []string{llmID},
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

	svc := newTestService(t, agentRepo, nil, nil)

	err := svc.DeleteByID(context.Background(), id)
	require.NoError(t, err)
}
