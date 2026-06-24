package memory

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
)

type Service struct {
	repository IRepository
}

func (s *Service) FindBy(ctx context.Context) (*database.Paginated[model.Memory], error) {
	return s.repository.FindBy(ctx, filter.Default[model.Memory]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Memory, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateMemoryDTO) (*model.Memory, error) {
	memory := &model.Memory{
		Name:          dto.Name,
		Type:          dto.Type,
		IsDefault:     dto.IsDefault,
		Provider:      dto.Provider,
		Configuration: dto.Configuration,
	}
	if err := s.repository.Create(ctx, memory); err != nil {
		return nil, err
	}
	return memory, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateMemoryDTO) error {
	updates := &model.Memory{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Type != nil {
		updates.Type = *dto.Type
	}
	if dto.IsDefault != nil {
		updates.IsDefault = *dto.IsDefault
	}
	if dto.Provider != nil {
		updates.Provider = *dto.Provider
	}
	if dto.Configuration != nil {
		updates.Configuration = *dto.Configuration
	}
	return s.repository.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repository.DeleteByID(ctx, id)
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
}
