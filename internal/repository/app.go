package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IAppRepository interface {
	IRepository[model.App]

	FindByCode(ctx context.Context, code string) (*model.App, error)
	UpdateByCode(ctx context.Context, code string, updates *model.App) error
	DeleteByCode(ctx context.Context, code string) error
}

type AppRepository struct {
	*Repository[model.App]
}

func NewAppRepository(db *gorm.DB) IAppRepository {
	return &AppRepository{
		Repository: NewRepository[model.App](db),
	}
}

func (r *AppRepository) FindByCode(ctx context.Context, code string) (*model.App, error) {
	paginated, err := r.Filter(ctx, filter.New[model.App](filter.WhereEq("code", code)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("app not found")
	}
	return paginated.First(), nil
}

func (r *AppRepository) UpdateByCode(ctx context.Context, code string, updates *model.App) error {
	affected, err := gorm.G[model.App](r.db(ctx)).Where("code = ?", code).Updates(ctx, *updates)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("entity not found")
	}
	return nil
}

func (r *AppRepository) DeleteByCode(ctx context.Context, code string) error {
	affected, err := gorm.G[model.App](r.db(ctx)).Where("code = ?", code).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("app not found")
	}
	return nil
}
