package repository

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IConversationMessageRepository interface {
	FilterInConversation(
		ctx context.Context,
		clientCode string,
		conversationID string,
		filter *filter.Options[model.ConversationMessage],
	) (*page.Paginated[model.ConversationMessage], error)
}

type ConversationMessageRepository struct {
	*Repository[model.ConversationMessage]
	clientRepo IClientRepository
}

func (r *ConversationMessageRepository) FilterInConversation(
	ctx context.Context,
	clientCode string,
	conversationID string,
	filterOptions *filter.Options[model.ConversationMessage],
) (*page.Paginated[model.ConversationMessage], error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	if filterOptions == nil {
		filterOptions = filter.Default[model.ConversationMessage]()
	}
	total, err := gorm.G[model.ConversationMessage](r.DB).
		Joins(clause.LeftJoin.Association("conversations"), nil).
		Where("conversation_id = ?", conversationID).
		Where("conversations.client_id = ?", client.ID).Count(ctx, "1 = 1")

	if err != nil {
		return nil, err
	}

	chain, err := filterOptions.ToGorm(gorm.G[model.ConversationMessage](r.DB))
	if err != nil {
		return nil, err
	}

	data, err := chain.Joins(clause.LeftJoin.Association("conversations"), nil).
		Where("conversation_id = ?", conversationID).
		Where("conversations.client_id = ?", client.ID).Find(ctx)

	return page.NewPaginated(data, total, int64(filterOptions.Skip), int64(filterOptions.Take)), err
}

func NewConversationMessageRepository(db *gorm.DB, clientRepo IClientRepository) IConversationMessageRepository {
	return &ConversationMessageRepository{
		Repository: NewRepository[model.ConversationMessage](db),
		clientRepo: clientRepo,
	}
}
