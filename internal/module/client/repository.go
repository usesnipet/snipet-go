package client

import (
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.Client]
}

type Repository struct {
	*database.Repository[model.Client]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.Client](db),
	}
}
