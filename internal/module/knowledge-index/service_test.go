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
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/util"
)

func newTestService(repo repository.IKnowledgeIndexRepository) *knowledgeindex.Service {
	return knowledgeindex.NewService(repo)
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

	svc := newTestService(repo)

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

	svc := newTestService(repo)

	result, err := svc.FindByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresIndexAndReturnsIt(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	config := util.JSONMap{"dimension": 1536}
	var stored *model.KnowledgeIndex

	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID string, index *model.KnowledgeIndex) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			stored = index
			index.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(repo)

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
	expectedErr := errors.New("create failed")
	repo := mocks.NewMockIKnowledgeIndexRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Return(expectedErr)

	svc := newTestService(repo)

	_, err := svc.Create(context.Background(), knowledgeID, knowledgeindex.CreateKnowledgeIndexDTO{
		Name:          "Index",
		Driver:        "pinecone",
		Configuration: util.JSONMap{},
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

	svc := newTestService(repo)

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

	svc := newTestService(repo)

	err := svc.DeleteByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
}
