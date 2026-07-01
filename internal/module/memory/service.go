package memory

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	repository repository.IMemoryRepository
}

func (s *Service) Filter(ctx context.Context, filter *filter.Options[model.Memory]) (*page.Paginated[model.Memory], error) {
	return s.repository.Filter(ctx, filter)
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
	return s.repository.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *Service) SetAsDefault(ctx context.Context, dto SetAsDefaultMemoryDTO) error {
	memory, err := s.repository.FindByID(ctx, dto.MemoryID)
	if err != nil {
		return err
	}
	if memory.IsDefault {
		return apperr.BadRequest("memory is already default")
	}

	return s.repository.SetAsDefault(ctx, memory.Type, dto.MemoryID)
}

func NewService(repository repository.IMemoryRepository) *Service {
	return &Service{repository: repository}
}
