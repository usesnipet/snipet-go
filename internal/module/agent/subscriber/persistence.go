package subscriber

import (
	"context"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/execution"
)

type Persistence struct {
	logger               *logger.Logger
	executionMessageRepo repository.IExecutionMessageRepository
	executionRepo        repository.IExecutionRepository

	executionID string
	tenantID    string
}

func NewPersistence(
	executionRepo repository.IExecutionRepository,
	executionMessageRepo repository.IExecutionMessageRepository,
	logger *logger.Logger,
	executionID string,
	tenantID string,
) *Persistence {
	return &Persistence{
		logger:               logger,
		executionRepo:        executionRepo,
		executionMessageRepo: executionMessageRepo,
		executionID:          executionID,
		tenantID:             tenantID,
	}
}

func (p *Persistence) Handle(ctx context.Context, event execution.IEvent) error {
	switch event := event.(type) {
	case execution.MessageAddedEvent:
		return p.handleExecutionMessageAddedEvent(ctx, event)
	case execution.TurnCompletedEvent:
		return p.handleExecutionTurnCompletedEvent(ctx, event)
	case execution.StatusChangedEvent:
		return p.handleExecutionStatusChangedEvent(ctx, event)
	case execution.FinishedEvent:
		return p.handleExecutionFinishedEvent(ctx)
	}

	return nil
}

func (p *Persistence) handleExecutionTurnCompletedEvent(
	ctx context.Context,
	event execution.TurnCompletedEvent,
) error {
	return p.executionRepo.UpdateByID(ctx, p.executionID, &model.Execution{Turns: event.Turn})
}

func (p *Persistence) handleExecutionMessageAddedEvent(
	ctx context.Context,
	event execution.MessageAddedEvent,
) error {
	message := model.ExecutionMessage{Message: event.Message, ExecutionID: p.executionID, TenantID: p.tenantID}
	return p.executionMessageRepo.CreateInExecution(ctx, p.executionID, message)
}

func (p *Persistence) handleExecutionStatusChangedEvent(
	ctx context.Context,
	event execution.StatusChangedEvent,
) error {
	return p.executionRepo.UpdateByID(ctx, p.executionID, &model.Execution{
		Status:       event.Status,
		ErrorMessage: event.ErrorMessage,
	})
}

func (p *Persistence) handleExecutionFinishedEvent(ctx context.Context) error {
	return p.executionRepo.UpdateByID(ctx, p.executionID, &model.Execution{
		Status: execution.StatusCompleted,
	})
}
