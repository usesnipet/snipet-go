package conversation

import (
	"context"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
)

type Service struct {
	repository IRepository
}

func (s *Service) FindBy(ctx context.Context) (*database.Paginated[model.Conversation], error) {
	return s.repository.FindBy(ctx, filter.Default[model.Conversation]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Conversation, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateConversationDTO) (*model.Conversation, error) {
	memoryID, err := uuid.Parse(dto.MemoryID)
	if err != nil {
		return nil, apperr.BadRequest("invalid memory id")
	}
	botID, err := uuid.Parse(dto.BotID)
	if err != nil {
		return nil, apperr.BadRequest("invalid bot id")
	}

	conversation := &model.Conversation{
		MemoryID: memoryID,
		BotID:    botID,
	}
	if err := s.repository.Create(ctx, conversation); err != nil {
		return nil, err
	}
	return conversation, nil
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repository.DeleteByID(ctx, id)
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
}
