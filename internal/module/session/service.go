package session

import (
	"context"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	sessionRepo        repository.ISessionRepository
	sessionMessageRepo repository.ISessionMessageRepository
	memoryRepo         repository.IMemoryRepository
}

func NewService(
	sessionRepo repository.ISessionRepository,
	sessionMessageRepo repository.ISessionMessageRepository,
	memoryRepo repository.IMemoryRepository,
) *Service {
	return &Service{
		sessionRepo:        sessionRepo,
		sessionMessageRepo: sessionMessageRepo,
		memoryRepo:         memoryRepo,
	}
}

func (s *Service) Filter(
	ctx context.Context,
	clientCode string,
	filter *filter.Options[model.Session],
) (*page.Paginated[model.Session], error) {
	principal := auth.GetPrincipal(ctx)
	if principal.GetType() == auth.PrincipalTypeJWT {
		return s.sessionRepo.FilterInClientWithUser(
			ctx,
			clientCode,
			principal.GetJWTClaims().Subject,
			filter,
		)
	}
	return s.sessionRepo.FilterInClient(ctx, clientCode, filter)
}

func (s *Service) FindByID(ctx context.Context, clientCode string, id string) (*model.Session, error) {
	return s.sessionRepo.FindByIDInClient(ctx, clientCode, id)
}

func (s *Service) Create(ctx context.Context, clientCode string, dto CreateSessionDTO) (*model.Session, error) {
	memoryUUID, err := uuid.Parse(dto.MemoryID)
	if err != nil {
		return nil, apperr.BadRequest("invalid memory id")
	}
	botUUID, err := uuid.Parse(dto.BotID)
	if err != nil {
		return nil, apperr.BadRequest("invalid bot id")
	}

	memory, err := s.memoryRepo.FindByID(ctx, dto.MemoryID)
	if err != nil {
		return nil, err
	}

	if memory.Type != model.MemoryTypeSession {
		return nil, apperr.BadRequest("invalid memory type")
	}

	session := &model.Session{
		MemoryID: memoryUUID,
		BotID:    botUUID,
		Metadata: dto.Metadata,
	}
	if err := s.sessionRepo.CreateInClient(ctx, clientCode, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) DeleteByID(ctx context.Context, clientCode string, id string) error {
	return s.sessionRepo.DeleteInClient(ctx, clientCode, id)
}

func (s *Service) FindMessages(
	ctx context.Context,
	clientCode string,
	sessionID string,
	filter *filter.Options[model.SessionMessage],
) (*page.Paginated[model.SessionMessage], error) {
	principal := auth.GetPrincipal(ctx)

	if principal.GetType() == auth.PrincipalTypeJWT {
		userID := principal.GetJWTClaims().Subject
		hasAccess, err := s.sessionRepo.CheckUserAccess(ctx, clientCode, userID, sessionID)
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, apperr.Forbidden("user does not have access to this session")
		}
	}

	return s.sessionMessageRepo.FilterInSession(
		ctx,
		clientCode,
		sessionID,
		filter,
	)
}
