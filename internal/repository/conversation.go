package repository

import (
	"context"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IConversationRepository interface {
	FilterInClient(
		ctx context.Context,
		clientID string,
		filter *filter.Options[model.Conversation],
	) (*page.Paginated[model.Conversation], error)
	FindByIDInClient(
		ctx context.Context,
		clientID string,
		id string,
	) (*model.Conversation, error)
	CreateInClient(
		ctx context.Context,
		clientID string,
		conversation *model.Conversation,
	) error
	DeleteInClient(
		ctx context.Context,
		clientID string,
		id string,
	) error
}

type ConversationRepository struct {
	*Repository[model.Conversation]
}

func (r *ConversationRepository) FilterInClient(
	ctx context.Context,
	clientID string,
	conversationFilter *filter.Options[model.Conversation],
) (*page.Paginated[model.Conversation], error) {
	return r.FilterBy(
		ctx,
		filter.Merge(
			conversationFilter,
			filter.New[model.Conversation](filter.WhereEq("client_id", clientID)),
		),
	)
}

func (r *ConversationRepository) FindByIDInClient(
	ctx context.Context,
	clientID string,
	id string,
) (*model.Conversation, error) {
	paginated, err := r.FilterInClient(ctx, clientID, filter.New[model.Conversation](filter.WhereEq("id", id)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("conversation not found")
	}
	return paginated.First(), nil
}

func (r *ConversationRepository) CreateInClient(
	ctx context.Context,
	clientID string,
	conversation *model.Conversation,
) error {
	var err error
	conversation.ClientID, err = r.parseClientID(clientID)
	if err != nil {
		return err
	}
	return r.Create(ctx, conversation)
}

func (r *ConversationRepository) DeleteInClient(
	ctx context.Context,
	clientID string,
	id string,
) error {
	affected, err := gorm.G[model.Conversation](r.DB).Where("client_id = ? AND id = ?", clientID, id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("conversation not found")
	}
	return nil
}

func (r ConversationRepository) parseClientID(clientID string) (uuid.UUID, error) {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return uuid.Nil, apperr.BadRequest("invalid client id")
	}
	return id, nil
}

func NewConversationRepository(db *gorm.DB) IConversationRepository {
	return &ConversationRepository{
		Repository: NewRepository[model.Conversation](db),
	}
}
