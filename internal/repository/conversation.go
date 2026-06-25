package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IConversationRepository interface {
	FilterInClient(
		ctx context.Context,
		clientCode string,
		filter *filter.Options[model.Conversation],
	) (*page.Paginated[model.Conversation], error)
	FindByIDInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) (*model.Conversation, error)
	CreateInClient(
		ctx context.Context,
		clientCode string,
		conversation *model.Conversation,
	) error
	DeleteInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) error
}

type ConversationRepository struct {
	*Repository[model.Conversation]
	clientRepo IClientRepository
}

func (r *ConversationRepository) FilterInClient(
	ctx context.Context,
	clientCode string,
	conversationFilter *filter.Options[model.Conversation],
) (*page.Paginated[model.Conversation], error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	return r.FilterBy(
		ctx,
		filter.Merge(
			conversationFilter,
			filter.New[model.Conversation](filter.WhereEq("client_id", client.ID)),
		),
	)
}

func (r *ConversationRepository) FindByIDInClient(
	ctx context.Context,
	clientCode string,
	id string,
) (*model.Conversation, error) {
	paginated, err := r.FilterInClient(
		ctx,
		clientCode,
		filter.New[model.Conversation](filter.WhereEq("id", id)),
	)
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
	clientCode string,
	conversation *model.Conversation,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	conversation.ClientID = client.ID
	return r.Create(ctx, conversation)
}

func (r *ConversationRepository) DeleteInClient(
	ctx context.Context,
	clientCode string,
	id string,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	affected, err := gorm.G[model.Conversation](r.DB).Where("client_id = ? AND id = ?", client.ID, id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("conversation not found")
	}
	return nil
}

func NewConversationRepository(db *gorm.DB, clientRepo IClientRepository) IConversationRepository {
	return &ConversationRepository{
		Repository: NewRepository[model.Conversation](db),
		clientRepo: clientRepo,
	}
}
