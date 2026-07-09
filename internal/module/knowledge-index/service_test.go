package knowledgeindex_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

var testIndexConfigSchema = util.JSONMap{
	"type": "object",
	"properties": util.JSONMap{
		"dimension": util.JSONMap{"type": "number"},
	},
	"required": []any{"dimension"},
}

type mockIndexDriver struct {
	mock.Mock
}

func (m *mockIndexDriver) Reader(config util.JSONMap) (runtime.IIndexReader, error) {
	args := m.Called(config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(runtime.IIndexReader), args.Error(1)
}

func (m *mockIndexDriver) Writer(config util.JSONMap) (runtime.IIndexWriter, error) {
	args := m.Called(config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(runtime.IIndexWriter), args.Error(1)
}

func (m *mockIndexDriver) TestConnection(ctx context.Context, config util.JSONMap) error {
	return m.Called(ctx, config).Error(0)
}

func (m *mockIndexDriver) GetConfigurationSchema(ctx context.Context) (util.JSONMap, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(util.JSONMap), args.Error(1)
}

type passthroughTxManager struct{}

func (m *passthroughTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockJobQueue struct {
	mock.Mock
}

func (m *mockJobQueue) Push(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (int64, error) {
	argsMock := m.Called(ctx, args, opts)
	return argsMock.Get(0).(int64), argsMock.Error(1)
}

func (m *mockJobQueue) JobGet(ctx context.Context, id int64) (*rivertype.JobRow, error) {
	argsMock := m.Called(ctx, id)
	return argsMock.Get(0).(*rivertype.JobRow), argsMock.Error(1)
}

type noopJobQueue struct{}

func (n *noopJobQueue) Push(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (int64, error) {
	return 0, nil
}

func (n *noopJobQueue) JobGet(ctx context.Context, id int64) (*rivertype.JobRow, error) {
	return nil, nil
}

type testServiceOptions struct {
	txManager   repository.ITxManager
	riverClient queue.IJobQueue
}

func expectSuccessfulIndexConnection(driver *mockIndexDriver, config util.JSONMap) {
	driver.On("GetConfigurationSchema", mock.Anything).Return(testIndexConfigSchema, nil).Once()
	driver.On("TestConnection", mock.Anything, config).Return(nil).Once()
}

func newTestService(
	t *testing.T,
	repo repository.IKnowledgeIndexRepository,
	drivers map[string]*mockIndexDriver,
	opts ...func(*testServiceOptions),
) *knowledgeindex.Service {
	t.Helper()

	options := testServiceOptions{
		txManager:   &passthroughTxManager{},
		riverClient: &noopJobQueue{},
	}
	for _, opt := range opts {
		opt(&options)
	}

	registry := runtime.NewRegistry[runtime.IIndexDriver]()
	for name, driver := range drivers {
		registry.MustRegister(name, driver)
	}
	indexManager := runtime.NewIndexManager(registry)
	return knowledgeindex.NewService(
		repo,
		mocks.NewMockIIndexedKnowledgeItemRepository(t),
		indexManager,
		options.riverClient,
		options.txManager,
	)
}

func TestFilterDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	opts := filter.Default[model.KnowledgeIndex]()
	expected := page.NewPaginated([]model.KnowledgeIndex{{Name: "Index A"}}, 1, 0, 10)
	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		FilterInKnowledge(mock.Anything, knowledgeID, opts).
		Return(expected, nil)

	svc := newTestService(t, repo, nil)

	result, err := svc.Filter(context.Background(), knowledgeID, opts)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	id := uuid.New().String()
	expected := &model.KnowledgeIndex{ID: id, Name: "Found"}
	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		FindByIDInKnowledge(mock.Anything, knowledgeID, id).
		Return(expected, nil)

	svc := newTestService(t, repo, nil)

	result, err := svc.FindByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresIndexAndReturnsIt(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	config := util.JSONMap{"dimension": 1536}
	var stored *model.KnowledgeIndex

	driver := new(mockIndexDriver)
	expectSuccessfulIndexConnection(driver, config)

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID string, index *model.KnowledgeIndex) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			stored = index
			index.ID = uuid.New().String()
		}).
		Return(nil)
	repo.EXPECT().
		FindByIDInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		RunAndReturn(func(_ context.Context, _ string, id string) (*model.KnowledgeIndex, error) {
			return &model.KnowledgeIndex{ID: id}, nil
		})

	riverClient := new(mockJobQueue)
	riverClient.On("Push", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil).Once()

	svc := newTestService(t, repo, map[string]*mockIndexDriver{"pinecone": driver}, func(o *testServiceOptions) {
		o.riverClient = riverClient
	})

	result, err := svc.Create(context.Background(), knowledgeID, knowledgeindex.CreateKnowledgeIndexDTO{
		Name:          "Docs Index",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "Docs Index", result.Name)
	assert.Equal(t, "pinecone", result.Driver)
	assert.Equal(t, config, result.Configuration)
	assert.Equal(t, knowledgeID, result.KnowledgeID)

	require.NotNil(t, stored)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Driver, stored.Driver)
	assert.Equal(t, result.Configuration, stored.Configuration)
	assert.Equal(t, result.KnowledgeID, stored.KnowledgeID)

	driver.AssertExpectations(t)
	riverClient.AssertExpectations(t)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	config := util.JSONMap{"dimension": 1536}
	expectedErr := errors.New("create failed")

	driver := new(mockIndexDriver)
	expectSuccessfulIndexConnection(driver, config)

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Return(expectedErr)

	svc := newTestService(t, repo, map[string]*mockIndexDriver{"pinecone": driver})

	_, err := svc.Create(context.Background(), knowledgeID, knowledgeindex.CreateKnowledgeIndexDTO{
		Name:          "Index",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.ErrorIs(t, err, expectedErr)

	driver.AssertExpectations(t)
}

func TestUpdateDelegatesNameToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	id := uuid.New().String()
	newName := "Updated Name"

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		UpdateNameInKnowledge(mock.Anything, knowledgeID, id, newName).
		Return(nil)

	svc := newTestService(t, repo, nil)

	err := svc.Update(context.Background(), knowledgeID, id, knowledgeindex.UpdateKnowledgeIndexDTO{
		Name: &newName,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	id := uuid.New().String()
	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		DeleteInKnowledge(mock.Anything, knowledgeID, id).
		Return(nil)

	svc := newTestService(t, repo, nil)

	err := svc.DeleteByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
}
