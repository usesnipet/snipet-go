package repository

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IApiKeyRepository interface {
	IRepository[model.APIKey]
	UpdateExpiration(ctx context.Context, tenantID, id string, expiresAt *time.Time) error
	ToggleActive(ctx context.Context, tenantID, id string, active bool) error
}

type ApiKeyRepository struct {
	*Repository[model.APIKey]
}

func (r *ApiKeyRepository) UpdateExpiration(ctx context.Context, tenantID, id string, expiresAt *time.Time) error {
	affected, err := gorm.G[model.APIKey](r.db(ctx)).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update(ctx, "expires_at", expiresAt)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("api key not found")
	}
	return nil
}

func (r *ApiKeyRepository) ToggleActive(ctx context.Context, tenantID, id string, active bool) error {
	affected, err := gorm.G[model.APIKey](r.db(ctx)).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update(ctx, "active", active)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("api key not found")
	}
	return nil
}

func NewApiKeyRepository(db *gorm.DB) IApiKeyRepository {
	return &ApiKeyRepository{
		Repository: NewRepository[model.APIKey](db),
	}
}
