package bot

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

func (s *Service) FindBy(ctx context.Context) (*database.Paginated[model.Bot], error) {
	return s.repository.FindBy(ctx, filter.Default[model.Bot]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Bot, error) {
	paginated, err := s.repository.FindBy(
		ctx,
		filter.New[model.Bot](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("bot not found")
	}
	return paginated.First(), nil
}

func (s *Service) Create(ctx context.Context, dto CreateBotDTO) (*model.Bot, error) {
	bot := &model.Bot{
		Name:          dto.Name,
		Description:   dto.Description,
		Configuration: dto.Configuration,
	}
	if err := s.repository.Create(ctx, bot); err != nil {
		return nil, err
	}
	return bot, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateBotDTO) error {
	updates := &model.Bot{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Description != nil {
		updates.Description = *dto.Description
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
