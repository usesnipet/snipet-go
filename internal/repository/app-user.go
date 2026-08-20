package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IAppUserRepository interface {
	FilterInApp(
		ctx context.Context,
		appCode string,
		filter *filter.Options[model.AppUser],
	) (*page.Paginated[model.AppUser], error)
	FindByIDInApp(
		ctx context.Context,
		appCode string,
		id string,
	) (*model.AppUser, error)
	FindByExternalIDInApp(
		ctx context.Context,
		appCode string,
		externalID string,
	) (*model.AppUser, error)

	CreateInApp(
		ctx context.Context,
		appCode string,
		user *model.AppUser,
		externalID *string,
	) error
}

type AppUserRepository struct {
	*Repository[model.AppUser]
	appRepo IAppRepository
}

func NewAppUserRepository(db *gorm.DB, appRepo IAppRepository) IAppUserRepository {
	return &AppUserRepository{
		Repository: NewRepository[model.AppUser](db),
		appRepo:    appRepo,
	}
}

func (r *AppUserRepository) FindByExternalIDInApp(
	ctx context.Context,
	appCode string,
	externalID string,
) (*model.AppUser, error) {
	app, err := r.appRepo.FindByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}

	userIDs := gorm.G[model.AppToAppUser](r.db(ctx)).
		Where("app_id = ?", app.ID).
		Where("external_id = ?", externalID).
		Select("app_user_id")

	user, err := gorm.G[model.AppUser](r.db(ctx)).
		Where("id IN (?)", userIDs).
		First(ctx)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AppUserRepository) FilterInApp(
	ctx context.Context,
	appCode string,
	userFilter *filter.Options[model.AppUser],
) (*page.Paginated[model.AppUser], error) {
	if userFilter == nil {
		userFilter = filter.Default[model.AppUser]()
	}
	app, err := r.appRepo.FindByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}

	userIDs := gorm.G[model.AppToAppUser](r.db(ctx)).
		Where("app_id = ?", app.ID).
		Select("app_user_id")

	total, err := gorm.G[model.AppUser](r.db(ctx)).
		Where("id IN (?)", userIDs).
		Count(ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	chain, err := userFilter.ToGorm(gorm.G[model.AppUser](r.db(ctx)))
	if err != nil {
		return nil, err
	}

	data, err := chain.
		Where("id IN (?)", userIDs).
		Find(ctx)

	return page.NewPaginated(data, total, int64(userFilter.Skip), int64(userFilter.Take)), err
}

func (r *AppUserRepository) FindByIDInApp(
	ctx context.Context,
	appCode string,
	id string,
) (*model.AppUser, error) {
	paginated, err := r.FilterInApp(
		ctx,
		appCode,
		filter.New[model.AppUser](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("app user not found")
	}
	return paginated.First(), nil
}

func (r *AppUserRepository) CreateInApp(
	ctx context.Context,
	appCode string,
	user *model.AppUser,
	externalID *string,
) error {
	app, err := r.appRepo.FindByCode(ctx, appCode)
	if err != nil {
		return err
	}
	return WithTransaction(ctx, r.DB, func(ctx context.Context) error {
		if err := gorm.G[model.AppUser](r.db(ctx)).Create(ctx, user); err != nil {
			return err
		}
		return gorm.G[model.AppToAppUser](r.db(ctx)).Create(ctx, &model.AppToAppUser{
			AppID:      app.ID,
			AppUserID:  user.ID,
			ExternalID: externalID,
		})
	})
}
