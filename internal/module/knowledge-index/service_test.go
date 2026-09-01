package knowledgeindex_test

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
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	knowledgemocks "github.com/usesnipet/snipet/pkg/driver/knowledge/mocks"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

var testIndexConfigSchema = jsonx.JSONMap{
	"type": "object",
	"properties": jsonx.JSONMap{
		"dimension": jsonx.JSONMap{"type": "number"},
	},
	"required": []any{"dimension"},
}

type noopPool struct{}

func (noopPool) Submit(context.Context, queue.Job) error { return nil }

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

type testServiceOptions struct {
	txManager repository.ITxManager
	pool      queue.IPool
}

func newMockIndex(t *testing.T, name string, schema jsonx.JSONMap) *knowledgemocks.MockIIndexDriver {
	t.Helper()

	indexDriver := knowledgemocks.NewMockIIndexDriver(t)
	indexDriver.EXPECT().Info().Return(driver.Info{
		Key:                 name,
		Name:                name,
		Description:         name,
		ConfigurationSchema: schema,
	})
	indexDriver.EXPECT().Validate().Return(nil)
	return indexDriver
}

// knowledgeRepoFor stubs IKnowledgeRepository.FindByID(knowledgeID) to
// return a Knowledge — Create verifies the parent exists before inserting
// the index.
func knowledgeRepoFor(t *testing.T, knowledgeID string) repository.IKnowledgeRepository {
	t.Helper()

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		FindByID(mock.Anything, knowledgeID).
		Return(&model.Knowledge{ID: knowledgeID}, nil)
	return repo
}

func newTestService(
	t *testing.T,
	repo repository.IKnowledgeIndexRepository,
	knowledgeRepo repository.IKnowledgeRepository,
	drivers map[string]knowledge.IIndexDriver,
	opts ...func(*testServiceOptions),
) *knowledgeindex.Service {
	t.Helper()

	options := testServiceOptions{
		txManager: newPassthroughTxManager(t),
		pool:      noopPool{},
	}
	for _, opt := range opts {
		opt(&options)
	}

	reg := driver.NewRegistry[knowledge.IIndexDriver](logger.NewLogger(logger.LevelError))
	for _, d := range drivers {
		reg.MustRegister(d, nil)
	}
	indexManager := manager.NewDriverManager(reg)
	syncWorker := knowledgeindex.NewSyncIndexWorker(
		indexManager,
		manager.NewDriverManager(driver.NewRegistry[knowledge.ISourceDriver](logger.NewLogger(logger.LevelError))),
		mocks.NewMockIKnowledgeRepository(t),
		mocks.NewMockIKnowledgeItemRepository(t),
		repo,
		mocks.NewMockIIndexedKnowledgeItemRepository(t),
		logger.NewLogger(logger.LevelError),
	)
	return knowledgeindex.NewService(
		repo,
		knowledgeRepo,
		mocks.NewMockIIndexedKnowledgeItemRepository(t),
		indexManager,
		options.pool,
		syncWorker,
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

	svc := newTestService(t, repo, mocks.NewMockIKnowledgeRepository(t), nil)

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

	svc := newTestService(t, repo, mocks.NewMockIKnowledgeRepository(t), nil)

	result, err := svc.FindByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresIndexAndReturnsIt(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	config := jsonx.JSONMap{"dimension": 1536}
	var stored *model.KnowledgeIndex

	indexDriver := newMockIndex(t, "pinecone", testIndexConfigSchema)
	indexDriver.EXPECT().TestConnection(mock.Anything, config).Return(nil)

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID string, index *model.KnowledgeIndex) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			stored = index
			index.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(t, repo, knowledgeRepoFor(t, knowledgeID), map[string]knowledge.IIndexDriver{"pinecone": indexDriver})

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
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	config := jsonx.JSONMap{"dimension": 1536}
	expectedErr := errors.New("create failed")

	indexDriver := newMockIndex(t, "pinecone", testIndexConfigSchema)
	indexDriver.EXPECT().TestConnection(mock.Anything, config).Return(nil)

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Return(expectedErr)

	svc := newTestService(t, repo, knowledgeRepoFor(t, knowledgeID), map[string]knowledge.IIndexDriver{"pinecone": indexDriver})

	_, err := svc.Create(context.Background(), knowledgeID, knowledgeindex.CreateKnowledgeIndexDTO{
		Name:          "Index",
		Driver:        "pinecone",
		Configuration: config,
	})
	require.ErrorIs(t, err, expectedErr)
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

	svc := newTestService(t, repo, mocks.NewMockIKnowledgeRepository(t), nil)

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

	svc := newTestService(t, repo, mocks.NewMockIKnowledgeRepository(t), nil)

	err := svc.DeleteByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
}

func TestListDriversReturnsIndexDrivers(t *testing.T) {
	t.Parallel()

	indexSchema := jsonx.JSONMap{"type": "object", "title": "index"}
	indexDriver := newMockIndex(t, "rag", indexSchema)

	svc := newTestService(t, mocks.NewMockIKnowledgeIndexRepository(t), mocks.NewMockIKnowledgeRepository(t), map[string]knowledge.IIndexDriver{"rag": indexDriver})

	result, err := svc.ListDrivers(context.Background())
	require.NoError(t, err)
	require.Len(t, result.IndexDrivers, 1)
	assert.Equal(t, "rag", result.IndexDrivers[0].Name)
	assert.Equal(t, indexSchema, result.IndexDrivers[0].ConfigurationSchema)
}
