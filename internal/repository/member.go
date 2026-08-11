package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type IMemberRepository interface {
	IRepository[model.Member]
	FindByUserAndTenant(ctx context.Context, userID, tenantID string) (*model.Member, error)
	FilterWithUser(ctx context.Context, opts *filter.Options[model.Member]) (*page.Paginated[model.Member], error)
}

type MemberRepository struct {
	*Repository[model.Member]
}

func NewMemberRepository(db *gorm.DB) IMemberRepository {
	return &MemberRepository{
		Repository: NewRepository[model.Member](db),
	}
}

func (r *MemberRepository) FindByUserAndTenant(ctx context.Context, userID, tenantID string) (*model.Member, error) {
	paginated, err := r.Filter(ctx, filter.New[model.Member](
		filter.WhereEq("user_id", userID),
		filter.WhereEq("tenant_id", tenantID),
	))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("member not found")
	}
	return paginated.First(), nil
}

func (r *MemberRepository) FilterWithUser(ctx context.Context, opts *filter.Options[model.Member]) (*page.Paginated[model.Member], error) {
	return r.Filter(ctx, filter.Merge(opts, filter.New[model.Member](filter.Include("User"))))
}
