package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IBotRepository interface {
	IRepository[model.Bot]
}

type BotRepository struct {
	*Repository[model.Bot]
}

func NewBotRepository(db *gorm.DB) IBotRepository {
	return &BotRepository{
		Repository: NewRepository[model.Bot](db),
	}
}
