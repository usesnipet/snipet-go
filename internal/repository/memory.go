package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type IMemoryRepository interface {
	IRepository[model.Memory]
	SetAsDefault(ctx context.Context, memoryType model.MemoryType, memoryID string) error
}

type MemoryRepository struct {
	*Repository[model.Memory]
}

func (r *MemoryRepository) SetAsDefault(ctx context.Context, memoryType model.MemoryType, memoryID string) error {
	return WithTransaction(ctx, r.DB, func(ctx context.Context) error {
		db := r.db(ctx)
		_, err := gorm.G[model.Memory](db).Where("type = ? AND is_default = ?", memoryType, true).Update(ctx, "is_default", false)
		if err != nil {
			return err
		}
		affected, err := gorm.G[model.Memory](db).Where("type = ? AND id = ?", memoryType, memoryID).Update(ctx, "is_default", true)
		if err != nil {
			return err
		}
		if affected == 0 {
			return apperr.NotFound("memory not found")
		}
		return nil
	})
}

func NewMemoryRepository(db *gorm.DB) IMemoryRepository {
	return &MemoryRepository{
		Repository: NewRepository[model.Memory](db),
	}
}
