package bot

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	botRepo repository.IBotRepository
}

func NewService(botRepo repository.IBotRepository) *Service {
	return &Service{
		botRepo: botRepo,
	}
}

func (s *Service) Filter(ctx context.Context) (*page.Paginated[model.Bot], error) {
	return s.botRepo.Filter(ctx, filter.Default[model.Bot]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Bot, error) {
	return s.botRepo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateBotDTO) (*model.Bot, error) {
	bot := &model.Bot{
		Name:          dto.Name,
		Description:   dto.Description,
		Configuration: dto.Configuration,
	}
	if err := s.botRepo.Create(ctx, bot); err != nil {
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
	return s.botRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.botRepo.DeleteByID(ctx, id)
}
