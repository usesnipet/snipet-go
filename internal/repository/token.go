package repository

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type ITokenRepository interface {
	IRepository[model.Token]
	FindByHashAndType(ctx context.Context, hash string, tokenType model.TokenType) (*model.Token, error)
	RevokeByID(ctx context.Context, id string) error
	FindByUserIDAndType(ctx context.Context, userID string, tokenType model.TokenType) ([]model.Token, error)
}

type TokenRepository struct {
	*Repository[model.Token]
}

func NewTokenRepository(db *gorm.DB) ITokenRepository {
	return &TokenRepository{
		Repository: NewRepository[model.Token](db),
	}
}

func (r *TokenRepository) FindByHashAndType(ctx context.Context, hash string, tokenType model.TokenType) (*model.Token, error) {
	paginated, err := r.Filter(ctx, filter.New[model.Token](
		filter.WhereEq("hash", hash),
		filter.WhereEq("type", tokenType),
	))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("token not found")
	}
	return paginated.First(), nil
}

func (r *TokenRepository) RevokeByID(ctx context.Context, id string) error {
	now := time.Now()
	affected, err := gorm.G[model.Token](r.db(ctx)).
		Where("id = ?", id).
		Update(ctx, "revoked_at", now)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("token not found")
	}
	return nil
}

func (r *TokenRepository) FindByUserIDAndType(ctx context.Context, userID string, tokenType model.TokenType) ([]model.Token, error) {
	paginated, err := r.Filter(ctx, filter.New[model.Token](
		filter.WhereEq("user_id", userID),
		filter.WhereEq("type", tokenType),
	))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("token not found")
	}
	return paginated.Data, nil
}
