package subscriber

import (
	"context"

	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
)

type Persistence struct {
	logger               *logger.Logger
	executionMessageRepo repository.IExecutionMessageRepository
	executionRepo        repository.IExecutionRepository

	executionID string
}

func NewPersistence(
	executionRepo repository.IExecutionRepository,
	executionMessageRepo repository.IExecutionMessageRepository,
	logger *logger.Logger,
	executionID string,
) *Persistence {
	return &Persistence{
		logger:               logger,
		executionRepo:        executionRepo,
		executionMessageRepo: executionMessageRepo,
		executionID:          executionID,
	}
}

func (p *Persistence) Handle(ctx context.Context, event runtime.IEvent) error {
	switch event := event.(type) {
	case runtime.ExecutionMessageAddedEvent:
		return p.handleExecutionMessageAddedEvent(ctx, event)
	case runtime.ExecutionTurnCompletedEvent:
		return p.handleExecutionTurnCompletedEvent(ctx, event)
	case runtime.ExecutionStatusChangedEvent:
		return p.handleExecutionStatusChangedEvent(ctx, event)
	case runtime.ExecutionFinishedEvent:
		return p.handleExecutionFinishedEvent(ctx)
	}

	return nil
}

func (p *Persistence) handleExecutionTurnCompletedEvent(
	ctx context.Context,
	event runtime.ExecutionTurnCompletedEvent,
) error {
	return p.executionRepo.UpdateByID(ctx, p.executionID, &model.Execution{Turns: event.Turn})
}

func (p *Persistence) handleExecutionMessageAddedEvent(
	ctx context.Context,
	event runtime.ExecutionMessageAddedEvent,
) error {
	message := model.ExecutionMessage{Message: event.Message, ExecutionID: p.executionID}
	return p.executionMessageRepo.CreateInExecution(ctx, p.executionID, message)
}

func (p *Persistence) handleExecutionStatusChangedEvent(
	ctx context.Context,
	event runtime.ExecutionStatusChangedEvent,
) error {
	return p.executionRepo.UpdateByID(ctx, p.executionID, &model.Execution{
		Status:       event.Status,
		ErrorMessage: event.ErrorMessage,
	})
}

func (p *Persistence) handleExecutionFinishedEvent(ctx context.Context) error {
	return p.executionRepo.UpdateByID(ctx, p.executionID, &model.Execution{
		Status: runtime.ExecutionStatusCompleted,
	})
}
