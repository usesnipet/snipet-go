package knowledge

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	repository repository.IKnowledgeRepository
}

func (s *Service) Filter(ctx context.Context, filter *filter.Options[model.Knowledge]) (*page.Paginated[model.Knowledge], error) {
	return s.repository.Filter(ctx, filter)
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Knowledge, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateKnowledgeDTO) (*model.Knowledge, error) {
	memory := &model.Knowledge{
		Name:          dto.Name,
		Description:   dto.Description,
		Driver:        dto.Driver,
		Configuration: dto.Configuration,
	}
	if err := s.repository.Create(ctx, memory); err != nil {
		return nil, err
	}
	return memory, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateKnowledgeDTO) error {
	updates := &model.Knowledge{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Description != nil {
		updates.Description = *dto.Description
	}
	return s.repository.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repository.DeleteByID(ctx, id)
}

func NewService(repository repository.IKnowledgeRepository) *Service {
	return &Service{repository: repository}
}
