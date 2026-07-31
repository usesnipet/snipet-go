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
	"github.com/usesnipet/snipet/internal/model"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	"github.com/usesnipet/snipet/internal/page"
	queuemocks "github.com/usesnipet/snipet/internal/queue/mocks"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
	knowledgemocks "github.com/usesnipet/snipet/pkg/driver/knowledge/mocks"
)

var testIndexConfigSchema = util.JSONMap{
	"type": "object",
	"properties": util.JSONMap{
		"dimension": util.JSONMap{"type": "number"},
	},
	"required": []any{"dimension"},
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

type testServiceOptions struct {
	txManager   repository.ITxManager
	riverClient *queuemocks.MockIJobQueue
}

func newMockIndex(t *testing.T, name string, schema util.JSONMap) *knowledgemocks.MockIIndexDriver {
	t.Helper()

	indexDriver := knowledgemocks.NewMockIIndexDriver(t)
	indexDriver.EXPECT().Info().Return(driver.Info{
		Name:                name,
		Description:         name,
		ConfigurationSchema: schema,
	})
	return indexDriver
}

func newTestService(
	t *testing.T,
	repo repository.IKnowledgeIndexRepository,
	drivers map[string]knowledge.IIndexDriver,
	opts ...func(*testServiceOptions),
) *knowledgeindex.Service {
	t.Helper()

	options := testServiceOptions{
		txManager:   newPassthroughTxManager(t),
		riverClient: newNoopJobQueue(t),
	}
	for _, opt := range opts {
		opt(&options)
	}

	reg := registry.New[knowledge.IIndexDriver]()
	for name, d := range drivers {
		reg.MustRegister(name, d)
	}
	indexManager := runtime.NewDriverManager(reg)
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

	svc := newTestService(t, repo, map[string]knowledge.IIndexDriver{"pinecone": indexDriver})

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
	config := util.JSONMap{"dimension": 1536}
	expectedErr := errors.New("create failed")

	indexDriver := newMockIndex(t, "pinecone", testIndexConfigSchema)
	indexDriver.EXPECT().TestConnection(mock.Anything, config).Return(nil)

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Return(expectedErr)

	svc := newTestService(t, repo, map[string]knowledge.IIndexDriver{"pinecone": indexDriver})

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

func TestListDriversReturnsIndexDrivers(t *testing.T) {
	t.Parallel()

	indexSchema := util.JSONMap{"type": "object", "title": "index"}
	indexDriver := newMockIndex(t, "rag", indexSchema)

	svc := newTestService(t, mocks.NewMockIKnowledgeIndexRepository(t), map[string]knowledge.IIndexDriver{"rag": indexDriver})

	result, err := svc.ListDrivers(context.Background())
	require.NoError(t, err)
	require.Len(t, result.IndexDrivers, 1)
	assert.Equal(t, "rag", result.IndexDrivers[0].Name)
	assert.Equal(t, indexSchema, result.IndexDrivers[0].ConfigurationSchema)
}
