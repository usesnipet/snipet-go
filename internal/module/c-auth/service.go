package c_auth

import (
	"context"
	"net/http"

	"github.com/usesnipet/snipet/internal/auth"
	auth_provider "github.com/usesnipet/snipet/internal/module/c-auth/auth-provider"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	registry   *auth_provider.Registry
	clientRepo repository.IClientRepository
	cUserRepo  repository.ICUserRepository
	jwtService *auth.JWTService
}

func NewService(registry *auth_provider.Registry, clientRepo repository.IClientRepository) *Service {
	return &Service{registry: registry, clientRepo: clientRepo}
}

func (s *Service) generateToken(ctx context.Context, clientCode string, externalID string) (string, error) {
	user, err := s.cUserRepo.FindByExternalIDInClient(ctx, clientCode, externalID)
	if err != nil {
		return "", err
	}

	token, _, err := s.jwtService.GenerateToken(clientCode, user)
	return token, err
}

func (s *Service) Authenticate(ctx context.Context, clientCode string, providerName auth_provider.ProviderName, req *http.Request) (AuthenticateResponse, error) {
	client, err := s.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return AuthenticateResponse{}, err
	}
	identity, err := s.registry.Authenticate(
		ctx,
		providerName,
		client.Code,
		&client.Config,
		req,
	)
	if err != nil {
		return AuthenticateResponse{}, err
	}
	token, err := s.generateToken(ctx, client.Code, identity.ExternalID)
	if err != nil {
		return AuthenticateResponse{}, err
	}
	return AuthenticateResponse{Token: token}, nil
}
