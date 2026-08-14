package repository

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"gorm.io/gorm"
)

type ITenantInvitationRepository interface {
	IRepository[model.TenantInvitation]
	FindPendingByEmailAndTenant(ctx context.Context, email, tenantID string) (*model.TenantInvitation, error)
	FindByToken(ctx context.Context, token string) (*model.TenantInvitation, error)
}

type TenantInvitationRepository struct {
	*Repository[model.TenantInvitation]
}

func NewTenantInvitationRepository(db *gorm.DB) ITenantInvitationRepository {
	return &TenantInvitationRepository{
		Repository: NewRepository[model.TenantInvitation](db),
	}
}

func (r *TenantInvitationRepository) FindPendingByEmailAndTenant(ctx context.Context, email, tenantID string) (*model.TenantInvitation, error) {
	paginated, err := r.Filter(ctx, filter.New[model.TenantInvitation](
		filter.WhereEq("email", email),
		filter.WhereEq("tenant_id", tenantID),
		filter.WhereEq("status", model.InvitationStatusPending),
		filter.WhereGte("expires_at", time.Now()),
	))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("invitation not found")
	}
	return paginated.First(), nil
}

func (r *TenantInvitationRepository) FindByToken(ctx context.Context, token string) (*model.TenantInvitation, error) {
	paginated, err := r.Filter(ctx, filter.New[model.TenantInvitation](
		filter.WhereEq("token", token),
		filter.Take(1),
	))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("invitation not found")
	}
	return paginated.First(), nil
}
