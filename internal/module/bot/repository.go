package bot

import (
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.Bot]
}

type Repository struct {
	*database.Repository[model.Bot]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.Bot](db),
	}
}
