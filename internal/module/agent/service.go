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
	"github.com/usesnipet/snipet/internal/runtime/driver"
	"github.com/usesnipet/snipet/internal/runtime/message"
	"github.com/usesnipet/snipet/internal/util"
)

type Service struct {
	agentRepo            repository.IAgentRepository
	engine               *runtime.Engine
	llmManager           *driver.Manager[driver.ILLM]
	toolManager          *driver.Manager[driver.ITool]
	executionRepo        repository.IExecutionRepository
	executionMessageRepo repository.IExecutionMessageRepository
	logger               *logger.Logger
}

func NewService(
	agentRepo repository.IAgentRepository,
	engine *runtime.Engine,
	llmManager *driver.Manager[driver.ILLM],
	toolManager *driver.Manager[driver.ITool],
	executionRepo repository.IExecutionRepository,
	executionMessageRepo repository.IExecutionMessageRepository,
	logger *logger.Logger,
) *Service {
	return &Service{
		agentRepo:            agentRepo,
		engine:               engine,
		llmManager:           llmManager,
		toolManager:          toolManager,
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
	agent := &model.Agent{
		Name:         dto.Name,
		Description:  dto.Description,
		Instructions: dto.Instructions,
		Configuration: model.AgentConfiguration{
			LLMs: util.Map(dto.LLMs, func(llm LLMConfigDTO) runtime.LLMConfig {
				return runtime.LLMConfig{
					Key:    llm.Key,
					Config: llm.Config,
				}
			}),
			Tools: runtime.ToolConfig(dto.Tools),
		},
	}

	if err := s.validateLLMs(ctx, agent.Configuration.LLMs); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}
	if err := s.validateTools(ctx, agent.Configuration.Tools); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	if err := s.agentRepo.Create(ctx, agent); err != nil {
		return nil, err
	}
	return agent, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateAgentDTO) error {
	updates := &model.Agent{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Description != nil {
		updates.Description = *dto.Description
	}
	if dto.Instructions != nil {
		updates.Instructions = *dto.Instructions
	}
	if len(dto.LLMs) > 0 {
		updates.Configuration.LLMs = util.Map(dto.LLMs, func(llm LLMConfigDTO) runtime.LLMConfig {
			return runtime.LLMConfig{
				Key:    llm.Key,
				Config: llm.Config,
			}
		})

		if err := s.validateLLMs(ctx, updates.Configuration.LLMs); err != nil {
			return apperr.BadRequest(err.Error())
		}
	}
	if dto.Tools != nil {
		updates.Configuration.Tools = runtime.ToolConfig(dto.Tools)
		if err := s.validateTools(ctx, updates.Configuration.Tools); err != nil {
			return apperr.BadRequest(err.Error())
		}
	}
	return s.agentRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.agentRepo.DeleteByID(ctx, id)
}

type EventHandler func(event runtime.IEvent) error

func (s *Service) Run(ctx context.Context, id string, dto RunAgentDTO, onEvent EventHandler) error {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	execution := &model.Execution{
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
				runtime.WithInitialMessages(
					message.NewMessage(message.MessageRoleUser, dto.Message),
				),
			},
			OnEvent: func(event runtime.IEvent) error {
				if err := s.handleExecutionEvent(ctx, execution, event); err != nil {
					return err
				}
				return onEvent(event)
			},
		},
	)
}

func (s *Service) handleExecutionEvent(ctx context.Context, execution *model.Execution, event runtime.IEvent) error {
	switch event := event.(type) {
	case runtime.ExecutionStatusChangedEvent:
		execution.Status = event.Status
		execution.ErrorMessage = event.ErrorMessage
		execution.Turns = event.Turns
		return s.executionRepo.UpdateByID(ctx, execution.ID, execution)
	case runtime.ExecutionMessageAddedEvent:
		return s.executionMessageRepo.CreateInExecution(
			ctx,
			execution.ID,
			util.Map(event.Messages, func(msg message.Message) model.ExecutionMessage {
				return *(&model.ExecutionMessage{}).FromRuntimeExecutionMessage(msg)
			}),
		)
	}
	return nil
}

func (s *Service) validateLLMs(ctx context.Context, llms []runtime.LLMConfig) error {
	return s.llmManager.ValidateMultipleConfigurationsByKey(util.Map(llms, func(cfg runtime.LLMConfig) driver.Configuration {
		return driver.Configuration(cfg)
	})...)
}

func (s *Service) validateTools(ctx context.Context, tools runtime.ToolConfig) error {
	configs := make([]driver.Configuration, 0, len(tools))
	for key, config := range tools {
		configs = append(configs, driver.Configuration{Key: key, Config: config})
	}
	return s.toolManager.ValidateMultipleConfigurationsByKey(configs...)
}
