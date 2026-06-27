package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/model"
	auth_provider "github.com/usesnipet/snipet/internal/module/auth/auth-provider"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	registry   *auth_provider.Registry
	clientRepo repository.IClientRepository
	userRepo   repository.IUserRepository
	jwtService *auth.JWTService
}

func NewService(
	registry *auth_provider.Registry,
	clientRepo repository.IClientRepository,
	userRepo repository.IUserRepository,
	jwtService *auth.JWTService,
) *Service {
	return &Service{registry: registry, clientRepo: clientRepo, userRepo: userRepo, jwtService: jwtService}
}

func (s *Service) generateToken(ctx context.Context, clientCode string, user *model.User) (*AuthenticateResponse, error) {
	token, _, err := s.jwtService.GenerateToken(clientCode, user)
	return &AuthenticateResponse{Token: token}, err
}

func (s *Service) Authenticate(ctx context.Context, clientCode string, providerName auth_provider.ProviderName, req *http.Request) (*AuthenticateResponse, error) {
	client, err := s.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	identity, err := s.registry.Authenticate(
		ctx,
		providerName,
		client.Code,
		&client.Config,
		req,
	)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByExternalIDInClient(ctx, clientCode, identity.ExternalID)
	if err != nil {
		return nil, err
	}
	return s.generateToken(ctx, client.Code, user)
}

func (s *Service) generateAnonymousName() string {
	return fmt.Sprintf("Anonymous %s", uuid.New().String()[:8])
}

func (s *Service) AuthenticateAnonymous(ctx context.Context, clientCode string, dto AuthenticateAnonymousDTO) (*AuthenticateResponse, error) {
	name := s.generateAnonymousName()
	if dto.Name != nil {
		name = *dto.Name
	}
	user := &model.User{
		Name:     name,
		Picture:  dto.Picture,
		Email:    dto.Email,
		Metadata: dto.Metadata,
	}
	if err := s.userRepo.CreateInClient(ctx, clientCode, user, nil); err != nil {
		return nil, err
	}
	return s.generateToken(ctx, clientCode, user)
}
