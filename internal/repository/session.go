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
		clientID string,
		userID string,
		sessionID string,
	) (bool, error)

	FilterInClient(
		ctx context.Context,
		clientID string,
		filter *filter.Options[model.Session],
	) (*page.Paginated[model.Session], error)
	FilterInClientWithUser(
		ctx context.Context,
		clientID string,
		userID string,
		filter *filter.Options[model.Session],
	) (*page.Paginated[model.Session], error)

	FindByIDInClient(
		ctx context.Context,
		clientID string,
		id string,
	) (*model.Session, error)

	DeleteInClient(
		ctx context.Context,
		clientID string,
		ID string,
	) error
}

type SessionRepository struct {
	*Repository[model.Session]
}

func NewSessionRepository(db *gorm.DB, clientRepo IClientRepository) ISessionRepository {
	return &SessionRepository{
		Repository: NewRepository[model.Session](db),
	}
}

func (r *SessionRepository) CheckUserAccess(
	ctx context.Context,
	clientId string,
	userID string,
	sessionID string,
) (bool, error) {
	var total int64
	err := r.db(ctx).Table("user_to_sessions").
		Joins("LEFT JOIN sessions s ON s.id = user_to_sessions.session_id").
		Where("s.client_id = ?", clientId).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Count(&total).Error
	return total > 0, err
}

func (r *SessionRepository) FilterInClient(
	ctx context.Context,
	clientId string,
	sessionFilter *filter.Options[model.Session],
) (*page.Paginated[model.Session], error) {
	return r.Filter(
		ctx,
		filter.Merge(
			sessionFilter,
			filter.New[model.Session](filter.WhereEq("client_id", clientId)),
		),
	)
}
func (r *SessionRepository) FilterInClientWithUser(
	ctx context.Context,
	clientID string,
	userID string,
	sessionFilter *filter.Options[model.Session],
) (*page.Paginated[model.Session], error) {
	if sessionFilter == nil {
		sessionFilter = filter.Default[model.Session]()
	}

	var total int64
	err := r.db(ctx).Table("sessions").
		Joins("LEFT JOIN user_to_sessions uts ON uts.session_id = sessions.id").
		Where("client_id = ?", clientID).
		Where("uts.user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, err
	}
	chain, err := sessionFilter.ToGormTx(r.db(ctx).Table("sessions"))
	if err != nil {
		return nil, err
	}
	var data []model.Session
	err = chain.Joins("LEFT JOIN user_to_sessions uts ON uts.session_id = sessions.id").
		Where("client_id = ?", clientID).
		Where("uts.user_id = ?", userID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return page.NewPaginated(data, total, int64(sessionFilter.Skip), int64(sessionFilter.Take)), err
}

func (r *SessionRepository) FindByIDInClient(
	ctx context.Context,
	clientId string,
	id string,
) (*model.Session, error) {
	paginated, err := r.FilterInClient(
		ctx,
		clientId,
		filter.New[model.Session](filter.WhereEq("id", id)),
	)
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("session not found")
	}
	return paginated.First(), nil
}

func (r *SessionRepository) DeleteInClient(
	ctx context.Context,
	clientId string,
	id string,
) error {
	affected, err := gorm.G[model.Session](r.db(ctx)).Where("client_id = ? AND id = ?", clientId, id).Delete(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.NotFound("session not found")
	}
	return nil
}
