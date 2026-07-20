package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	agent "github.com/usesnipet/snipet/internal/module/agent"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
)

func newTestService(
	agentRepo repository.IAgentRepository,
) *agent.Service {
	return agent.NewService(agentRepo)
}

func assertAppError(t *testing.T, err error, statusCode int, message string) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, statusCode, appErr.StatusCode)
	assert.Equal(t, message, appErr.Message)
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

	config := model.AgentConfiguration{LLM: model.LLMConfig{Key: "gpt-4"}}
	var stored *model.Agent

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, a *model.Agent) {
			stored = a
			a.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(agentRepo, nil)

	result, err := svc.Create(context.Background(), agent.CreateAgentDTO{
		Name:          "Support Agent",
		Description:   "Handles support tickets",
		Configuration: config,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Support Agent", result.Name)
	assert.Equal(t, "Handles support tickets", result.Description)
	assert.Equal(t, config, result.Configuration)

	require.NotNil(t, stored)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Description, stored.Description)
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
		Name:          "Agent",
		Configuration: model.AgentConfiguration{},
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
			assert.Empty(t, updates.Configuration.LLM)
		}).
		Return(nil)

	svc := newTestService(agentRepo, nil)

	err := svc.Update(context.Background(), id, agent.UpdateAgentDTO{
		Name:        &newName,
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestUpdateDelegatesConfigurationToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	config := model.AgentConfiguration{LLM: model.LLMConfig{Key: "claude"}}

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, _ string, updates *model.Agent) {
			assert.Equal(t, config, updates.Configuration)
		}).
		Return(nil)

	svc := newTestService(agentRepo, nil)

	err := svc.Update(context.Background(), id, agent.UpdateAgentDTO{
		Configuration: &config,
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
