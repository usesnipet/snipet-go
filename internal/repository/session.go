package repository

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"gorm.io/gorm"
)

type ISessionRepository interface {
	IRepository[model.Session]
	CheckUserAccess(
		ctx context.Context,
		appID string,
		userID string,
		sessionID string,
	) (bool, error)

	FilterInApp(
		ctx context.Context,
		appID string,
		filter *filter.Options[model.Session],
	) (*page.Paginated[model.Session], error)
	FilterInAppWithUser(
		ctx context.Context,
		appID string,
		userID string,
		filter *filter.Options[model.Session],
	) (*page.Paginated[model.Session], error)

	FindByIDInApp(
		ctx context.Context,
		appID string,
		id string,
		opts *filter.Options[model.Session],
	) (*model.Session, error)

	DeleteInApp(
		ctx context.Context,
		appID string,
		ID string,
	) error
}

type SessionRepository struct {
	*Repository[model.Session]
}

func NewSessionRepository(db *gorm.DB, appRepo IAppRepository) ISessionRepository {
	return &SessionRepository{
		Repository: NewRepository[model.Session](db),
	}
}

func (r *SessionRepository) CheckUserAccess(
	ctx context.Context,
	appID string,
	userID string,
	sessionID string,
) (bool, error) {
	var total int64
	err := r.db(ctx).Table("app_user_to_sessions").
		Joins("LEFT JOIN sessions s ON s.id = app_user_to_sessions.session_id").
		Where("s.app_id = ?", appID).
		Where("app_user_id = ? AND session_id = ?", userID, sessionID).
		Count(&total).Error
	return total > 0, err
}

func (r *SessionRepository) FilterInApp(
	ctx context.Context,
	appID string,
	sessionFilter *filter.Options[model.Session],
) (*page.Paginated[model.Session], error) {
	return r.Filter(
		ctx,
		filter.Merge(
			sessionFilter,
			filter.New[model.Session](filter.WhereEq("app_id", appID)),
		),
	)
}
func (r *SessionRepository) FilterInAppWithUser(
	ctx context.Context,
	appID string,
	userID string,
	sessionFilter *filter.Options[model.Session],
) (*page.Paginated[model.Session], error) {
	if sessionFilter == nil {
		sessionFilter = filter.Default[model.Session]()
	}

	var total int64
	err := r.db(ctx).Table("sessions").
		Joins("LEFT JOIN app_user_to_sessions uts ON uts.session_id = sessions.id").
		Where("app_id = ?", appID).
		Where("uts.app_user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, err
	}
	chain, err := sessionFilter.ToGormTx(r.db(ctx).Model(&model.Session{}))
	if err != nil {
		return nil, err
	}
	var data []model.Session
	err = chain.Joins("LEFT JOIN app_user_to_sessions uts ON uts.session_id = sessions.id").
		Where("app_id = ?", appID).
		Where("uts.app_user_id = ?", userID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return page.NewPaginated(data, total, int64(sessionFilter.Skip), int64(sessionFilter.Take)), err
}

func (r *SessionRepository) FindByIDInApp(
	ctx context.Context,
	appID string,
	id string,
	opts *filter.Options[model.Session],
) (*model.Session, error) {
	paginated, err := r.FilterInApp(
		ctx,
		appID,
		filter.Merge(
			opts,
			filter.New[model.Session](filter.WhereEq("id", id)),
		),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("session not found")
	}
	return paginated.First(), nil
}

func (r *SessionRepository) DeleteInApp(
	ctx context.Context,
	appID string,
	id string,
) error {
	affected, err := gorm.G[model.Session](r.db(ctx)).Where("app_id = ? AND id = ?", appID, id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("session not found")
	}
	return nil
}
