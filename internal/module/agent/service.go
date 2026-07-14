package agent

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	agentRepo repository.IAgentRepository
}

func NewService(agentRepo repository.IAgentRepository) *Service {
	return &Service{
		agentRepo: agentRepo,
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
		Name:          dto.Name,
		Description:   dto.Description,
		Configuration: dto.Configuration,
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
	if dto.Configuration != nil {
		updates.Configuration = *dto.Configuration
	}
	return s.agentRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.agentRepo.DeleteByID(ctx, id)
}
