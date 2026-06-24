package bot

import (
	"context"
	"errors"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.Bot]
	LinkClientToBot(ctx context.Context, clientID, botID uuid.UUID) error
}

type Repository struct {
	*database.Repository[model.Bot]
}

func (r *Repository) LinkClientToBot(ctx context.Context, clientID, botID uuid.UUID) error {
	clientBot := &model.ClientBot{
		ClientID: clientID,
		BotID:    botID,
	}
	err := gorm.G[model.ClientBot](r.DB).Create(ctx, clientBot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperr.Conflict("client already linked to bot")
		}
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.Bot](db),
	}
}
