package repository

import (
	"context"

	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IExecutionMessageRepository interface {
	CreateInExecution(ctx context.Context, executionID string, messages []model.ExecutionMessage) error
	CountByExecutionID(ctx context.Context, executionID string) (int64, error)
}

type ExecutionMessageRepository struct {
	*Repository[model.ExecutionMessage]
}

func NewExecutionMessageRepository(db *gorm.DB) IExecutionMessageRepository {
	return &ExecutionMessageRepository{
		Repository: NewRepository[model.ExecutionMessage](db),
	}
}

func (r *ExecutionMessageRepository) CountByExecutionID(ctx context.Context, executionID string) (int64, error) {
	return gorm.G[model.ExecutionMessage](r.db(ctx)).
		Where("execution_id = ?", executionID).
		Count(ctx, "1 = 1")
}

func (r *ExecutionMessageRepository) CreateInExecution(ctx context.Context, executionID string, messages []model.ExecutionMessage) error {
	for i := range messages {
		messages[i].ExecutionID = executionID
	}
	return gorm.G[model.ExecutionMessage](r.db(ctx)).CreateInBatches(ctx, &messages, 100)
}
