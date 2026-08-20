package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IAppRepository interface {
	IRepository[model.App]
}

type AppRepository struct {
	*Repository[model.App]
}

func NewAppRepository(db *gorm.DB) IAppRepository {
	return &AppRepository{
		Repository: NewRepository[model.App](db),
	}
}
