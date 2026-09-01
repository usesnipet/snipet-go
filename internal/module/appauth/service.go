package appauth

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
	auth_provider "github.com/usesnipet/snipet/internal/module/appauth/auth-provider"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/pkg/jsonx"
	"gorm.io/gorm"
)

type Service struct {
	registry             *auth_provider.Registry
	appRepo              repository.IAppRepository
	userRepo             repository.IAppUserRepository
	userRefreshTokenRepo repository.IAppUserRefreshTokenRepository
	jwtService           *auth.JWTService[*auth.AppUserClaims]
	tokenService         *auth.TokenService
	authConfig           config.AuthConfig
}

func NewService(
	registry *auth_provider.Registry,
	appRepo repository.IAppRepository,
	userRepo repository.IAppUserRepository,
	userRefreshTokenRepo repository.IAppUserRefreshTokenRepository,
	jwtService *auth.JWTService[*auth.AppUserClaims],
	refreshTokenService *auth.TokenService,
	authConfig config.AuthConfig,
) *Service {
	return &Service{
		registry:             registry,
		appRepo:              appRepo,
		userRepo:             userRepo,
		userRefreshTokenRepo: userRefreshTokenRepo,
		jwtService:           jwtService,
		tokenService:         refreshTokenService,
		authConfig:           authConfig,
	}
}

func (s *Service) issueTokens(ctx context.Context, appCode string, user *model.AppUser, metadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	claims := &auth.AppUserClaims{
		BaseClaims: auth.NewBaseClaims(s.authConfig, user.ID),
		AppCode:    appCode,
	}
	accessToken, accessTokenExpiresAt, err := s.jwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = jsonx.JSONMap{}
	}

	record := &model.AppUserRefreshToken{
		Hash:      s.tokenService.HashToken(refreshToken),
		AppUserID: user.ID,
		ExpiresAt: time.Now().Add(s.authConfig.RefreshTokenExpiration),
		Metadata:  metadata,
	}
	if err := s.userRefreshTokenRepo.Create(ctx, record); err != nil {
		return nil, err
	}

	return &AuthenticateResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: time.Now().Add(s.authConfig.RefreshTokenExpiration),
		User:                  *user,
	}, nil
}

func (s *Service) Authenticate(ctx context.Context, appCode string, providerName auth_provider.ProviderName, req *http.Request, refreshMetadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	app, err := s.appRepo.FindByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}
	identity, err := s.registry.Authenticate(
		ctx,
		providerName,
		app.Code,
		&app.AuthConfig,
		req,
	)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByExternalIDInApp(ctx, appCode, identity.ExternalID)
	if err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, app.Code, user, refreshMetadata)
}

func (s *Service) generateAnonymousName() string {
	return fmt.Sprintf("Anonymous %s", uuid.New().String()[:8])
}

func (s *Service) AuthenticateAnonymous(ctx context.Context, appCode string, dto AuthenticateAnonymousDTO, refreshMetadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	app, err := s.appRepo.FindByCode(ctx, appCode)
	if err != nil {
		return nil, err
	}
	if !app.AuthConfig.Anonymous.Enabled {
		return nil, apperr.Unauthorized("anonymous authentication is not enabled for this app")
	}

	name := s.generateAnonymousName()
	if dto.Name != nil {
		name = *dto.Name
	}
	metadata := dto.Metadata
	if metadata == nil {
		metadata = jsonx.JSONMap{}
	}
	user := &model.AppUser{
		Name:     name,
		Picture:  dto.Picture,
		Email:    dto.Email,
		Metadata: metadata,
	}
	if err := s.userRepo.CreateInApp(ctx, appCode, user, nil); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, appCode, user, refreshMetadata)
}

func (s *Service) Refresh(ctx context.Context, appCode string, dto RefreshDTO, refreshMetadata jsonx.JSONMap) (*AuthenticateResponse, error) {
	hash := s.tokenService.HashToken(dto.RefreshToken)
	token, err := s.userRefreshTokenRepo.FindByHash(ctx, hash)
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

	user, err := s.userRepo.FindByIDInApp(ctx, appCode, token.AppUserID)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, apperr.Unauthorized("invalid refresh token")
		}
		return nil, err
	}

	if err := s.userRefreshTokenRepo.RevokeByID(ctx, token.ID); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, appCode, user, refreshMetadata)
}
