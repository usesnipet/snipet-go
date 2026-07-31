package repository

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IExecutionMessageRepository interface {
	CreateInExecution(ctx context.Context, executionID string, messages []model.ExecutionMessage) error
	CountByExecutionID(ctx context.Context, executionID string) (int64, error)
	FilterInSession(
		ctx context.Context,
		sessionID string,
		filter *filter.Options[model.ExecutionMessage],
	) (*page.Paginated[model.ExecutionMessage], error)
	ListBySessionID(ctx context.Context, sessionID string) ([]model.ExecutionMessage, error)
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

func (r *ExecutionMessageRepository) FilterInSession(
	ctx context.Context,
	sessionID string,
	filterOptions *filter.Options[model.ExecutionMessage],
) (*page.Paginated[model.ExecutionMessage], error) {
	if filterOptions == nil {
		filterOptions = filter.Default[model.ExecutionMessage]()
	}
	if len(filterOptions.Order.Fields) == 0 {
		filterOptions = filter.Merge(
			filterOptions,
			filter.New[model.ExecutionMessage](filter.OrderAsc("timestamp")),
		)
	}

	sessionExecutions := r.db(ctx).Table("executions").Select("id").Where("session_id = ?", sessionID)
	base := func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&model.ExecutionMessage{}).
			Where("execution_id IN (?)", sessionExecutions)
	}

	var total int64
	if err := base(r.db(ctx)).Count(&total).Error; err != nil {
		return nil, err
	}

	chain, err := filterOptions.ToGormTx(base(r.db(ctx)))
	if err != nil {
		return nil, err
	}

	var data []model.ExecutionMessage
	if err := chain.Find(&data).Error; err != nil {
		return nil, err
	}

	return page.NewPaginated(data, total, int64(filterOptions.Skip), int64(filterOptions.Take)), nil
}

func (r *ExecutionMessageRepository) ListBySessionID(ctx context.Context, sessionID string) ([]model.ExecutionMessage, error) {
	var data []model.ExecutionMessage
	err := r.db(ctx).Table("execution_messages").
		Select("execution_messages.*").
		Joins("INNER JOIN executions ON executions.id = execution_messages.execution_id").
		Where("executions.session_id = ?", sessionID).
		Order("execution_messages.timestamp ASC").
		Order("execution_messages.sequence ASC").
		Find(&data).Error
	return data, err
}
