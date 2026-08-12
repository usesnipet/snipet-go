package repository

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IClientUserRefreshTokenRepository interface {
	Create(ctx context.Context, token *model.ClientUserRefreshToken) error
	FindByHash(ctx context.Context, hash string) (*model.ClientUserRefreshToken, error)
	RevokeByID(ctx context.Context, id string) error
}

type ClientUserRefreshTokenRepository struct {
	*Repository[model.ClientUserRefreshToken]
}

func NewClientUserRefreshTokenRepository(db *gorm.DB) IClientUserRefreshTokenRepository {
	return &ClientUserRefreshTokenRepository{
		Repository: NewRepository[model.ClientUserRefreshToken](db),
	}
}

func (r *ClientUserRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*model.ClientUserRefreshToken, error) {
	token, err := gorm.G[model.ClientUserRefreshToken](r.db(ctx)).
		Where("hash = ?", hash).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *ClientUserRefreshTokenRepository) RevokeByID(ctx context.Context, id string) error {
	now := time.Now()
	affected, err := gorm.G[model.ClientUserRefreshToken](r.db(ctx)).
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
