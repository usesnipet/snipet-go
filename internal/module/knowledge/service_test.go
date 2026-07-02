package knowledge_test

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
	knowledge "github.com/usesnipet/snipet/internal/module/knowledge"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/util"
)

func newTestService(repo repository.IKnowledgeRepository) *knowledge.Service {
	return knowledge.NewService(repo)
}

func TestFilterDelegatesToRepository(t *testing.T) {
	t.Parallel()

	opts := filter.Default[model.Knowledge]()
	expected := page.NewPaginated([]model.Knowledge{{Name: "Knowledge A"}}, 1, 0, 10)
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Filter(mock.Anything, opts).
		Return(expected, nil)

	svc := newTestService(repo)

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

	svc := newTestService(repo)

	result, err := svc.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresKnowledgeAndReturnsIt(t *testing.T) {
	t.Parallel()

	config := util.JSONMap{"index": "docs"}
	var stored *model.Knowledge

	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Run(func(_ context.Context, k *model.Knowledge) {
			stored = k
			k.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(repo)

	result, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
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

	expectedErr := errors.New("create failed")
	repo := mocks.NewMockIKnowledgeRepository(t)
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expectedErr)

	svc := newTestService(repo)

	_, err := svc.Create(context.Background(), knowledge.CreateKnowledgeDTO{
		Name:          "Knowledge",
		Driver:        "pinecone",
		Configuration: util.JSONMap{},
	})
	require.ErrorIs(t, err, expectedErr)
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

	svc := newTestService(repo)

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

	svc := newTestService(repo)

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

	svc := newTestService(repo)

	err := svc.DeleteByID(context.Background(), id)
	require.NoError(t, err)
}
