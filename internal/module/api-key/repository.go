package apikey

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.APIKey]
	UpdateExpiration(ctx context.Context, id string, expiresAt *time.Time) error
	ToggleActive(ctx context.Context, id string, active bool) error
}

type Repository struct {
	*database.Repository[model.APIKey]
}

func (r *Repository) UpdateExpiration(ctx context.Context, id string, expiresAt *time.Time) error {
	affected, err := gorm.G[model.APIKey](r.DB).
		Where("id = ?", id).
		Update(ctx, "expires_at", expiresAt)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("api key not found")
	}
	return nil
}

func (r *Repository) ToggleActive(ctx context.Context, id string, active bool) error {
	_, err := gorm.G[model.APIKey](r.DB).
		Where("id = ?", id).
		Update(ctx, "active", active)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.APIKey](db),
	}
}
