package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	IRepository[model.Account]
	FindByProviderAndExternalID(ctx context.Context, provider, externalID string) (*model.Account, error)
}

type AccountRepository struct {
	*Repository[model.Account]
}

func NewAccountRepository(db *gorm.DB) IAccountRepository {
	return &AccountRepository{
		Repository: NewRepository[model.Account](db),
	}
}

func (r *AccountRepository) FindByProviderAndExternalID(ctx context.Context, provider, externalID string) (*model.Account, error) {
	paginated, err := r.Filter(ctx, filter.New[model.Account](
		filter.WhereEq("provider", provider),
		filter.WhereEq("external_id", externalID),
	))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("account not found")
	}
	return paginated.First(), nil
}
