package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/model"
	auth_provider "github.com/usesnipet/snipet/internal/module/auth/auth-provider"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/pkg/jsonx"
	"gorm.io/gorm"
)

type Service struct {
	registry            *auth_provider.Registry
	clientRepo          repository.IClientRepository
	userRepo            repository.IClientUserRepository
	refreshTokenRepo    repository.IRefreshTokenRepository
	jwtService          *auth.JWTService
	refreshTokenService *auth.RefreshTokenService
	authConfig          config.AuthConfig
}

func NewService(
	registry *auth_provider.Registry,
	clientRepo repository.IClientRepository,
	userRepo repository.IClientUserRepository,
	refreshTokenRepo repository.IRefreshTokenRepository,
	jwtService *auth.JWTService,
	refreshTokenService *auth.RefreshTokenService,
	authConfig config.AuthConfig,
) *Service {
	return &Service{
		registry:            registry,
		clientRepo:          clientRepo,
		userRepo:            userRepo,
		refreshTokenRepo:    refreshTokenRepo,
		jwtService:          jwtService,
		refreshTokenService: refreshTokenService,
		authConfig:          authConfig,
	}
}

func (s *Service) issueTokens(ctx context.Context, clientCode string, user *model.ClientUser, metadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	accessToken, accessTokenExpiresAt, _, err := s.jwtService.GenerateToken(clientCode, user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiresAt, err := s.refreshTokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = jsonx.JSONMap{}
	}

	record := &model.RefreshToken{
		Hash:      s.refreshTokenService.HashRefreshToken(refreshToken),
		UserID:    user.ID,
		ExpiresAt: refreshTokenExpiresAt,
		Metadata:  metadata,
	}
	if err := s.refreshTokenRepo.Create(ctx, record); err != nil {
		return nil, err
	}

	return &AuthenticateResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		User:                  *user,
	}, nil
}

func (s *Service) Authenticate(ctx context.Context, clientCode string, providerName auth_provider.ProviderName, req *http.Request, refreshMetadata jsonx.JSONMap) (*AuthenticateResponse, error) {
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
	return s.issueTokens(ctx, client.Code, user, refreshMetadata)
}

func (s *Service) generateAnonymousName() string {
	return fmt.Sprintf("Anonymous %s", uuid.New().String()[:8])
}

func (s *Service) AuthenticateAnonymous(ctx context.Context, clientCode string, dto AuthenticateAnonymousDTO, refreshMetadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	client, err := s.clientRepo.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	if !client.Config.Anonymous.Enabled {
		return nil, apperr.Unauthorized("anonymous authentication is not enabled for this client")
	}

	name := s.generateAnonymousName()
	if dto.Name != nil {
		name = *dto.Name
	}
	metadata := dto.Metadata
	if metadata == nil {
		metadata = jsonx.JSONMap{}
	}
	user := &model.ClientUser{
		Name:     name,
		Picture:  dto.Picture,
		Email:    dto.Email,
		Metadata: metadata,
	}
	if err := s.userRepo.CreateInClient(ctx, clientCode, user, nil); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, clientCode, user, refreshMetadata)
}

func (s *Service) Refresh(ctx context.Context, clientCode string, dto RefreshDTO, refreshMetadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	hash := s.refreshTokenService.HashRefreshToken(dto.RefreshToken)
	token, err := s.refreshTokenRepo.FindByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.Unauthorized("invalid refresh token")
		}
		return nil, err
	}
	if token.RevokedAt != nil {
		return nil, apperr.Unauthorized("refresh token revoked")
	}
	if token.ExpiresAt.Before(time.Now()) {
		return nil, apperr.Unauthorized("refresh token expired")
	}

	user, err := s.userRepo.FindByIDInClient(ctx, clientCode, token.UserID)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, apperr.Unauthorized("invalid refresh token")
		}
		return nil, err
	}

	if err := s.refreshTokenRepo.RevokeByID(ctx, token.ID); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, clientCode, user, refreshMetadata)
}
