package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IClientRepository interface {
	IFilterableRepository[model.Client]
	ICreatableRepository[model.Client]
	FilterByUserID(
		ctx context.Context,
		userID string,
		filter *filter.Options[model.Client],
	) (*page.Paginated[model.Client], error)
	FindByCode(ctx context.Context, code string) (*model.Client, error)
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

func (r *ClientRepository) FilterByUserID(
	ctx context.Context,
	userID string,
	filterOptions *filter.Options[model.Client],
) (*page.Paginated[model.Client], error) {
	if filterOptions == nil {
		filterOptions = filter.Default[model.Client]()
	}

	var total int64
	err := r.DB.WithContext(ctx).Table("clients").
		Joins("LEFT JOIN client_to_users ctu ON ctu.client_id = clients.id").
		Where("ctu.user_id = ?", userID).
		Count(&total).Error

	if err != nil {
		return nil, err
	}

	var data []model.Client
	tx, err := filterOptions.ToGormTx(r.DB.WithContext(ctx).Table("clients"))
	if err != nil {
		return nil, err
	}
	err = tx.Joins(
		"LEFT JOIN client_to_users ctu ON ctu.client_id = clients.id",
	).Where("ctu.user_id = ?", userID).Find(&data).Error

	return page.NewPaginated(data, total, int64(filterOptions.Skip), int64(filterOptions.Take)), err
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
