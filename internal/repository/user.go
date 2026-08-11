package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	IRepository[model.User]
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository struct {
	*Repository[model.User]
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		Repository: NewRepository[model.User](db),
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	paginated, err := r.Filter(ctx, filter.New[model.User](filter.WhereEq("email", email)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("user not found")
	}
	return paginated.First(), nil
}
