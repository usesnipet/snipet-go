package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IClientRepository interface {
	IRepository[model.Client]
}

type ClientRepository struct {
	*Repository[model.Client]
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &ClientRepository{
		Repository: NewRepository[model.Client](db),
	}
}
