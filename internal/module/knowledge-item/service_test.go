package knowledgeitem_test

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
	knowledgeitem "github.com/usesnipet/snipet/internal/module/knowledge-item"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/util"
)

func newTestService(repo repository.IKnowledgeItemRepository) *knowledgeitem.Service {
	return knowledgeitem.NewService(repo)
}

func TestFilterDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	opts := filter.Default[model.KnowledgeItem]()
	expected := page.NewPaginated([]model.KnowledgeItem{{Name: "Item A"}}, 1, 0, 10)
	repo := mocks.NewMockIKnowledgeItemRepository(t)
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
	expected := &model.KnowledgeItem{ID: id, Name: "Found"}
	repo := mocks.NewMockIKnowledgeItemRepository(t)
	repo.EXPECT().
		FindByIDInKnowledge(mock.Anything, knowledgeID, id).
		Return(expected, nil)

	svc := newTestService(repo)

	result, err := svc.FindByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresItemAndReturnsIt(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	metadata := util.JSONMap{"source": "drive"}
	var stored *model.KnowledgeItem

	repo := mocks.NewMockIKnowledgeItemRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID string, item *model.KnowledgeItem) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			stored = item
			item.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(repo)

	result, err := svc.Create(context.Background(), knowledgeID, knowledgeitem.CreateKnowledgeItemDTO{
		ExternalID: "ext-1",
		Name:       "Document",
		Hash:       "abc123",
		Metadata:   metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "ext-1", result.ExternalID)
	assert.Equal(t, "Document", result.Name)
	assert.Equal(t, "abc123", result.Hash)
	assert.Equal(t, metadata, result.Metadata)
	assert.Equal(t, knowledgeID, result.KnowledgeID)

	require.NotNil(t, stored)
	assert.Equal(t, result.ExternalID, stored.ExternalID)
	assert.Equal(t, result.Name, stored.Name)
	assert.Equal(t, result.Hash, stored.Hash)
	assert.Equal(t, result.Metadata, stored.Metadata)
	assert.Equal(t, result.KnowledgeID, stored.KnowledgeID)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	expectedErr := errors.New("create failed")
	repo := mocks.NewMockIKnowledgeItemRepository(t)
	repo.EXPECT().
		CreateInKnowledge(mock.Anything, knowledgeID, mock.Anything).
		Return(expectedErr)

	svc := newTestService(repo)

	_, err := svc.Create(context.Background(), knowledgeID, knowledgeitem.CreateKnowledgeItemDTO{})
	require.ErrorIs(t, err, expectedErr)
}

func TestUpdateDelegatesPartialFieldsToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	id := uuid.New().String()
	newName := "Updated Name"
	newHash := "def456"

	repo := mocks.NewMockIKnowledgeItemRepository(t)
	repo.EXPECT().
		UpdateInKnowledge(mock.Anything, knowledgeID, id, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID, gotID string, updates *model.KnowledgeItem) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			assert.Equal(t, id, gotID)
			assert.Equal(t, newName, updates.Name)
			assert.Equal(t, newHash, updates.Hash)
			assert.Empty(t, updates.ExternalID)
			assert.Nil(t, updates.Metadata)
		}).
		Return(nil)

	svc := newTestService(repo)

	err := svc.Update(context.Background(), knowledgeID, id, knowledgeitem.UpdateKnowledgeItemDTO{
		Name: &newName,
		Hash: &newHash,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	id := uuid.New().String()
	repo := mocks.NewMockIKnowledgeItemRepository(t)
	repo.EXPECT().
		DeleteInKnowledge(mock.Anything, knowledgeID, id).
		Return(nil)

	svc := newTestService(repo)

	err := svc.DeleteByID(context.Background(), knowledgeID, id)
	require.NoError(t, err)
}
