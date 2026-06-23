package client

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.Client]
	FindByID(ctx context.Context, id string) (*model.Client, error)
}

type Repository struct {
	*database.Repository[model.Client]
}

func (r *Repository) FindByID(ctx context.Context, id string) (*model.Client, error) {
	paginated, err := r.Repository.FindBy(ctx, filter.New[model.Client](filter.WhereEq("id", id)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("client not found")
	}
	return paginated.First(), nil
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.Client](db),
	}
}
