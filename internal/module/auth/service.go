package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/email"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	config      config.AuthConfig
	userRepo    repository.IUserRepository
	accountRepo repository.IAccountRepository
	tokenRepo   repository.ITokenRepository
	jwtService  *auth.JWTService[*UserClaims]
	tokenGen    *auth.TokenService
	emailSvc    *email.Service
}

func NewService(
	config config.AuthConfig,
	userRepo repository.IUserRepository,
	accountRepo repository.IAccountRepository,
	tokenRepo repository.ITokenRepository,
	jwtService *auth.JWTService[*UserClaims],
	tokenGen *auth.TokenService,
	emailSvc *email.Service,
) *Service {
	return &Service{
		config:      config,
		userRepo:    userRepo,
		accountRepo: accountRepo,
		tokenRepo:   tokenRepo,
		jwtService:  jwtService,
		tokenGen:    tokenGen,
		emailSvc:    emailSvc,
	}
}

// Register creates an inactive User (email + password) and emails an
// activate_account token. No access/refresh tokens are issued — login stays
// blocked until Activate is called.
func (s *Service) Register(ctx context.Context, dto RegisterDTO) (*RegisterResponse, error) {
	if _, err := s.userRepo.FindByEmail(ctx, dto.Email); err == nil {
		return nil, apperr.Conflict("email already in use")
	} else {
		var appErr *apperr.Error
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	passwordHash, err := auth.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	created := &model.User{
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: &passwordHash,
		Challenges:   []model.Challenge{model.ChallengeActiveAccount},
	}
	if err := s.userRepo.Create(ctx, created); err != nil {
		return nil, err
	}

	if err := s.issueActivationToken(ctx, created); err != nil {
		return nil, err
	}

	return &RegisterResponse{User: *created}, nil
}

// issueActionToken creates a one-time DB-backed token (activate_account,
// reset_password, ...) and returns its plaintext for embedding in an email
// link. Only the SHA-256 hash is persisted.
func (s *Service) issueActionToken(ctx context.Context, u *model.User, tokenType model.TokenType, expiration time.Duration) (string, error) {
	plainToken, err := s.tokenGen.GenerateToken()
	if err != nil {
		return "", err
	}

	token := &model.Token{
		Type:      tokenType,
		Hash:      s.tokenGen.HashToken(plainToken),
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(expiration),
	}
	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return "", err
	}
	return plainToken, nil
}

func (s *Service) issueActivationToken(ctx context.Context, u *model.User) error {
	plainToken, err := s.issueActionToken(ctx, u, model.TokenTypeActivateAccount, s.config.ActivateAccountTokenExpiration)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/activate?token=%s", strings.TrimRight(s.config.AppURL, "/"), plainToken)
	return s.emailSvc.SendTemplate(ctx, u.Email, email.TemplateActivateAccount, email.ActivateAccountData{
		Name: u.Name,
		Link: link,
	})
}

// Login rejects with a generic "invalid credentials" error for unknown
// emails and OAuth-only accounts (PasswordHash == nil), without leaking
// which. Only once the password itself checks out is the more specific
// "account not activated" error safe to return.
func (s *Service) Login(ctx context.Context, dto LoginDTO) (*AuthenticateResponse, error) {
	invalidCredentials := apperr.Unauthorized("invalid credentials")

	found, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil || found.PasswordHash == nil {
		return nil, invalidCredentials
	}
	if err := auth.ComparePassword(*found.PasswordHash, dto.Password); err != nil {
		return nil, invalidCredentials
	}

	if slices.Contains(found.Challenges, model.ChallengeActiveAccount) {
		return nil, apperr.Forbidden("account not activated")
	}

	return s.issueTokens(ctx, found)
}

func (s *Service) issueTokens(ctx context.Context, u *model.User) (*AuthenticateResponse, error) {
	claims := &UserClaims{BaseClaims: auth.NewBaseClaims(s.config, u.ID)}
	accessToken, accessTokenExpiresAt, err := s.jwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}

	refreshTokenExpiresAt := time.Now().Add(s.config.RefreshTokenExpiration)
	record := &model.Token{
		Type:      model.TokenTypeRefresh,
		Hash:      s.tokenGen.HashToken(refreshToken),
		UserID:    u.ID,
		ExpiresAt: refreshTokenExpiresAt,
	}
	if err := s.tokenRepo.Create(ctx, record); err != nil {
		return nil, err
	}

	return &AuthenticateResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		User:                  *u,
	}, nil
}

func (s *Service) GetAuthorizationURL(ctx context.Context, provider ProviderName) (string, error) {
	panic("not implemented")
}

func (s *Service) AuthenticateCallback(ctx context.Context, provider ProviderName, code string) (*AuthenticateResponse, error) {
	panic("not implemented")
}

// findRefreshToken resolves a plaintext refresh token to its Token record,
// collapsing "not found" into the same generic error Refresh/Logout return
// for any other invalid-token case.
func (s *Service) findRefreshToken(ctx context.Context, plainToken string) (*model.Token, error) {
	hash := s.tokenGen.HashToken(plainToken)
	token, err := s.tokenRepo.FindByHashAndType(ctx, hash, model.TokenTypeRefresh)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, apperr.Unauthorized("invalid refresh token")
		}
		return nil, err
	}
	return token, nil
}

func (s *Service) Refresh(ctx context.Context, dto RefreshDTO) (*AuthenticateResponse, error) {
	token, err := s.findRefreshToken(ctx, dto.RefreshToken)
	if err != nil {
		return nil, err
	}
	if token.RevokedAt != nil || token.ExpiresAt.Before(time.Now()) {
		return nil, apperr.Unauthorized("invalid refresh token")
	}

	found, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, apperr.Unauthorized("invalid refresh token")
		}
		return nil, err
	}

	if err := s.tokenRepo.RevokeByID(ctx, token.ID); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, found)
}

func (s *Service) Logout(ctx context.Context, dto RefreshDTO) error {
	token, err := s.findRefreshToken(ctx, dto.RefreshToken)
	if err != nil {
		return err
	}
	return s.tokenRepo.RevokeByID(ctx, token.ID)
}

// currentUserID resolves the subject of the bearer JWT on ctx.
func (s *Service) currentUserID(ctx context.Context) (string, error) {
	principal, ok := auth.GetPrincipal(ctx)
	if !ok || principal.GetType() != auth.PrincipalTypeJWT || principal.GetJWTClaims() == nil {
		return "", apperr.Unauthorized("unauthorized")
	}
	subject, err := principal.GetJWTClaims().GetSubject()
	if err != nil || subject == "" {
		return "", apperr.Unauthorized("unauthorized")
	}
	return subject, nil
}

// SetPassword is a "set", not a "change" — no current-password check,
// regardless of whether PasswordHash was already set. Covers both an
// OAuth-only user setting their first password and an existing user
// changing theirs.
func (s *Service) SetPassword(ctx context.Context, dto SetPasswordDTO) error {
	userID, err := s.currentUserID(ctx)
	if err != nil {
		return err
	}

	passwordHash, err := auth.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdateByID(ctx, userID, &model.User{PasswordHash: &passwordHash})
}

// ForgotPassword always succeeds regardless of whether the email exists, to
// avoid leaking which emails are registered.
func (s *Service) ForgotPassword(ctx context.Context, dto ForgotPasswordDTO) error {
	found, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}

	plainToken, err := s.issueActionToken(ctx, found, model.TokenTypeResetPassword, s.config.ResetPasswordTokenExpiration)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimRight(s.config.AppURL, "/"), plainToken)
	return s.emailSvc.SendTemplate(ctx, found.Email, email.TemplateResetPassword, email.ResetPasswordData{
		Name: found.Name,
		Link: link,
	})
}

func (s *Service) ResetPassword(ctx context.Context, dto ResetPasswordDTO) error {
	invalidToken := apperr.BadRequest("invalid or expired reset token")

	hash := s.tokenGen.HashToken(dto.Token)
	token, err := s.tokenRepo.FindByHashAndType(ctx, hash, model.TokenTypeResetPassword)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return invalidToken
		}
		return err
	}
	if token.RevokedAt != nil || token.ExpiresAt.Before(time.Now()) {
		return invalidToken
	}

	passwordHash, err := auth.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}
	if err := s.userRepo.UpdateByID(ctx, token.UserID, &model.User{PasswordHash: &passwordHash}); err != nil {
		return err
	}

	return s.tokenRepo.RevokeByID(ctx, token.ID)
}

func (s *Service) Activate(ctx context.Context, dto ActivateAccountDTO) error {
	panic("not implemented")
}

func (s *Service) ResendActivation(ctx context.Context, dto ResendActivationDTO) error {
	panic("not implemented")
}
