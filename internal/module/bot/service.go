package bot

import (
	"context"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/client"
)

type Service struct {
	repository       IRepository
	clientRepository client.IRepository
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

func (s *Service) LinkClientToBot(ctx context.Context, dto LinkClientToBotDTO) error {
	// region Check if ids are valid
	clientUUID, err := uuid.Parse(dto.ClientID)
	if err != nil {
		return apperr.BadRequest("invalid client id")
	}
	botUUID, err := uuid.Parse(dto.BotID)
	if err != nil {
		return apperr.BadRequest("invalid bot id")
	}
	// endregion

	// region Check if client exists
	if _, err = s.clientRepository.FindByID(ctx, dto.ClientID); err != nil {
		return err
	}
	// endregion

	// region Check if bot exists
	if _, err = s.repository.FindByID(ctx, dto.BotID); err != nil {
		return err
	}
	// endregion

	return s.repository.LinkClientToBot(ctx, clientUUID, botUUID)
}

func NewService(repository IRepository, clientRepository client.IRepository) *Service {
	return &Service{
		repository:       repository,
		clientRepository: clientRepository,
	}
}
