package conversation

import (
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.Conversation]
}

type Repository struct {
	*database.Repository[model.Conversation]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.Conversation](db),
	}
}
