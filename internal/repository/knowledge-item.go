package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IKnowledgeItemRepository interface {
	FilterInKnowledge(
		ctx context.Context,
		knowledgeID string,
		filter *filter.Options[model.KnowledgeItem],
	) (*page.Paginated[model.KnowledgeItem], error)

	FindByIDInKnowledge(
		ctx context.Context,
		knowledgeID string,
		id string,
	) (*model.KnowledgeItem, error)

	CreateInKnowledge(
		ctx context.Context,
		knowledgeID string,
		item *model.KnowledgeItem,
	) error

	UpdateInKnowledge(
		ctx context.Context,
		knowledgeID string,
		id string,
		item *model.KnowledgeItem,
	) error

	DeleteInKnowledge(
		ctx context.Context,
		knowledgeID string,
		id string,
	) error
}

type KnowledgeItemRepository struct {
	*Repository[model.KnowledgeItem]
}

func NewKnowledgeItemRepository(db *gorm.DB) IKnowledgeItemRepository {
	return &KnowledgeItemRepository{
		Repository: NewRepository[model.KnowledgeItem](db),
	}
}

func (r *KnowledgeItemRepository) FilterInKnowledge(
	ctx context.Context,
	knowledgeID string,
	itemFilter *filter.Options[model.KnowledgeItem],
) (*page.Paginated[model.KnowledgeItem], error) {
	return r.Filter(
		ctx,
		filter.Merge(
			itemFilter,
			filter.New[model.KnowledgeItem](filter.WhereEq("knowledge_id", knowledgeID)),
		),
	)
}

func (r *KnowledgeItemRepository) FindByIDInKnowledge(
	ctx context.Context,
	knowledgeID string,
	id string,
) (*model.KnowledgeItem, error) {
	paginated, err := r.FilterInKnowledge(
		ctx,
		knowledgeID,
		filter.New[model.KnowledgeItem](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("knowledge item not found")
	}
	return paginated.First(), nil
}

func (r *KnowledgeItemRepository) CreateInKnowledge(
	ctx context.Context,
	knowledgeID string,
	item *model.KnowledgeItem,
) error {
	item.KnowledgeID = knowledgeID
	if err := r.Create(ctx, item); err != nil {
		return err
	}
	return nil
}

func (r *KnowledgeItemRepository) UpdateInKnowledge(
	ctx context.Context,
	knowledgeID string,
	id string,
	item *model.KnowledgeItem,
) error {
	affected, err := gorm.G[model.KnowledgeItem](r.db(ctx)).
		Where("knowledge_id = ? AND id = ?", knowledgeID, id).
		Updates(ctx, *item)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("knowledge item not found")
	}
	return nil
}

func (r *KnowledgeItemRepository) DeleteInKnowledge(
	ctx context.Context,
	knowledgeID string,
	id string,
) error {
	affected, err := gorm.G[model.KnowledgeItem](r.db(ctx)).
		Where("knowledge_id = ? AND id = ?", knowledgeID, id).
		Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("knowledge item not found")
	}
	return nil
}
