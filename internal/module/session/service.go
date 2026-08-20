package session

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/agent"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/execution"
)

type Service struct {
	sessionRepo          repository.ISessionRepository
	executionMessageRepo repository.IExecutionMessageRepository
	clientService        *client.Service
	agentService         *agent.Service
}

func NewService(
	sessionRepo repository.ISessionRepository,
	executionMessageRepo repository.IExecutionMessageRepository,
	clientService *client.Service,
	agentService *agent.Service,
) *Service {
	return &Service{
		sessionRepo:          sessionRepo,
		executionMessageRepo: executionMessageRepo,
		clientService:        clientService,
		agentService:         agentService,
	}
}

// resolveClient returns the full Client row (not just its ID) — Create
// needs client.TenantID to stamp the new Session's denormalized tenant_id.
func (s *Service) resolveClient(ctx context.Context, clientCode string) (*model.App, error) {
	return s.clientService.FindByCode(ctx, clientCode)
}

func (s *Service) ensureSessionUserAccess(ctx context.Context, clientID string, sessionID string) error {
	if auth.HasApiKey(ctx) {
		return nil
	}

	identity, err := auth.CurrentClientUser(ctx)
	if err != nil {
		return err
	}
	hasAccess, err := s.sessionRepo.CheckUserAccess(ctx, clientID, identity.UserID, sessionID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return apperr.Forbidden("user does not have access to this session")
	}
	return nil
}

func (s *Service) Filter(ctx context.Context, clientCode string, filter *filter.Options[model.Session]) (*page.Paginated[model.Session], error) {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	clientID := resolvedClient.ID
	if auth.HasApiKey(ctx) {
		return s.sessionRepo.FilterInClient(ctx, clientID, filter)
	}
	identity, err := auth.CurrentClientUser(ctx)
	if err != nil {
		return nil, err
	}
	return s.sessionRepo.FilterInClientWithUser(
		ctx,
		clientID,
		identity.UserID,
		filter,
	)
}

func (s *Service) FindByID(
	ctx context.Context,
	clientCode string,
	id string,
	opts *filter.Options[model.Session],
) (*model.Session, error) {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	clientID := resolvedClient.ID
	if err := s.ensureSessionUserAccess(ctx, clientID, id); err != nil {
		return nil, err
	}

	return s.sessionRepo.FindByIDInClient(ctx, clientID, id, opts)
}

func (s *Service) Create(ctx context.Context, clientCode string, dto CreateSessionDTO) (*model.Session, error) {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	clientID := resolvedClient.ID

	session := &model.Session{
		TenantID: resolvedClient.TenantID,
		AgentID:  dto.AgentID,
		Metadata: dto.Metadata,
		ClientID: clientID,
	}

	if identity, err := auth.CurrentClientUser(ctx); err == nil {
		session.ClientUserToSessions = append(session.ClientUserToSessions, model.ClientUserToSession{
			ClientUserID: identity.UserID,
		})
	} else if !auth.HasApiKey(ctx) {
		return nil, apperr.Unauthorized("unauthorized")
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) DeleteByID(ctx context.Context, clientCode string, id string) error {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return err
	}
	clientID := resolvedClient.ID
	if err := s.ensureSessionUserAccess(ctx, clientID, id); err != nil {
		return err
	}

	return s.sessionRepo.DeleteInClient(ctx, clientID, id)
}

func (s *Service) UpdateByID(ctx context.Context, clientCode string, id string, dto UpdateSessionDTO) error {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return err
	}
	clientID := resolvedClient.ID
	if err := s.ensureSessionUserAccess(ctx, clientID, id); err != nil {
		return err
	}
	if _, err := s.sessionRepo.FindByIDInClient(ctx, clientID, id, nil); err != nil {
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
	clientCode string,
	sessionID string,
	filter *filter.Options[model.ExecutionMessage],
) (*page.Paginated[model.ExecutionMessage], error) {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	clientID := resolvedClient.ID
	if err := s.ensureSessionUserAccess(ctx, clientID, sessionID); err != nil {
		return nil, err
	}
	if _, err := s.sessionRepo.FindByIDInClient(ctx, clientID, sessionID, nil); err != nil {
		return nil, err
	}

	return s.executionMessageRepo.FilterInSession(ctx, sessionID, filter)
}

func (s *Service) Run(
	ctx context.Context,
	clientCode string,
	sessionID string,
	dto RunSessionDTO,
	subscribers ...execution.Subscriber,
) error {
	resolvedClient, err := s.resolveClient(ctx, clientCode)
	if err != nil {
		return err
	}
	clientID := resolvedClient.ID
	if err := s.ensureSessionUserAccess(ctx, clientID, sessionID); err != nil {
		return err
	}

	session, err := s.sessionRepo.FindByIDInClient(ctx, clientID, sessionID, nil)
	if err != nil {
		return err
	}

	return s.agentService.Run(ctx, agent.RunInput{
		AgentID:   session.AgentID,
		SessionID: &session.ID,
		Message:   dto.Message,
	}, subscribers...)
}
