package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IClientRepository interface {
	IFilterableRepository[model.Client]
	FindByCode(ctx context.Context, code string) (*model.Client, error)

	ICreatableRepository[model.Client]

	UpdateByCode(ctx context.Context, code string, updates *model.Client) error

	DeleteByCode(ctx context.Context, code string) error
}

type ClientRepository struct {
	*Repository[model.Client]
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &ClientRepository{
		Repository: NewRepository[model.Client](db),
	}
}

func (r *ClientRepository) FindByCode(ctx context.Context, code string) (*model.Client, error) {
	paginated, err := r.Filter(ctx, filter.New[model.Client](filter.WhereEq("code", code)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("client not found")
	}
	return paginated.First(), nil
}

func (r *ClientRepository) UpdateByCode(ctx context.Context, code string, updates *model.Client) error {
	affected, err := gorm.G[model.Client](r.DB).Where("code = ?", code).Updates(ctx, *updates)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("entity not found")
	}
	return nil
}

func (r *ClientRepository) DeleteByCode(ctx context.Context, code string) error {
	affected, err := gorm.G[model.Client](r.DB).Where("code = ?", code).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("client not found")
	}
	return nil
}
