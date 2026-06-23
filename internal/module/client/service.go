package client

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
)

type Service struct {
	repository IRepository
}

func (s *Service) FindBy(ctx context.Context) (*database.Paginated[model.Client], error) {
	return s.repository.FindBy(ctx, filter.Default[model.Client]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Client, error) {
	paginated, err := s.repository.FindBy(
		ctx,
		filter.New[model.Client](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("client not found")
	}
	return paginated.First(), nil
}

func (s *Service) Create(ctx context.Context, dto CreateClientDTO) (*model.Client, error) {
	client := &model.Client{
		Name: dto.Name,
	}
	if err := s.repository.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateClientDTO) error {
	updates := &model.Client{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	return s.repository.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repository.DeleteByID(ctx, id)
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
}
