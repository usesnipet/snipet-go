package database

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"gorm.io/gorm"
)

type IRepository[T any] interface {
	Create(ctx context.Context, model *T) error
	FindByID(ctx context.Context, id string) (*T, error)
	FindBy(ctx context.Context, filter *filter.Options[T]) (*Paginated[T], error)
	UpdateByID(ctx context.Context, id string, model *T) error
	DeleteByID(ctx context.Context, id string) error
}

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(ctx context.Context, model *T) error {
	return gorm.G[T](r.DB).Create(ctx, model)
}

func (r *Repository[T]) FindBy(ctx context.Context, filter *filter.Options[T]) (*Paginated[T], error) {
	total, err := gorm.G[T](r.DB).Count(ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	chain, err := filter.ToGorm(gorm.G[T](r.DB))
	if err != nil {
		return nil, err
	}
	orgs, err := chain.Find(ctx)
	return NewPaginated(orgs, total, int64(filter.Skip), int64(filter.Take)), err
}

func (r *Repository[T]) UpdateByID(ctx context.Context, id string, model *T) error {
	affected, err := gorm.G[T](r.DB).Where("id = ?", id).Updates(ctx, *model)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("entity not found")
	}
	return nil
}

func (r *Repository[T]) DeleteByID(ctx context.Context, id string) error {
	affected, err := gorm.G[T](r.DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("entity not found")
	}
	return nil
}

func (r *Repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	paginated, err := r.FindBy(ctx, filter.New[T](filter.WhereEq("id", id)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("entity not found")
	}
	return paginated.First(), nil
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{DB: db}
}
