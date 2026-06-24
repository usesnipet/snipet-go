package client

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	clientRepo repository.IClientRepository
}

func (s *Service) FilterBy(ctx context.Context) (*page.Paginated[model.Client], error) {
	return s.clientRepo.FilterBy(ctx, filter.Default[model.Client]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Client, error) {
	return s.clientRepo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateClientDTO) (*model.Client, error) {
	client := &model.Client{
		Name: dto.Name,
	}
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateClientDTO) error {
	updates := &model.Client{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	return s.clientRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.clientRepo.DeleteByID(ctx, id)
}

func NewService(clientRepo repository.IClientRepository) *Service {
	return &Service{clientRepo: clientRepo}
}
