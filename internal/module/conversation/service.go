package conversation

import (
	"context"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	conversationRepo        repository.IConversationRepository
	conversationMessageRepo repository.IConversationMessageRepository
	memoryRepo              repository.IMemoryRepository
}

func (s *Service) FindBy(ctx context.Context, clientCode string) (*page.Paginated[model.Conversation], error) {
	return s.conversationRepo.FilterInClient(ctx, clientCode, filter.Default[model.Conversation]())
}

func (s *Service) FindByID(ctx context.Context, clientCode string, id string) (*model.Conversation, error) {
	return s.conversationRepo.FindByIDInClient(ctx, clientCode, id)
}

func (s *Service) Create(ctx context.Context, clientCode string, dto CreateConversationDTO) (*model.Conversation, error) {
	memoryUUID, err := uuid.Parse(dto.MemoryID)
	if err != nil {
		return nil, apperr.BadRequest("invalid memory id")
	}
	botUUID, err := uuid.Parse(dto.BotID)
	if err != nil {
		return nil, apperr.BadRequest("invalid bot id")
	}

	memory, err := s.memoryRepo.FindByID(ctx, dto.MemoryID)
	if err != nil {
		return nil, err
	}

	if memory.Type != model.MemoryTypeConversation {
		return nil, apperr.BadRequest("invalid memory type")
	}

	conversation := &model.Conversation{
		MemoryID: memoryUUID,
		BotID:    botUUID,
		Metadata: dto.Metadata,
	}
	if err := s.conversationRepo.CreateInClient(ctx, clientCode, conversation); err != nil {
		return nil, err
	}
	return conversation, nil
}

func (s *Service) DeleteByID(ctx context.Context, clientCode string, id string) error {
	return s.conversationRepo.DeleteInClient(ctx, clientCode, id)
}

func (s *Service) FindMessages(
	ctx context.Context,
	clientCode string,
	conversationID string,
	messageFilter FindMessagesFilterDTO,
) (*page.Paginated[model.ConversationMessage], error) {
	return s.conversationMessageRepo.FilterInConversation(
		ctx,
		clientCode,
		conversationID,
		filter.New[model.ConversationMessage](
			filter.Take(messageFilter.Take),
			filter.Skip(messageFilter.Skip),
			filter.OrderBy("created_at", messageFilter.Sort),
		),
	)
}

func NewService(
	conversationRepo repository.IConversationRepository,
	conversationMessageRepo repository.IConversationMessageRepository,
	memoryRepo repository.IMemoryRepository,
) *Service {
	return &Service{
		conversationRepo:        conversationRepo,
		conversationMessageRepo: conversationMessageRepo,
		memoryRepo:              memoryRepo,
	}
}
