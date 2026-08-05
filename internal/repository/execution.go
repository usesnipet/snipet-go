package repository

import (
	"context"

	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IExecutionRepository interface {
	IRepository[model.Execution]
	Upsert(ctx context.Context, execution *model.Execution) error
}

type ExecutionRepository struct {
	*Repository[model.Execution]
}

func (r *ExecutionRepository) Upsert(ctx context.Context, execution *model.Execution) error {
	return r.db(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"status", "turns", "error_message", "metadata", "updated_at",
			}),
		}).
		Create(execution).
		Error
}

func NewExecutionRepository(db *gorm.DB) IExecutionRepository {
	return &ExecutionRepository{
		Repository: NewRepository[model.Execution](db),
	}
}
