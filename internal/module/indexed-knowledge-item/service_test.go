package indexedknowledgeitem_test

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
	indexedknowledgeitem "github.com/usesnipet/snipet/internal/module/indexed-knowledge-item"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/repository/mocks"
	"github.com/usesnipet/snipet/internal/util"
)

func newTestService(repo repository.IIndexedKnowledgeItemRepository) *indexedknowledgeitem.Service {
	return indexedknowledgeitem.NewService(repo)
}

func TestFilterDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	indexID := uuid.New().String()
	opts := filter.Default[model.IndexedKnowledgeItem]()
	expected := page.NewPaginated([]model.IndexedKnowledgeItem{{Status: model.IndexedStatusPending}}, 1, 0, 10)
	repo := mocks.NewMockIIndexedKnowledgeItemRepository(t)
	repo.EXPECT().
		FilterInIndex(mock.Anything, knowledgeID, indexID, opts).
		Return(expected, nil)

	svc := newTestService(repo)

	result, err := svc.Filter(context.Background(), knowledgeID, indexID, opts)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	indexID := uuid.New().String()
	id := uuid.New().String()
	expected := &model.IndexedKnowledgeItem{ID: id, Status: model.IndexedStatusIndexed}
	repo := mocks.NewMockIIndexedKnowledgeItemRepository(t)
	repo.EXPECT().
		FindByIDInIndex(mock.Anything, knowledgeID, indexID, id).
		Return(expected, nil)

	svc := newTestService(repo)

	result, err := svc.FindByID(context.Background(), knowledgeID, indexID, id)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateStoresItemAndReturnsIt(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	indexID := uuid.New().String()
	knowledgeItemID := uuid.New().String()
	metadata := util.JSONMap{"chunk": 1}
	var stored *model.IndexedKnowledgeItem

	repo := mocks.NewMockIIndexedKnowledgeItemRepository(t)
	repo.EXPECT().
		CreateInIndex(mock.Anything, knowledgeID, indexID, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID, gotIndexID string, item *model.IndexedKnowledgeItem) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			assert.Equal(t, indexID, gotIndexID)
			stored = item
			item.ID = uuid.New().String()
		}).
		Return(nil)

	svc := newTestService(repo)

	result, err := svc.Create(context.Background(), knowledgeID, indexID, indexedknowledgeitem.CreateIndexedKnowledgeItemDTO{
		Status:          model.IndexedStatusPending,
		Version:         1,
		Hash:            "abc123",
		Metadata:        metadata,
		KnowledgeItemID: knowledgeItemID,
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, model.IndexedStatusPending, result.Status)
	assert.Equal(t, 1, result.Version)
	assert.Equal(t, "abc123", result.Hash)
	assert.Equal(t, metadata, result.Metadata)
	assert.Equal(t, indexID, result.IndexID)
	assert.Equal(t, knowledgeItemID, result.KnowledgeItemID)

	require.NotNil(t, stored)
	assert.Equal(t, result.Status, stored.Status)
	assert.Equal(t, result.Version, stored.Version)
	assert.Equal(t, result.Hash, stored.Hash)
	assert.Equal(t, result.Metadata, stored.Metadata)
	assert.Equal(t, result.IndexID, stored.IndexID)
	assert.Equal(t, result.KnowledgeItemID, stored.KnowledgeItemID)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	indexID := uuid.New().String()
	expectedErr := errors.New("create failed")
	repo := mocks.NewMockIIndexedKnowledgeItemRepository(t)
	repo.EXPECT().
		CreateInIndex(mock.Anything, knowledgeID, indexID, mock.Anything).
		Return(expectedErr)

	svc := newTestService(repo)

	_, err := svc.Create(context.Background(), knowledgeID, indexID, indexedknowledgeitem.CreateIndexedKnowledgeItemDTO{
		KnowledgeItemID: uuid.New().String(),
	})
	require.ErrorIs(t, err, expectedErr)
}

func TestUpdateDelegatesPartialFieldsToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	indexID := uuid.New().String()
	id := uuid.New().String()
	newStatus := model.IndexedStatusIndexed
	newVersion := 2

	repo := mocks.NewMockIIndexedKnowledgeItemRepository(t)
	repo.EXPECT().
		UpdateInIndex(mock.Anything, knowledgeID, indexID, id, mock.Anything).
		Run(func(_ context.Context, gotKnowledgeID, gotIndexID, gotID string, updates *model.IndexedKnowledgeItem) {
			assert.Equal(t, knowledgeID, gotKnowledgeID)
			assert.Equal(t, indexID, gotIndexID)
			assert.Equal(t, id, gotID)
			assert.Equal(t, newStatus, updates.Status)
			assert.Equal(t, newVersion, updates.Version)
			assert.Empty(t, updates.Hash)
			assert.Nil(t, updates.Metadata)
		}).
		Return(nil)

	svc := newTestService(repo)

	err := svc.Update(context.Background(), knowledgeID, indexID, id, indexedknowledgeitem.UpdateIndexedKnowledgeItemDTO{
		Status:  &newStatus,
		Version: &newVersion,
	})
	require.NoError(t, err)
}

func TestDeleteByIDDelegatesToRepository(t *testing.T) {
	t.Parallel()

	knowledgeID := uuid.New().String()
	indexID := uuid.New().String()
	id := uuid.New().String()
	repo := mocks.NewMockIIndexedKnowledgeItemRepository(t)
	repo.EXPECT().
		DeleteInIndex(mock.Anything, knowledgeID, indexID, id).
		Return(nil)

	svc := newTestService(repo)

	err := svc.DeleteByID(context.Background(), knowledgeID, indexID, id)
	require.NoError(t, err)
}
