package clientuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	userRepo repository.IClientUserRepository
}

func NewService(userRepo repository.IClientUserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) generateAnonymousName() string {
	return fmt.Sprintf("Anonymous %s", uuid.New().String()[:8])
}

func (s *Service) CreateAnonymous(ctx context.Context, clientCode string, dto CreateAnonymousClientUserDTO) error {
	name := s.generateAnonymousName()
	if dto.Name != nil {
		name = *dto.Name
	}
	user := &model.ClientUser{
		Name:     name,
		Metadata: dto.Metadata,
	}
	if err := s.userRepo.CreateInClient(ctx, clientCode, user, nil); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateAuthenticated(ctx context.Context, clientCode string, dto CreateAuthenticatedClientUserDTO) error {
	user := &model.ClientUser{
		Name:     dto.Name,
		Metadata: dto.Metadata,
		Email:    &dto.Email,
		Picture:  dto.Picture,
	}
	if err := s.userRepo.CreateInClient(ctx, clientCode, user, &dto.ExternalID); err != nil {
		return err
	}
	return nil
}

func (s *Service) FilterInClient(ctx context.Context, clientCode string, filter *filter.Options[model.ClientUser]) (*page.Paginated[model.ClientUser], error) {
	return s.userRepo.FilterInClient(ctx, clientCode, filter)
}

func (s *Service) Me(ctx context.Context, clientCode string) (*model.ClientUser, error) {
	principal, ok := auth.GetPrincipal(ctx)
	if !ok || principal.GetType() != auth.PrincipalTypeJWT || principal.GetJWTClaims() == nil {
		return nil, apperr.Unauthorized("unauthorized")
	}
	claims := principal.GetJWTClaims()
	if claims.ClientCode != clientCode {
		return nil, apperr.Forbidden("client code mismatch")
	}
	return s.userRepo.FindByIDInClient(ctx, clientCode, claims.Subject)
}
