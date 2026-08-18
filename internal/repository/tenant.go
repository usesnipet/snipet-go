package repository

import (
	"context"

	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type ITenantRepository interface {
	IRepository[model.Tenant]

	FilterByUserIDWithMember(ctx context.Context, userID string) (*page.Paginated[model.Tenant], error)
}

type TenantRepository struct {
	*Repository[model.Tenant]
}

func NewTenantRepository(db *gorm.DB) ITenantRepository {
	return &TenantRepository{
		Repository: NewRepository[model.Tenant](db),
	}
}

func (r *TenantRepository) FilterByUserIDWithMember(
	ctx context.Context,
	userID string,
) (*page.Paginated[model.Tenant], error) {
	user, err := gorm.G[model.User](r.db(ctx)).Where("id = ?", userID).Preload("Members.Tenant", nil).First(ctx)
	if err != nil {
		return nil, err
	}

	tenants := make([]model.Tenant, 0, len(user.Members))
	for _, member := range user.Members {
		tenant := *member.Tenant
		member.Tenant = nil
		tenant.Members = append(tenant.Members, member)
		tenants = append(tenants, tenant)
	}
	return page.NewPaginated(tenants, int64(len(tenants)), 0, int64(len(tenants))), nil
}
