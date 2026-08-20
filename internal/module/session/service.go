package session

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/agent"
	appmodule "github.com/usesnipet/snipet/internal/module/app"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/execution"
)

type Service struct {
	sessionRepo          repository.ISessionRepository
	executionMessageRepo repository.IExecutionMessageRepository
	appService           *appmodule.Service
	agentService         *agent.Service
}

func NewService(
	sessionRepo repository.ISessionRepository,
	executionMessageRepo repository.IExecutionMessageRepository,
	appService *appmodule.Service,
	agentService *agent.Service,
) *Service {
	return &Service{
		sessionRepo:          sessionRepo,
		executionMessageRepo: executionMessageRepo,
		appService:           appService,
		agentService:         agentService,
	}
}

// resolveApp returns the full App row (not just its ID) — Create needs
// app.TenantID to stamp the new Session's denormalized tenant_id.
func (s *Service) resolveApp(ctx context.Context, appCode string) (*model.App, error) {
	app, err := s.appService.FindByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}

	if auth.HasAppKey(ctx) {
		appIdentity, err := auth.CurrentAppKey(ctx)
		if err != nil {
			return nil, err
		}

		if appIdentity.AppID != app.ID {
			return nil, apperr.Forbidden("this app key does not have access to this app")
		}
		return app, nil
	}
	if auth.HasAppUser(ctx) {
		userIdentity, err := auth.CurrentAppUser(ctx)
		if err != nil {
			return nil, err
		}
		if userIdentity.AppCode != app.Code {
			return nil, apperr.Forbidden("you do not have access to this app")
		}
		return app, nil
	}
	return nil, apperr.Unauthorized("you must be logged to access this resource")
}

func (s *Service) ensureSessionAccess(ctx context.Context, appID string, sessionID string) error {
	if auth.HasAppKey(ctx) {
		appIdentity, err := auth.CurrentAppKey(ctx)
		if err != nil {
			return err
		}
		if appIdentity.AppID != appID {
			return apperr.Forbidden("app does not have access to this session")
		}
		return nil
	}

	userIdentity, err := auth.CurrentAppUser(ctx)
	if err != nil {
		return err
	}
	hasAccess, err := s.sessionRepo.CheckUserAccess(ctx, appID, userIdentity.UserID, sessionID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return apperr.Forbidden("user does not have access to this session")
	}
	return nil
}

func (s *Service) Filter(ctx context.Context, appCode string, filter *filter.Options[model.Session]) (*page.Paginated[model.Session], error) {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return nil, err
	}
	appID := resolvedApp.ID
	if auth.HasAppKey(ctx) {
		return s.sessionRepo.FilterInApp(ctx, appID, filter)
	}
	identity, err := auth.CurrentAppUser(ctx)
	if err != nil {
		return nil, err
	}
	return s.sessionRepo.FilterInAppWithUser(
		ctx,
		appID,
		identity.UserID,
		filter,
	)
}

func (s *Service) FindByID(
	ctx context.Context,
	appCode string,
	id string,
	opts *filter.Options[model.Session],
) (*model.Session, error) {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return nil, err
	}
	appID := resolvedApp.ID
	if err := s.ensureSessionAccess(ctx, appID, id); err != nil {
		return nil, err
	}

	return s.sessionRepo.FindByIDInApp(ctx, appID, id, opts)
}

func (s *Service) Create(ctx context.Context, appCode string, dto CreateSessionDTO) (*model.Session, error) {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return nil, err
	}
	appID := resolvedApp.ID

	session := &model.Session{
		TenantID: resolvedApp.TenantID,
		AgentID:  dto.AgentID,
		Metadata: dto.Metadata,
		AppID:    appID,
	}

	if identity, err := auth.CurrentAppUser(ctx); err == nil {
		session.AppUserToSessions = append(session.AppUserToSessions, model.AppUserToSession{
			AppUserID: identity.UserID,
		})
	} else if !auth.HasAppKey(ctx) {
		return nil, apperr.Unauthorized("unauthorized")
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) DeleteByID(ctx context.Context, appCode string, id string) error {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return err
	}
	appID := resolvedApp.ID
	if err := s.ensureSessionAccess(ctx, appID, id); err != nil {
		return err
	}

	return s.sessionRepo.DeleteInApp(ctx, appID, id)
}

func (s *Service) UpdateByID(ctx context.Context, appCode string, id string, dto UpdateSessionDTO) error {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return err
	}
	appID := resolvedApp.ID
	if err := s.ensureSessionAccess(ctx, appID, id); err != nil {
		return err
	}
	if _, err := s.sessionRepo.FindByIDInApp(ctx, appID, id, nil); err != nil {
		return err
	}

	updates := &model.Session{}
	if dto.AgentID != nil {
		updates.AgentID = *dto.AgentID
	}
	if dto.Metadata != nil {
		updates.Metadata = dto.Metadata
	}
	return s.sessionRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) FindMessages(
	ctx context.Context,
	appCode string,
	sessionID string,
	filter *filter.Options[model.ExecutionMessage],
) (*page.Paginated[model.ExecutionMessage], error) {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return nil, err
	}
	appID := resolvedApp.ID
	if err := s.ensureSessionAccess(ctx, appID, sessionID); err != nil {
		return nil, err
	}
	if _, err := s.sessionRepo.FindByIDInApp(ctx, appID, sessionID, nil); err != nil {
		return nil, err
	}

	return s.executionMessageRepo.FilterInSession(ctx, sessionID, filter)
}

func (s *Service) Run(
	ctx context.Context,
	appCode string,
	sessionID string,
	dto RunSessionDTO,
	subscribers ...execution.Subscriber,
) error {
	resolvedApp, err := s.resolveApp(ctx, appCode)
	if err != nil {
		return err
	}
	appID := resolvedApp.ID
	if err := s.ensureSessionAccess(ctx, appID, sessionID); err != nil {
		return err
	}

	session, err := s.sessionRepo.FindByIDInApp(ctx, appID, sessionID, nil)
	if err != nil {
		return err
	}

	return s.agentService.Run(ctx, agent.RunInput{
		AgentID:   session.AgentID,
		SessionID: &session.ID,
		Message:   dto.Message,
	}, subscribers...)
}
