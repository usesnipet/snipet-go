package repository

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ISessionMessageRepository interface {
	ICreatableRepository[model.SessionMessage]

	FilterInSession(
		ctx context.Context,
		clientCode string,
		sessionID string,
		filter *filter.Options[model.SessionMessage],
	) (*page.Paginated[model.SessionMessage], error)
}

type SessionMessageRepository struct {
	*Repository[model.SessionMessage]
	clientRepo IClientRepository
}

func NewSessionMessageRepository(db *gorm.DB, clientRepo IClientRepository) ISessionMessageRepository {
	return &SessionMessageRepository{
		Repository: NewRepository[model.SessionMessage](db),
		clientRepo: clientRepo,
	}
}

func (r *SessionMessageRepository) FilterInSession(
	ctx context.Context,
	clientCode string,
	sessionID string,
	filterOptions *filter.Options[model.SessionMessage],
) (*page.Paginated[model.SessionMessage], error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	if filterOptions == nil {
		filterOptions = filter.Default[model.SessionMessage]()
	}
	total, err := gorm.G[model.SessionMessage](r.db(ctx)).
		Joins(clause.LeftJoin.Association("Session"), nil).
		Where("session_id = ?", sessionID).
		Where(`"Session"."client_id" = ?`, client.ID).Count(ctx, "1 = 1")

	if err != nil {
		return nil, err
	}

	chain, err := filterOptions.ToGorm(gorm.G[model.SessionMessage](r.db(ctx)))
	if err != nil {
		return nil, err
	}

	data, err := chain.Joins(clause.LeftJoin.Association("Session"), nil).
		Where("session_id = ?", sessionID).
		Where(`"Session"."client_id" = ?`, client.ID).Find(ctx)

	return page.NewPaginated(data, total, int64(filterOptions.Skip), int64(filterOptions.Take)), err
}
