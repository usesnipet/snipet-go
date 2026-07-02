package session

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/client"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	sessionRepo        repository.ISessionRepository
	sessionMessageRepo repository.ISessionMessageRepository
	clientService      *client.Service
}

func NewService(
	sessionRepo repository.ISessionRepository,
	sessionMessageRepo repository.ISessionMessageRepository,
	clientService *client.Service,
) *Service {
	return &Service{
		sessionRepo:        sessionRepo,
		sessionMessageRepo: sessionMessageRepo,
		clientService:      clientService,
	}
}

func (s *Service) ensureSessionUserAccess(ctx context.Context, clientID string, sessionID string) (*auth.UserClaims, error) {
	principal, ok := auth.GetPrincipal(ctx)
	if !ok {
		return nil, apperr.Unauthorized("unauthorized")
	}
	if principal.GetType() == auth.PrincipalTypeAPIKey {
		return nil, nil
	}

	userId := principal.GetJWTClaims().Subject
	hasAccess, err := s.sessionRepo.CheckUserAccess(ctx, clientID, userId, sessionID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, apperr.Forbidden("user does not have access to this session")
	}
	return principal.GetJWTClaims(), nil
}

func (s *Service) Filter(ctx context.Context, clientCode string, filter *filter.Options[model.Session]) (*page.Paginated[model.Session], error) {
	client, err := s.clientService.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	principal, ok := auth.GetPrincipal(ctx)
	if !ok {
		return nil, apperr.Unauthorized("unauthorized")
	}
	if principal.GetType() == auth.PrincipalTypeAPIKey {
		return s.sessionRepo.FilterInClient(ctx, client.ID, filter)
	}
	return s.sessionRepo.FilterInClientWithUser(
		ctx,
		client.ID,
		principal.GetJWTClaims().Subject,
		filter,
	)
}

func (s *Service) FindByID(ctx context.Context, clientCode string, id string) (*model.Session, error) {
	client, err := s.clientService.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	if _, err := s.ensureSessionUserAccess(ctx, client.ID, id); err != nil {
		return nil, err
	}

	return s.sessionRepo.FindByIDInClient(ctx, client.ID, id)
}

func (s *Service) Create(ctx context.Context, clientCode string, dto CreateSessionDTO) (*model.Session, error) {
	client, err := s.clientService.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		BotID:    dto.BotID,
		Metadata: dto.Metadata,
		ClientID: client.ID,
	}

	principal, ok := auth.GetPrincipal(ctx)
	if !ok {
		return nil, apperr.Unauthorized("unauthorized")
	}
	if principal.GetType() == auth.PrincipalTypeJWT {
		session.UserToSessions = append(session.UserToSessions, model.UserToSession{
			UserID: principal.GetJWTClaims().Subject,
		})
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) DeleteByID(ctx context.Context, clientCode string, id string) error {
	if _, err := s.ensureSessionUserAccess(ctx, clientCode, id); err != nil {
		return err
	}

	return s.sessionRepo.DeleteInClient(ctx, clientCode, id)
}

func (s *Service) FindMessages(
	ctx context.Context,
	clientCode string,
	sessionID string,
	filter *filter.Options[model.SessionMessage],
) (*page.Paginated[model.SessionMessage], error) {
	if _, err := s.ensureSessionUserAccess(ctx, clientCode, sessionID); err != nil {
		return nil, err
	}

	return s.sessionMessageRepo.FilterInSession(
		ctx,
		clientCode,
		sessionID,
		filter,
	)
}

func (s *Service) SendMessage(ctx context.Context, clientCode string, sessionID string, dto SendMessageDTO) error {
	principal, err := s.ensureSessionUserAccess(ctx, clientCode, sessionID)
	if err != nil {
		return err
	}

	message := &model.SessionMessage{
		UserID:    principal.Subject,
		SessionID: sessionID,
		Role:      "user",
		Parts: []model.SessionMessagePart{
			{
				Type:    model.SessionMessagePartTypeText,
				Content: dto.Message,
			},
		},
	}
	return s.sessionMessageRepo.Create(ctx, message)
}
