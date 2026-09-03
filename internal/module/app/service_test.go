package app_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	appmodule "github.com/usesnipet/snipet/internal/module/app"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository/mocks"
)

func assertAppError(t *testing.T, err error, statusCode int) {
	t.Helper()
	var appErr *apperr.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, statusCode, appErr.StatusCode)
}

// passthroughTx returns a MockITxManager that just runs the callback inline.
func passthroughTx(t *testing.T) *mocks.MockITxManager {
	t.Helper()
	tx := mocks.NewMockITxManager(t)
	tx.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()
	return tx
}

func newAppService(
	t *testing.T,
	appRepo *mocks.MockIAppRepository,
	agentRepo *mocks.MockIAgentRepository,
) *appmodule.Service {
	t.Helper()
	if agentRepo == nil {
		agentRepo = mocks.NewMockIAgentRepository(t)
	}
	return appmodule.NewService(
		appRepo,
		agentRepo,
		passthroughTx(t),
		auth.NewAPIKeyGenerator(),
		auth.NewKeyHasher(),
		logger.NewLogger(logger.LevelError),
	)
}

func TestFindByCodeDelegatesToRepository(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme"}, nil)

	svc := newAppService(t, appRepo, nil)

	data, err := svc.FindByCode(context.Background(), "acme")
	require.NoError(t, err)
	require.Equal(t, "acme", data.Code)
}

func TestFindPublicByCodeReturnsPublicFields(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme", Name: "Acme", Description: "desc", Public: true}, nil)

	svc := newAppService(t, appRepo, nil)

	data, err := svc.FindPublicByCode(context.Background(), "acme")
	require.NoError(t, err)
	require.Equal(t, &appmodule.PublicAppDTO{Code: "acme", Name: "Acme", Description: "desc"}, data)
}

func TestFindPublicByCodeIncludesLinkedAgents(t *testing.T) {
	t.Parallel()

	links := []model.AppToAgent{
		{AppID: "app-1", AgentID: "agent-1", Agent: &model.Agent{ID: "agent-1", Name: "Support"}},
	}
	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{
			ID: "app-1", Code: "acme", Name: "Acme", Description: "desc",
			Public: true, AppToAgents: links,
		}, nil)

	svc := newAppService(t, appRepo, nil)

	data, err := svc.FindPublicByCode(context.Background(), "acme")
	require.NoError(t, err)
	require.Equal(t, links, data.AppToAgents)
	require.Equal(t, "Support", data.AppToAgents[0].Agent.Name)
}

func TestFindPublicByCodeRejectsNonPublicApp(t *testing.T) {
	t.Parallel()

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{Code: "acme", Public: false}, nil)

	svc := newAppService(t, appRepo, nil)

	_, err := svc.FindPublicByCode(context.Background(), "acme")
	assertAppError(t, err, http.StatusNotFound)
}

func TestLinkAgentsRejectsUnknownAgentID(t *testing.T) {
	t.Parallel()

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.Agent{{ID: "agent-1"}}, 1, 0, 0), nil)

	svc := newAppService(t, mocks.NewMockIAppRepository(t), agentRepo)

	err := svc.LinkAgents(context.Background(), "acme", []string{"agent-1", "agent-missing"})
	assertAppError(t, err, http.StatusBadRequest)
}

func TestLinkAgentsRejectsDuplicateAgentID(t *testing.T) {
	t.Parallel()

	svc := newAppService(t, mocks.NewMockIAppRepository(t), mocks.NewMockIAgentRepository(t))

	err := svc.LinkAgents(context.Background(), "acme", []string{"agent-1", "agent-1"})
	assertAppError(t, err, http.StatusBadRequest)
}

func TestLinkAgentsReplacesTheWholeSet(t *testing.T) {
	t.Parallel()

	agentRepo := mocks.NewMockIAgentRepository(t)
	agentRepo.EXPECT().
		Filter(mock.Anything, mock.Anything).
		Return(page.NewPaginated([]model.Agent{{ID: "agent-1"}, {ID: "agent-2"}}, 2, 0, 0), nil)

	appRepo := mocks.NewMockIAppRepository(t)
	appRepo.EXPECT().
		FindByCode(mock.Anything, "acme").
		Return(&model.App{ID: "app-1", Code: "acme"}, nil)
	appRepo.EXPECT().
		ReplaceAgents(mock.Anything, "app-1", []string{"agent-1", "agent-2"}).
		Return(nil)

	svc := newAppService(t, appRepo, agentRepo)

	require.NoError(t, svc.LinkAgents(context.Background(), "acme", []string{"agent-1", "agent-2"}))
}
