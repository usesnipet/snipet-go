package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IBotRepository interface {
	IRepository[model.Bot]
	LinkBotToClient(ctx context.Context, clientID, botID uuid.UUID) error
}

type BotRepository struct {
	*Repository[model.Bot]
}

func (r *BotRepository) LinkBotToClient(ctx context.Context, clientID, botID uuid.UUID) error {
	clientBot := &model.ClientToBot{
		ClientID: clientID,
		BotID:    botID,
	}
	err := gorm.G[model.ClientToBot](r.db(ctx)).Create(ctx, clientBot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperr.Conflict("client already linked to bot")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return apperr.NotFound("client or bot not found")
		}
		return err
	}
	return nil
}

func NewBotRepository(db *gorm.DB) IBotRepository {
	return &BotRepository{
		Repository: NewRepository[model.Bot](db),
	}
}
