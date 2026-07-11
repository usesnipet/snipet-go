package knowledge_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	knowledge "github.com/usesnipet/snipet/internal/module/knowledge"
	"github.com/usesnipet/snipet/internal/page"
	queuemocks "github.com/usesnipet/snipet/internal/queue/mocks"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime"
	runtimemocks "github.com/usesnipet/snipet/internal/runtime/mocks"
	"github.com/usesnipet/snipet/internal/util"
)

var testConfigSchema = util.JSONMap{
	"type": "object",
	"properties": util.JSONMap{
		"index": util.JSONMap{"type": "string"},
	},
	"required": []any{"index"},
}

func newPassthroughTxManager(t *testing.T) *mocks.MockITxManager {
	t.Helper()

	tx := mocks.NewMockITxManager(t)
	tx.EXPECT().
		WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()
	return tx
}

func newNoopJobQueue(t *testing.T) *queuemocks.MockIJobQueue {
	t.Helper()

	q := queuemocks.NewMockIJobQueue(t)
	q.EXPECT().Push(mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil).Maybe()
	q.EXPECT().JobGet(mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	return q
}

func newTestService(
	t *testing.T,
	repo repository.IKnowledgeRepository,
	drivers map[string]*runtimemocks.MockISourceDriver,
	opts ...func(*testServiceOptions),
) *knowledge.Service {
	t.Helper()

	options := testServiceOptions{
		txManager:   newPassthroughTxManager(t),
		riverClient: newNoopJobQueue(t),
	}
	for _, opt := range opts {
		opt(&options)
	}

	registry := runtime.NewRegistry[runtime.ISourceDriver]()
	for name, driver := range drivers {
		registry.MustRegister(name, driver)
	}
	sourceManager := runtime.NewSourceManager(registry)
	return knowledge.NewService(
		options.txManager,
		repo,
		mocks.NewMockIKnowledgeItemRepository(t),
		sourceManager,
		options.riverClient,
	)
}

type testServiceOptions struct {
	txManager   repository.ITxManager
	riverClient *queuemocks.MockIJobQueue
}

func expectSuccessfulConnection(driver *runtimemocks.MockISourceDriver, config util.JSONMap) {
	driver.EXPECT().GetConfigurationSchema(mock.Anything).Return(testConfigSchema, nil).Once()
	driver.EXPECT().TestConnection(mock.Anything, config).Return(nil).Once()
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

	opts := filter.Default[model.Knowledge]()
	expected := page.NewPaginated([]model.Knowledge{{Name: "Knowledge A"}}, 1, 0, 10)
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Filter(mock.Anything, opts).
		Return(expected, nil)

	svc := newTestService(t, repo, nil)

	result, err := svc.Filter(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	expected := &model.Knowledge{ID: id, Name: "Found"}
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		FindByID(mock.Anything, id).
		Return(expected, nil)

	svc := newTestService(t, repo, nil)

	result, err := svc.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresKnowledgeAndReturnsIt(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	var stored *model.Knowledge

	driver := runtimemocks.NewMockISourceDriver(t)
	expectSuccessfulConnection(driver, config)

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, k *model.Knowledge) {
			stored = k
			k.ID = uuid.New().String()
		}).
		Return(nil)
	repo.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, id string) (*model.Knowledge, error) {
			return &model.Knowledge{ID: id}, nil
		})

	riverClient := queuemocks.NewMockIJobQueue(t)
	riverClient.EXPECT().Push(mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil).Once()

	svc := newTestService(t, repo, map[string]*runtimemocks.MockISourceDriver{"pinecone": driver}, func(o *testServiceOptions) {
		o.riverClient = riverClient
	})

	result, _, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Product Docs",
		Description:   "Internal documentation",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Product Docs", result.Name)
	assert.Equal(t, "Internal documentation", result.Description)
	assert.Equal(t, "pinecone", result.Driver)
	assert.Equal(t, config, result.Configuration)

	require.NotNil(t, stored)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Description, stored.Description)
	assert.Equal(t, result.Driver, stored.Driver)
	assert.Equal(t, result.Configuration, stored.Configuration)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	expectedErr := errors.New("create failed")

	driver := runtimemocks.NewMockISourceDriver(t)
	expectSuccessfulConnection(driver, config)

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expectedErr)

	svc := newTestService(t, repo, map[string]*runtimemocks.MockISourceDriver{"pinecone": driver})

	_, _, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Knowledge",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.ErrorIs(t, err, expectedErr)
}

func TestCreateReturnsBadRequestWhenConnectionFails(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	connectionErr := errors.New("connection refused")

	driver := runtimemocks.NewMockISourceDriver(t)
	driver.EXPECT().GetConfigurationSchema(mock.Anything).Return(testConfigSchema, nil).Once()
	driver.EXPECT().TestConnection(mock.Anything, config).Return(connectionErr).Once()

	repo := mocks.NewMockIKnowledgeRepository(t)

	svc := newTestService(t, repo, map[string]*runtimemocks.MockISourceDriver{"pinecone": driver})

	_, _, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Knowledge",
		Driver:        "pinecone",
		Configuration: config,
	})
	assertAppError(t, err, http.StatusBadRequest, "bad request: "+runtime.ErrConnectionFailed.Error())
}

func TestUpdateDelegatesPartialFieldsToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	newName := "Updated Name"
	newDescription := "Updated description"

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, gotID string, updates *model.Knowledge) {
			assert.Equal(t, id, gotID)
			assert.Equal(t, newName, updates.Name)
			assert.Equal(t, newDescription, updates.Description)
			assert.Empty(t, updates.Driver)
			assert.Nil(t, updates.Configuration)
		}).
		Return(nil)

	svc := newTestService(t, repo, nil)

	err := svc.Update(context.Background(), id, knowledge.UpdateKnowledgeDTO{
		Name:        &newName,
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestUpdateDelegatesDescriptionToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	newDescription := "Updated description"

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		UpdateByID(mock.Anything, id, mock.Anything).
		Run(func(_ context.Context, gotID string, updates *model.Knowledge) {
			assert.Equal(t, id, gotID)
			assert.Empty(t, updates.Name)
			assert.Equal(t, newDescription, updates.Description)
		}).
		Return(nil)

	svc := newTestService(t, repo, nil)

	err := svc.Update(context.Background(), id, knowledge.UpdateKnowledgeDTO{
		Description: &newDescription,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		DeleteByID(mock.Anything, id).
		Return(nil)

	svc := newTestService(t, repo, nil)

	err := svc.DeleteByID(context.Background(), id)
	require.NoError(t, err)
}

func TestTestConnectionSucceedsWhenDriverValidatesAndConnects(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}

	driver := runtimemocks.NewMockISourceDriver(t)
	expectSuccessfulConnection(driver, config)

	svc := newTestService(t, mocks.NewMockIKnowledgeRepository(t), map[string]*runtimemocks.MockISourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", config)
	require.NoError(t, err)
}

func TestTestConnectionReturnsNotFoundForUnknownDriver(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, mocks.NewMockIKnowledgeRepository(t), nil)

	err := svc.TestConnection(context.Background(), "unknown", util.JSONMap{"index": "docs"})
	assertAppError(t, err, http.StatusNotFound, runtime.ErrDriverNotFound.Error())
}

func TestTestConnectionReturnsBadRequestWhenSchemaFetchFails(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	schemaErr := errors.New("schema unavailable")

	driver := runtimemocks.NewMockISourceDriver(t)
	driver.EXPECT().GetConfigurationSchema(mock.Anything).Return(nil, schemaErr).Once()

	svc := newTestService(t, mocks.NewMockIKnowledgeRepository(t), map[string]*runtimemocks.MockISourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", config)
	assertAppError(t, err, http.StatusBadRequest, schemaErr.Error())
}

func TestTestConnectionReturnsBadRequestWhenConfigurationInvalid(t *testing.T) {
	t.Parallel()

	driver := runtimemocks.NewMockISourceDriver(t)
	driver.EXPECT().GetConfigurationSchema(mock.Anything).Return(testConfigSchema, nil).Once()

	svc := newTestService(t, mocks.NewMockIKnowledgeRepository(t), map[string]*runtimemocks.MockISourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", util.JSONMap{})
	assertAppError(t, err, http.StatusBadRequest, runtime.ErrInvalidConfiguration.Error())
}

func TestTestConnectionReturnsDriverConnectionError(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	connectionErr := errors.New("connection refused")

	driver := runtimemocks.NewMockISourceDriver(t)
	driver.EXPECT().GetConfigurationSchema(mock.Anything).Return(testConfigSchema, nil).Once()
	driver.EXPECT().TestConnection(mock.Anything, config).Return(connectionErr).Once()

	svc := newTestService(t, mocks.NewMockIKnowledgeRepository(t), map[string]*runtimemocks.MockISourceDriver{"pinecone": driver})

	err := svc.TestConnection(context.Background(), "pinecone", config)
	assertAppError(t, err, http.StatusBadRequest, runtime.ErrConnectionFailed.Error())
}

func TestListDriversReturnsSourceDrivers(t *testing.T) {
	t.Parallel()

	sourceSchema := util.JSONMap{"type": "object"}

	sourceDriver := runtimemocks.NewMockISourceDriver(t)
	sourceDriver.EXPECT().GetConfigurationSchema(mock.Anything).Return(sourceSchema, nil).Once()

	svc := newTestService(
		t,
		mocks.NewMockIKnowledgeRepository(t),
		map[string]*runtimemocks.MockISourceDriver{"fs": sourceDriver},
	)

	result, err := svc.ListDrivers(context.Background())
	require.NoError(t, err)
	require.Len(t, result.SourceDrivers, 1)
	assert.Equal(t, "fs", result.SourceDrivers[0].Name)
	assert.Equal(t, sourceSchema, result.SourceDrivers[0].ConfigurationSchema)
}
