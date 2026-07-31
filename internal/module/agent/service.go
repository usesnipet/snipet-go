package agent

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/internal/util/set"
	"github.com/usesnipet/snipet/pkg/msg"
)

type Service struct {
	agentRepo            repository.IAgentRepository
	llmRepo              repository.ILLMRepository
	txManager            repository.ITxManager
	engine               *runtime.Engine
	executionRepo        repository.IExecutionRepository
	executionMessageRepo repository.IExecutionMessageRepository
	logger               *logger.Logger
}

func NewService(
	agentRepo repository.IAgentRepository,
	llmRepo repository.ILLMRepository,
	txManager repository.ITxManager,
	engine *runtime.Engine,
	executionRepo repository.IExecutionRepository,
	executionMessageRepo repository.IExecutionMessageRepository,
	logger *logger.Logger,
) *Service {
	return &Service{
		agentRepo:            agentRepo,
		llmRepo:              llmRepo,
		txManager:            txManager,
		engine:               engine,
		executionRepo:        executionRepo,
		executionMessageRepo: executionMessageRepo,
		logger:               logger,
	}
}

func (s *Service) Filter(ctx context.Context) (*page.Paginated[model.Agent], error) {
	return s.agentRepo.Filter(ctx, filter.Default[model.Agent]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Agent, error) {
	return s.agentRepo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateAgentDTO) (*model.Agent, error) {
	if err := s.validateLLMIDs(ctx, dto.LLMIDs); err != nil {
		return nil, err
	}

	agent := &model.Agent{
		Name:         dto.Name,
		Description:  dto.Description,
		Instructions: dto.Instructions,
	}

	err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.agentRepo.Create(ctx, agent); err != nil {
			return err
		}
		return s.agentRepo.ReplaceLLMs(ctx, agent.ID, dto.LLMIDs)
	})
	if err != nil {
		return nil, err
	}

	return s.agentRepo.FindByID(ctx, agent.ID)
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateAgentDTO) error {
	if dto.LLMIDs != nil {
		if err := s.validateLLMIDs(ctx, dto.LLMIDs); err != nil {
			return err
		}
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		updates := &model.Agent{}
		hasScalarUpdates := false
		if dto.Name != nil {
			updates.Name = *dto.Name
			hasScalarUpdates = true
		}
		if dto.Description != nil {
			updates.Description = *dto.Description
			hasScalarUpdates = true
		}
		if dto.Instructions != nil {
			updates.Instructions = *dto.Instructions
			hasScalarUpdates = true
		}

		if hasScalarUpdates {
			if err := s.agentRepo.UpdateByID(ctx, id, updates); err != nil {
				return err
			}
		} else if _, err := s.agentRepo.FindByID(ctx, id); err != nil {
			return err
		}

		if dto.LLMIDs != nil {
			return s.agentRepo.ReplaceLLMs(ctx, id, dto.LLMIDs)
		}
		return nil
	})
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.agentRepo.DeleteByID(ctx, id)
}

type EventHandler func(event runtime.IEvent) error

func (s *Service) Run(ctx context.Context, input RunInput, onEvent EventHandler) error {
	agent, err := s.agentRepo.FindByID(ctx, input.AgentID)
	if err != nil {
		return err
	}

	historyIDSet := set.New[string]()
	initialMessages := make([]msg.Message, 0)

	if input.SessionID != nil {
		history, err := s.executionMessageRepo.ListBySessionID(ctx, *input.SessionID)
		if err != nil {
			return err
		}
		for _, em := range history {
			historyIDSet.Add(em.ID)
			initialMessages = append(initialMessages, em.Message)
		}
	}

	initialMessages = append(
		initialMessages,
		msg.NewMessage(msg.RoleUser, input.Message),
	)

	execution := &model.Execution{
		SessionID:    input.SessionID,
		AgentID:      agent.ID,
		Status:       runtime.ExecutionStatusRunning,
		ErrorMessage: "",
		Turns:        0,
		Metadata:     util.JSONMap{},
	}
	if err := s.executionRepo.Create(ctx, execution); err != nil {
		return err
	}

	if onEvent == nil {
		onEvent = func(runtime.IEvent) error { return nil }
	}

	return s.engine.Start(
		ctx,
		runtime.StartOptions{
			Agent: agent.ToRuntimeAgent(),
			ExecutionOptions: []runtime.ExecutionOption{
				runtime.WithInitialMessages(initialMessages...),
			},
			OnEvent: func(event runtime.IEvent) error {
				forward, err := s.handleExecutionEvent(ctx, execution, historyIDSet, event)
				if err != nil {
					return err
				}
				if forward == nil {
					return nil
				}
				return onEvent(forward)
			},
		},
	)
}

// handleExecutionEvent persists execution updates and returns the event to forward to
// callers. History messages re-emitted as initial context are skipped so they are not
// duplicated in the new execution or streamed again to the client.
func (s *Service) handleExecutionEvent(
	ctx context.Context,
	execution *model.Execution,
	historyIDSet set.Set[string],
	event runtime.IEvent,
) (runtime.IEvent, error) {
	switch event := event.(type) {
	case runtime.ExecutionStatusChangedEvent:
		execution.Status = event.Status
		execution.ErrorMessage = event.ErrorMessage
		execution.Turns = event.Turns
		if err := s.executionRepo.UpdateByID(ctx, execution.ID, execution); err != nil {
			return nil, err
		}
		if event.Status == runtime.ExecutionStatusFailed {
			errorMessage := msg.NewMessage(msg.RoleAssistant, "An error occurred while processing the request.")
			if err := s.executionMessageRepo.CreateInExecution(ctx, execution.ID, []model.ExecutionMessage{{Message: errorMessage}}); err != nil {
				return nil, err
			}
		}
		return event, nil
	case runtime.ExecutionMessageAddedEvent:
		newMessages := make([]model.ExecutionMessage, 0, len(event.Messages))
		for _, m := range event.Messages {
			if historyIDSet.Contains(m.ID) {
				continue
			}
			executionMessage := model.ExecutionMessage{Message: m}
			newMessages = append(newMessages, executionMessage)
		}
		if len(newMessages) == 0 {
			return nil, nil
		}
		if err := s.executionMessageRepo.CreateInExecution(ctx, execution.ID, newMessages); err != nil {
			return nil, err
		}

		return event, nil
	}
	return event, nil
}

func (s *Service) validateLLMIDs(ctx context.Context, llmIDs []string) error {
	seen := make(map[string]struct{}, len(llmIDs))
	for _, id := range llmIDs {
		if _, ok := seen[id]; ok {
			return apperr.BadRequest("duplicate llm id")
		}
		seen[id] = struct{}{}
	}

	ids := make([]any, len(llmIDs))
	for i, id := range llmIDs {
		ids[i] = id
	}
	found, err := s.llmRepo.Filter(ctx, filter.New[model.LLM](
		filter.WhereIn("id", ids...),
		filter.Take(len(llmIDs)),
	))
	if err != nil {
		return err
	}
	if len(found.Data) != len(llmIDs) {
		return apperr.BadRequest("one or more llm ids were not found")
	}
	return nil
}
