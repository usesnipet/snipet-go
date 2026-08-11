package repository

import (
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type ITenantRepository interface {
	IRepository[model.Tenant]
}

type TenantRepository struct {
	*Repository[model.Tenant]
}

func NewTenantRepository(db *gorm.DB) ITenantRepository {
	return &TenantRepository{
		Repository: NewRepository[model.Tenant](db),
	}
}
