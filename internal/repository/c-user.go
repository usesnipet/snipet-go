package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICUserRepository interface {
	FilterInClient(
		ctx context.Context,
		clientCode string,
		filter *filter.Options[model.CUser],
	) (*page.Paginated[model.CUser], error)
	FindByIDInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) (*model.CUser, error)
	CreateInClient(
		ctx context.Context,
		clientCode string,
		cUser *model.CUser,
		externalID *string,
	) error
	FindByExternalIDInClient(
		ctx context.Context,
		clientCode string,
		externalID string,
	) (*model.CUser, error)
}

type CUserRepository struct {
	*Repository[model.CUser]
	clientRepo IClientRepository
}

func NewCUserRepository(db *gorm.DB, clientRepo IClientRepository) ICUserRepository {
	return &CUserRepository{
		Repository: NewRepository[model.CUser](db),
		clientRepo: clientRepo,
	}
}

func (r *CUserRepository) FindByExternalIDInClient(
	ctx context.Context,
	clientCode string,
	externalID string,
) (*model.CUser, error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	user, err := gorm.G[model.CUser](r.DB).
		Joins(clause.LeftJoin.Association("client_to_users"), nil).
		Where("client_to_users.client_id = ?", client.ID).
		Where("client_to_users.external_id = ?", externalID).
		First(ctx)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *CUserRepository) FilterInClient(
	ctx context.Context,
	clientCode string,
	cUserFilter *filter.Options[model.CUser],
) (*page.Paginated[model.CUser], error) {
	if cUserFilter == nil {
		cUserFilter = filter.Default[model.CUser]()
	}
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	total, err := gorm.G[model.CUser](r.DB).
		Joins(clause.LeftJoin.Association("client_to_users"), nil).
		Where("client_to_users.client_id = ?", client.ID).Count(ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	chain, err := cUserFilter.ToGorm(gorm.G[model.CUser](r.DB))
	if err != nil {
		return nil, err
	}

	data, err := chain.Joins(clause.LeftJoin.Association("client_to_users"), nil).
		Where("client_to_users.client_id = ?", client.ID).Find(ctx)

	return page.NewPaginated(data, total, int64(cUserFilter.Skip), int64(cUserFilter.Take)), err
}

func (r *CUserRepository) FindByIDInClient(
	ctx context.Context,
	clientCode string,
	id string,
) (*model.CUser, error) {
	paginated, err := r.FilterInClient(
		ctx,
		clientCode,
		filter.New[model.CUser](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("c user not found")
	}
	return paginated.First(), nil
}

func (r *CUserRepository) CreateInClient(
	ctx context.Context,
	clientCode string,
	cUser *model.CUser,
	externalID *string,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	cUser.ClientToUsers = []model.ClientToUser{
		{
			ClientID:     client.ID,
			ClientUserID: cUser.ID,
			ExternalID:   externalID,
		},
	}
	return r.Create(ctx, cUser)
}
