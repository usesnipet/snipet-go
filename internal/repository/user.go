package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IClientUserRepository interface {
	FilterInClient(
		ctx context.Context,
		clientCode string,
		filter *filter.Options[model.ClientUser],
	) (*page.Paginated[model.ClientUser], error)
	FindByIDInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) (*model.ClientUser, error)
	FindByExternalIDInClient(
		ctx context.Context,
		clientCode string,
		externalID string,
	) (*model.ClientUser, error)

	CreateInClient(
		ctx context.Context,
		clientCode string,
		user *model.ClientUser,
		externalID *string,
	) error
}

type ClientUserRepository struct {
	*Repository[model.ClientUser]
	clientRepo IClientRepository
}

func NewClientUserRepository(db *gorm.DB, clientRepo IClientRepository) IClientUserRepository {
	return &ClientUserRepository{
		Repository: NewRepository[model.ClientUser](db),
		clientRepo: clientRepo,
	}
}

func (r *ClientUserRepository) FindByExternalIDInClient(
	ctx context.Context,
	clientCode string,
	externalID string,
) (*model.ClientUser, error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	userIDs := gorm.G[model.ClientToClientUser](r.db(ctx)).
		Where("client_id = ?", client.ID).
		Where("external_id = ?", externalID).
		Select("user_id")

	user, err := gorm.G[model.ClientUser](r.db(ctx)).
		Where("id IN (?)", userIDs).
		First(ctx)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *ClientUserRepository) FilterInClient(
	ctx context.Context,
	clientCode string,
	userFilter *filter.Options[model.ClientUser],
) (*page.Paginated[model.ClientUser], error) {
	if userFilter == nil {
		userFilter = filter.Default[model.ClientUser]()
	}
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	userIDs := gorm.G[model.ClientToClientUser](r.db(ctx)).
		Where("client_id = ?", client.ID).
		Select("user_id")

	total, err := gorm.G[model.ClientUser](r.db(ctx)).
		Where("id IN (?)", userIDs).
		Count(ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	chain, err := userFilter.ToGorm(gorm.G[model.ClientUser](r.db(ctx)))
	if err != nil {
		return nil, err
	}

	data, err := chain.
		Where("id IN (?)", userIDs).
		Find(ctx)

	return page.NewPaginated(data, total, int64(userFilter.Skip), int64(userFilter.Take)), err
}

func (r *ClientUserRepository) FindByIDInClient(
	ctx context.Context,
	clientCode string,
	id string,
) (*model.ClientUser, error) {
	paginated, err := r.FilterInClient(
		ctx,
		clientCode,
		filter.New[model.ClientUser](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("c user not found")
	}
	return paginated.First(), nil
}

func (r *ClientUserRepository) CreateInClient(
	ctx context.Context,
	clientCode string,
	user *model.ClientUser,
	externalID *string,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	return WithTransaction(ctx, r.DB, func(ctx context.Context) error {
		if err := gorm.G[model.ClientUser](r.db(ctx)).Create(ctx, user); err != nil {
			return err
		}
		return gorm.G[model.ClientToClientUser](r.db(ctx)).Create(ctx, &model.ClientToClientUser{
			ClientID:     client.ID,
			ClientUserID: user.ID,
			ExternalID:   externalID,
		})
	})
}
