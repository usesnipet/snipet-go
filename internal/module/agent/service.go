package agent

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/authz"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/agent/subscriber"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/runtime/execution"
	"github.com/usesnipet/snipet/pkg/jsonx"
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

func (s *Service) Filter(ctx context.Context, tenantID string) (*page.Paginated[model.Agent], error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	return s.agentRepo.Filter(ctx, filter.New[model.Agent](filter.WhereEq("tenant_id", tenantID)))
}

// findInTenant fetches by id then verifies the row belongs to tenantID.
func (s *Service) findInTenant(ctx context.Context, tenantID, id string) (*model.Agent, error) {
	found, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if found.TenantID != tenantID {
		return nil, apperr.NotFound("agent not found")
	}
	return found, nil
}

func (s *Service) FindByID(ctx context.Context, tenantID, id string) (*model.Agent, error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	return s.findInTenant(ctx, tenantID, id)
}

func (s *Service) Create(ctx context.Context, tenantID string, dto CreateAgentDTO) (*model.Agent, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return nil, err
	}
	if err := s.validateLLMIDs(ctx, tenantID, dto.LLMIDs); err != nil {
		return nil, err
	}

	agent := &model.Agent{
		TenantID:     tenantID,
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

func (s *Service) Update(ctx context.Context, tenantID, id string, dto UpdateAgentDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}
	if _, err := s.findInTenant(ctx, tenantID, id); err != nil {
		return err
	}
	if dto.LLMIDs != nil {
		if err := s.validateLLMIDs(ctx, tenantID, dto.LLMIDs); err != nil {
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
		}

		if dto.LLMIDs != nil {
			return s.agentRepo.ReplaceLLMs(ctx, id, dto.LLMIDs)
		}
		return nil
	})
}

func (s *Service) DeleteByID(ctx context.Context, tenantID, id string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}
	if _, err := s.findInTenant(ctx, tenantID, id); err != nil {
		return err
	}
	return s.agentRepo.DeleteByID(ctx, id)
}

// Run stays reachable both directly (POST /agent/{id}/run, API-key
// authenticated) and internally via session.Service.Run (client-widget
// end-user JWT, or an API key on the session's anyAuthMiddleware). Unlike
// the CRUD methods above it takes no tenantID param — there's no
// tenant-staff caller here, no /tenants/{tenant_id}/... URL. When the
// caller authenticated with an API key, its TenantID must match the
// target agent's (404 on mismatch, not 403, to avoid confirming the id
// exists in another tenant). When called through a client-widget identity
// instead, no additional check is needed here — the session/client
// ownership chain that got the caller this far already scoped it.
func (s *Service) Run(ctx context.Context, input RunInput, subscribers ...execution.Subscriber) error {
	agent, err := s.agentRepo.FindByID(ctx, input.AgentID)
	if err != nil {
		return err
	}

	if identity, err := auth.CurrentApiKey(ctx); err == nil {
		if agent.TenantID != identity.TenantID {
			return apperr.NotFound("agent not found")
		}
	}

	initialMessages := make([]msg.Message, 0)

	if input.SessionID != nil {
		history, err := s.executionMessageRepo.FilterInSession(
			ctx,
			*input.SessionID,
			filter.New[model.ExecutionMessage](
				filter.OrderDesc("created_at"),
				filter.Take(20),
			),
		)
		if err != nil {
			return err
		}
		for _, em := range history.Data {
			initialMessages = append(initialMessages, em.Message)
		}
	}

	userMessage := msg.NewMessage(msg.RoleUser, input.Message)
	initialMessages = append(initialMessages, userMessage)

	executionModel := &model.Execution{
		TenantID:     agent.TenantID,
		SessionID:    input.SessionID,
		AgentID:      agent.ID,
		Status:       execution.StatusRunning,
		ErrorMessage: "",
		Turns:        0,
		Metadata:     jsonx.JSONMap{},
	}
	executionRuntime, err := executionModel.ToRuntimeExecution(
		execution.WithAgent(agent.ToRuntimeAgent()),
		execution.WithInitialMessages(initialMessages...),
	)
	if err != nil {
		return err
	}
	if err := s.executionRepo.Create(ctx, executionModel); err != nil {
		return err
	}
	err = s.executionMessageRepo.CreateInExecution(ctx, executionModel.ID, model.ExecutionMessage{Message: userMessage, TenantID: agent.TenantID})
	if err != nil {
		return err
	}

	executionRuntime.Subscribe(
		subscriber.NewPersistence(s.executionRepo, s.executionMessageRepo, s.logger, executionModel.ID, agent.TenantID),
	)
	executionRuntime.Subscribe(subscribers...)

	return s.engine.Start(ctx, executionRuntime)
}

func (s *Service) validateLLMIDs(ctx context.Context, tenantID string, llmIDs []string) error {
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
		filter.WhereEq("tenant_id", tenantID),
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
