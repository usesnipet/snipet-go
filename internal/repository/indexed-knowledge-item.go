package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IIndexedKnowledgeItemRepository interface {
	FilterInIndex(
		ctx context.Context,
		knowledgeID string,
		indexID string,
		filter *filter.Options[model.IndexedKnowledgeItem],
	) (*page.Paginated[model.IndexedKnowledgeItem], error)

	FindByIDInIndex(
		ctx context.Context,
		knowledgeID string,
		indexID string,
		id string,
	) (*model.IndexedKnowledgeItem, error)

	CreateInIndex(
		ctx context.Context,
		knowledgeID string,
		indexID string,
		item *model.IndexedKnowledgeItem,
	) error

	UpdateInIndex(
		ctx context.Context,
		knowledgeID string,
		indexID string,
		id string,
		item *model.IndexedKnowledgeItem,
	) error

	DeleteInIndex(
		ctx context.Context,
		knowledgeID string,
		indexID string,
		id string,
	) error
}

type IndexedKnowledgeItemRepository struct {
	*Repository[model.IndexedKnowledgeItem]
}

func NewIndexedKnowledgeItemRepository(db *gorm.DB) IIndexedKnowledgeItemRepository {
	return &IndexedKnowledgeItemRepository{
		Repository: NewRepository[model.IndexedKnowledgeItem](db),
	}
}

func (r *IndexedKnowledgeItemRepository) ensureIndexInKnowledge(
	ctx context.Context,
	knowledgeID string,
	indexID string,
) error {
	count, err := gorm.G[model.KnowledgeIndex](r.db(ctx)).
		Where("knowledge_id = ? AND id = ?", knowledgeID, indexID).
		Count(ctx, "1 = 1")
	if err != nil {
		return err
	}
	if count == 0 {
		return apperr.NotFound("knowledge index not found")
	}
	return nil
}

func (r *IndexedKnowledgeItemRepository) FilterInIndex(
	ctx context.Context,
	knowledgeID string,
	indexID string,
	itemFilter *filter.Options[model.IndexedKnowledgeItem],
) (*page.Paginated[model.IndexedKnowledgeItem], error) {
	if err := r.ensureIndexInKnowledge(ctx, knowledgeID, indexID); err != nil {
		return nil, err
	}
	return r.Filter(
		ctx,
		filter.Merge(
			itemFilter,
			filter.New[model.IndexedKnowledgeItem](filter.WhereEq("index_id", indexID)),
		),
	)
}

func (r *IndexedKnowledgeItemRepository) FindByIDInIndex(
	ctx context.Context,
	knowledgeID string,
	indexID string,
	id string,
) (*model.IndexedKnowledgeItem, error) {
	paginated, err := r.FilterInIndex(
		ctx,
		knowledgeID,
		indexID,
		filter.New[model.IndexedKnowledgeItem](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("indexed knowledge item not found")
	}
	return paginated.First(), nil
}

func (r *IndexedKnowledgeItemRepository) CreateInIndex(
	ctx context.Context,
	knowledgeID string,
	indexID string,
	item *model.IndexedKnowledgeItem,
) error {
	if err := r.ensureIndexInKnowledge(ctx, knowledgeID, indexID); err != nil {
		return err
	}
	item.IndexID = indexID
	if err := r.Create(ctx, item); err != nil {
		return err
	}
	return nil
}

func (r *IndexedKnowledgeItemRepository) UpdateInIndex(
	ctx context.Context,
	knowledgeID string,
	indexID string,
	id string,
	item *model.IndexedKnowledgeItem,
) error {
	if err := r.ensureIndexInKnowledge(ctx, knowledgeID, indexID); err != nil {
		return err
	}
	affected, err := gorm.G[model.IndexedKnowledgeItem](r.db(ctx)).
		Where("index_id = ? AND id = ?", indexID, id).
		Updates(ctx, *item)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("indexed knowledge item not found")
	}
	return nil
}

func (r *IndexedKnowledgeItemRepository) DeleteInIndex(
	ctx context.Context,
	knowledgeID string,
	indexID string,
	id string,
) error {
	if err := r.ensureIndexInKnowledge(ctx, knowledgeID, indexID); err != nil {
		return err
	}
	affected, err := gorm.G[model.IndexedKnowledgeItem](r.db(ctx)).
		Where("index_id = ? AND id = ?", indexID, id).
		Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("indexed knowledge item not found")
	}
	return nil
}
