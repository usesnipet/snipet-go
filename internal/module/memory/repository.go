package memory

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IRepository interface {
	database.IRepository[model.Memory]
	SetAsDefault(ctx context.Context, memoryID string) error
}

type Repository struct {
	*database.Repository[model.Memory]
}

func (r *Repository) SetAsDefault(ctx context.Context, memoryID string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[model.Memory](tx).Where("is_default = ?", true).Update(ctx, "is_default", false)
		if err != nil {
			return err
		}
		affected, err := gorm.G[model.Memory](tx).Where("id = ?", memoryID).Update(ctx, "is_default", true)
		if err != nil {
			if affected == 0 {
				return apperr.NotFound("memory not found")
			}
			return err
		}
		return nil
	})
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		Repository: database.NewRepository[model.Memory](db),
	}
}
