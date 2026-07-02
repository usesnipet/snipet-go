package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IKnowledgeIndexRepository interface {
	FilterInKnowledge(
		ctx context.Context,
		knowledgeID string,
		filter *filter.Options[model.KnowledgeIndex],
	) (*page.Paginated[model.KnowledgeIndex], error)

	FindByIDInKnowledge(
		ctx context.Context,
		knowledgeID string,
		id string,
	) (*model.KnowledgeIndex, error)

	CreateInKnowledge(
		ctx context.Context,
		knowledgeID string,
		index *model.KnowledgeIndex,
	) error

	UpdateNameInKnowledge(
		ctx context.Context,
		knowledgeID string,
		id string,
		name string,
	) error

	DeleteInKnowledge(
		ctx context.Context,
		knowledgeID string,
		id string,
	) error
}

type KnowledgeIndexRepository struct {
	*Repository[model.KnowledgeIndex]
}

func NewKnowledgeIndexRepository(db *gorm.DB) IKnowledgeIndexRepository {
	return &KnowledgeIndexRepository{
		Repository: NewRepository[model.KnowledgeIndex](db),
	}
}

func (r *KnowledgeIndexRepository) FilterInKnowledge(
	ctx context.Context,
	knowledgeID string,
	indexFilter *filter.Options[model.KnowledgeIndex],
) (*page.Paginated[model.KnowledgeIndex], error) {
	return r.Filter(
		ctx,
		filter.Merge(
			indexFilter,
			filter.New[model.KnowledgeIndex](filter.WhereEq("knowledge_id", knowledgeID)),
		),
	)
}

func (r *KnowledgeIndexRepository) FindByIDInKnowledge(
	ctx context.Context,
	knowledgeID string,
	id string,
) (*model.KnowledgeIndex, error) {
	paginated, err := r.FilterInKnowledge(
		ctx,
		knowledgeID,
		filter.New[model.KnowledgeIndex](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("knowledge index not found")
	}
	return paginated.First(), nil
}

func (r *KnowledgeIndexRepository) CreateInKnowledge(
	ctx context.Context,
	knowledgeID string,
	index *model.KnowledgeIndex,
) error {
	index.KnowledgeID = knowledgeID
	if err := r.Create(ctx, index); err != nil {
		return err
	}
	return nil
}

func (r *KnowledgeIndexRepository) UpdateNameInKnowledge(
	ctx context.Context,
	knowledgeID string,
	id string,
	name string,
) error {
	affected, err := gorm.G[model.KnowledgeIndex](r.db(ctx)).
		Where("knowledge_id = ? AND id = ?", knowledgeID, id).
		Updates(ctx, model.KnowledgeIndex{Name: name})
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("knowledge index not found")
	}
	return nil
}

func (r *KnowledgeIndexRepository) DeleteInKnowledge(
	ctx context.Context,
	knowledgeID string,
	id string,
) error {
	affected, err := gorm.G[model.KnowledgeIndex](r.db(ctx)).
		Where("knowledge_id = ? AND id = ?", knowledgeID, id).
		Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("knowledge index not found")
	}
	return nil
}
