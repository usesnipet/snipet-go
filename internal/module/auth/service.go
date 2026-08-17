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
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/license"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/module/email"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	config           config.AuthConfig
	userRepo         repository.IUserRepository
	accountRepo      repository.IAccountRepository
	tokenRepo        repository.ITokenRepository
	jwtService       *auth.JWTService[*auth.PlatformUserClaims]
	tokenGen         *auth.TokenService
	emailSvc         *email.Service
	providerRegistry *ProviderRegistry
	license          *license.Service
}

func NewService(
	config config.AuthConfig,
	userRepo repository.IUserRepository,
	accountRepo repository.IAccountRepository,
	tokenRepo repository.ITokenRepository,
	jwtService *auth.JWTService[*auth.PlatformUserClaims],
	tokenGen *auth.TokenService,
	emailSvc *email.Service,
	providerRegistry *ProviderRegistry,
	license *license.Service,
) *Service {
	return &Service{
		config:           config,
		userRepo:         userRepo,
		accountRepo:      accountRepo,
		tokenRepo:        tokenRepo,
		jwtService:       jwtService,
		tokenGen:         tokenGen,
		emailSvc:         emailSvc,
		providerRegistry: providerRegistry,
		license:          license,
	}
}

// Register creates an inactive User (email + password) and emails an
// activate_account token. No access/refresh tokens are issued — login stays
// blocked until Activate is called. Only available on licensed (multi-tenant
// capable) instances: unlicensed single-tenant instances have no self-serve
// path into a Tenant (no invitation flow with just one tenant), so a tenant
// admin creates accounts directly instead (member.Service.Create).
func (s *Service) Register(ctx context.Context, dto RegisterDTO) (*RegisterResponse, error) {
	if !s.license.Info().Valid {
		return nil, apperr.Forbidden("self-registration is disabled for this instance — ask a tenant admin to create your account")
	}

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

	link := fmt.Sprintf("%s/auth/activate?token=%s", strings.TrimRight(s.config.AppURL, "/"), plainToken)
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
		if tokens, err := s.tokenRepo.FindByUserIDAndType(ctx, found.ID, model.TokenTypeActivateAccount); err == nil {
			if len(tokens) == 0 || tokens[0].RevokedAt != nil || tokens[0].ExpiresAt.Before(time.Now()) {
				if err := s.issueActivationToken(ctx, found); err != nil {
					return nil, err
				}
			}
		}
		return nil, apperr.Forbidden("account not activated")
	}

	return s.issueTokens(ctx, found)
}

func (s *Service) issueTokens(ctx context.Context, u *model.User) (*AuthenticateResponse, error) {
	claims := &auth.PlatformUserClaims{BaseClaims: auth.NewBaseClaims(s.config, u.ID)}
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
	p, err := s.providerRegistry.Get(provider)
	if err != nil {
		return "", err
	}

	state, err := s.tokenGen.GenerateToken()
	if err != nil {
		return "", err
	}

	return p.AuthorizationURL(state), nil
}

// AuthenticateCallback exchanges code for the provider identity. A matching
// Account issues tokens for its User directly. No Account but a matching
// email links a new Account to that existing User. Neither found creates
// both — OAuth logins skip the activation gate entirely (no PasswordHash,
// no active_account challenge) since the provider already verified the
// email.
func (s *Service) AuthenticateCallback(ctx context.Context, provider ProviderName, code string) (*AuthenticateResponse, error) {
	p, err := s.providerRegistry.Get(provider)
	if err != nil {
		return nil, err
	}

	identity, err := p.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	account, err := s.accountRepo.FindByProviderAndExternalID(ctx, string(provider), identity.ExternalID)
	if err == nil {
		found, err := s.userRepo.FindByID(ctx, account.UserID)
		if err != nil {
			return nil, err
		}
		return s.issueTokens(ctx, found)
	}
	var appErr *apperr.Error
	if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
		return nil, err
	}

	found, err := s.userRepo.FindByEmail(ctx, identity.Email)
	if err != nil {
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}

		name := identity.Name
		if name == "" {
			name = identity.Email
		}

		found = &model.User{
			Name:       name,
			Email:      identity.Email,
			Challenges: []model.Challenge{},
		}
		if identity.Picture != "" {
			found.Picture = &identity.Picture
		}
		if err := s.userRepo.Create(ctx, found); err != nil {
			return nil, err
		}
	}

	newAccount := &model.Account{
		UserID:     found.ID,
		Provider:   string(provider),
		ExternalID: identity.ExternalID,
	}
	if err := s.accountRepo.Create(ctx, newAccount); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, found)
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

// SetPassword is a "set", not a "change" — no current-password check,
// regardless of whether PasswordHash was already set. Covers both an
// OAuth-only user setting their first password and an existing user
// changing theirs.
func (s *Service) SetPassword(ctx context.Context, dto SetPasswordDTO) error {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return err
	}

	passwordHash, err := auth.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdateByID(ctx, identity.User.ID, &model.User{PasswordHash: &passwordHash})
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
	invalidToken := apperr.BadRequest("invalid or expired activation token")

	hash := s.tokenGen.HashToken(dto.Token)
	token, err := s.tokenRepo.FindByHashAndType(ctx, hash, model.TokenTypeActivateAccount)
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

	found, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return err
	}

	remaining := slices.DeleteFunc(found.Challenges, func(c model.Challenge) bool {
		return c == model.ChallengeActiveAccount
	})
	if remaining == nil {
		remaining = []model.Challenge{}
	}
	if err := s.userRepo.UpdateByID(ctx, found.ID, &model.User{Challenges: remaining}); err != nil {
		return err
	}

	return s.tokenRepo.RevokeByID(ctx, token.ID)
}

// ResendActivation revokes any prior un-consumed activate_account token for
// the user before issuing a new one, so valid tokens don't stack. Silently
// no-ops for an unknown email, same enumeration-safety rationale as
// ForgotPassword.
func (s *Service) ResendActivation(ctx context.Context, dto ResendActivationDTO) error {
	found, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}

	pending, err := s.tokenRepo.Filter(ctx, filter.New[model.Token](
		filter.WhereEq("user_id", found.ID),
		filter.WhereEq("type", model.TokenTypeActivateAccount),
	))
	if err != nil {
		return err
	}
	for _, t := range pending.Data {
		if t.RevokedAt == nil {
			if err := s.tokenRepo.RevokeByID(ctx, t.ID); err != nil {
				return err
			}
		}
	}

	return s.issueActivationToken(ctx, found)
}
