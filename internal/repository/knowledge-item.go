package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	UpsertMany(
		ctx context.Context,
		items []model.KnowledgeItem,
		batchSize int,
	) error

	HashesByExternalIDInKnowledge(
		ctx context.Context,
		knowledgeID string,
	) (map[string]string, error)

	DeleteByExternalIDsInKnowledge(
		ctx context.Context,
		knowledgeID string,
		externalIDs []string,
	) (int64, error)

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

func (r *KnowledgeItemRepository) UpsertMany(ctx context.Context, items []model.KnowledgeItem, batchSize int) error {
	if len(items) == 0 {
		return nil
	}
	return r.db(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "knowledge_id"}, {Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "hash", "metadata", "last_modified"}),
		}).
		CreateInBatches(items, batchSize).
		Error
}

func (r *KnowledgeItemRepository) HashesByExternalIDInKnowledge(
	ctx context.Context,
	knowledgeID string,
) (map[string]string, error) {
	items, err := gorm.G[model.KnowledgeItem](r.db(ctx)).
		Select("external_id, hash").
		Where("knowledge_id = ?", knowledgeID).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	hashes := make(map[string]string, len(items))
	for _, item := range items {
		hashes[item.ExternalID] = item.Hash
	}
	return hashes, nil
}

func (r *KnowledgeItemRepository) DeleteByExternalIDsInKnowledge(
	ctx context.Context,
	knowledgeID string,
	externalIDs []string,
) (int64, error) {
	if len(externalIDs) == 0 {
		return 0, nil
	}
	affected, err := gorm.G[model.KnowledgeItem](r.db(ctx)).
		Where("knowledge_id = ? AND external_id IN ?", knowledgeID, externalIDs).
		Delete(ctx)
	if err != nil {
		return 0, err
	}
	return int64(affected), nil
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
