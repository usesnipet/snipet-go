package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FilterInClient(
		ctx context.Context,
		clientCode string,
		filter *filter.Options[model.User],
	) (*page.Paginated[model.User], error)
	FindByIDInClient(
		ctx context.Context,
		clientCode string,
		id string,
	) (*model.User, error)
	FindByExternalIDInClient(
		ctx context.Context,
		clientCode string,
		externalID string,
	) (*model.User, error)

	CreateInClient(
		ctx context.Context,
		clientCode string,
		user *model.User,
		externalID *string,
	) error
}

type UserRepository struct {
	*Repository[model.User]
	clientRepo IClientRepository
}

func NewUserRepository(db *gorm.DB, clientRepo IClientRepository) IUserRepository {
	return &UserRepository{
		Repository: NewRepository[model.User](db),
		clientRepo: clientRepo,
	}
}

func (r *UserRepository) FindByExternalIDInClient(
	ctx context.Context,
	clientCode string,
	externalID string,
) (*model.User, error) {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	userIDs := gorm.G[model.ClientToUser](r.DB).
		Where("client_id = ?", client.ID).
		Where("external_id = ?", externalID).
		Select("user_id")

	user, err := gorm.G[model.User](r.DB).
		Where("id IN (?)", userIDs).
		First(ctx)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FilterInClient(
	ctx context.Context,
	clientCode string,
	userFilter *filter.Options[model.User],
) (*page.Paginated[model.User], error) {
	if userFilter == nil {
		userFilter = filter.Default[model.User]()
	}
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	userIDs := gorm.G[model.ClientToUser](r.DB).
		Where("client_id = ?", client.ID).
		Select("user_id")

	total, err := gorm.G[model.User](r.DB).
		Where("id IN (?)", userIDs).
		Count(ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	chain, err := userFilter.ToGorm(gorm.G[model.User](r.DB))
	if err != nil {
		return nil, err
	}

	data, err := chain.
		Where("id IN (?)", userIDs).
		Find(ctx)

	return page.NewPaginated(data, total, int64(userFilter.Skip), int64(userFilter.Take)), err
}

func (r *UserRepository) FindByIDInClient(
	ctx context.Context,
	clientCode string,
	id string,
) (*model.User, error) {
	paginated, err := r.FilterInClient(
		ctx,
		clientCode,
		filter.New[model.User](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("c user not found")
	}
	return paginated.First(), nil
}

func (r *UserRepository) CreateInClient(
	ctx context.Context,
	clientCode string,
	user *model.User,
	externalID *string,
) error {
	client, err := r.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[model.User](tx).Create(ctx, user); err != nil {
			return err
		}
		return gorm.G[model.ClientToUser](tx).Create(ctx, &model.ClientToUser{
			ClientID:   client.ID,
			UserID:     user.ID,
			ExternalID: externalID,
		})
	})
}
