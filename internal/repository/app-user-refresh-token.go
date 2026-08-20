package repository

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IAppUserRefreshTokenRepository interface {
	Create(ctx context.Context, token *model.AppUserRefreshToken) error
	FindByHash(ctx context.Context, hash string) (*model.AppUserRefreshToken, error)
	RevokeByID(ctx context.Context, id string) error
}

type AppUserRefreshTokenRepository struct {
	*Repository[model.AppUserRefreshToken]
}

func NewAppUserRefreshTokenRepository(db *gorm.DB) IAppUserRefreshTokenRepository {
	return &AppUserRefreshTokenRepository{
		Repository: NewRepository[model.AppUserRefreshToken](db),
	}
}

func (r *AppUserRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*model.AppUserRefreshToken, error) {
	token, err := gorm.G[model.AppUserRefreshToken](r.db(ctx)).
		Where("hash = ?", hash).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AppUserRefreshTokenRepository) RevokeByID(ctx context.Context, id string) error {
	now := time.Now()
	affected, err := gorm.G[model.AppUserRefreshToken](r.db(ctx)).
		Where("id = ?", id).
		Update(ctx, "revoked_at", now)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("refresh token not found")
	}
	return nil
}
