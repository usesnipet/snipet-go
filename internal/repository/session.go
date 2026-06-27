package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type ISessionRepository interface {
	FilterInClient(
		ctx context.Context,
		clientCode string,
		filter *filter.Options[model.Session],
	) (*page.Paginated[model.Session], error)
	FindByIDInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) (*model.Session, error)
	CreateInClient(
		ctx context.Context,
		clientCode string,
		session *model.Session,
	) error
	DeleteInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) error
}

type SessionRepository struct {
	*Repository[model.Session]
	clientRepo IClientRepository
}

func (r *SessionRepository) FilterInClient(
	ctx context.Context,
	clientCode string,
	sessionFilter *filter.Options[model.Session],
) (*page.Paginated[model.Session], error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	return r.FilterBy(
		ctx,
		filter.Merge(
			sessionFilter,
			filter.New[model.Session](filter.WhereEq("client_id", client.ID)),
		),
	)
}

func (r *SessionRepository) FindByIDInClient(
	ctx context.Context,
	clientCode string,
	id string,
) (*model.Session, error) {
	paginated, err := r.FilterInClient(
		ctx,
		clientCode,
		filter.New[model.Session](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("session not found")
	}
	return paginated.First(), nil
}

func (r *SessionRepository) CreateInClient(
	ctx context.Context,
	clientCode string,
	session *model.Session,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	session.ClientID = client.ID
	return r.Create(ctx, session)
}

func (r *SessionRepository) DeleteInClient(
	ctx context.Context,
	clientCode string,
	id string,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	affected, err := gorm.G[model.Session](r.DB).Where("client_id = ? AND id = ?", client.ID, id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("session not found")
	}
	return nil
}

func NewSessionRepository(db *gorm.DB, clientRepo IClientRepository) ISessionRepository {
	return &SessionRepository{
		Repository: NewRepository[model.Session](db),
		clientRepo: clientRepo,
	}
}
