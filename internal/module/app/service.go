package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/authz"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	appRepo   repository.IAppRepository
	generator *auth.APIKeyGenerator
	hasher    *auth.KeyHasher
	logger    *logger.Logger
}

func NewService(
	appRepo repository.IAppRepository,
	generator *auth.APIKeyGenerator,
	hasher *auth.KeyHasher,
	logger *logger.Logger,
) *Service {
	return &Service{appRepo: appRepo, generator: generator, hasher: hasher, logger: logger}
}

func (s *Service) Filter(ctx context.Context, tenantID string, opts *filter.Options[model.App]) (*page.Paginated[model.App], error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	return s.appRepo.Filter(ctx, filter.Merge(opts, filter.New[model.App](filter.WhereEq("tenant_id", tenantID))))
}

// FindByCode is the tenant-agnostic lookup used internally (e.g. Ping) where
// there's no tenant-staff caller to scope against.
func (s *Service) FindByCode(ctx context.Context, code string) (*model.App, error) {
	return s.appRepo.FindByCode(ctx, code)
}

// findByCodeInTenant is the admin-path lookup — verifies the app belongs
// to tenantID (404, not 403, to avoid confirming the code exists elsewhere).
func (s *Service) findByCodeInTenant(ctx context.Context, tenantID, code string) (*model.App, error) {
	found, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if found.TenantID != tenantID {
		return nil, apperr.NotFound("app not found")
	}
	return found, nil
}

func (s *Service) FindByCodeInTenant(ctx context.Context, tenantID, code string) (*model.App, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	return s.findByCodeInTenant(ctx, tenantID, code)
}

// FindPublicByCode is the unauthenticated lookup for an app's public info
// (e.g. a frontend-only widget looking itself up by code). Only apps marked
// Public are exposed; a non-public app 404s the same as one that doesn't
// exist, so callers can't use this to probe for private app codes.
func (s *Service) FindPublicByCode(ctx context.Context, code string) (*PublicAppDTO, error) {
	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if !app.Public {
		return nil, apperr.NotFound("app not found")
	}
	return &PublicAppDTO{Code: app.Code, Name: app.Name, Description: app.Description}, nil
}

func (s *Service) GenerateCode(ctx context.Context) (string, error) {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}
	code := base64.RawURLEncoding.EncodeToString(randomBytes)[:10] // first 10 characters

	// check if code is already in use
	_, err := s.FindByCode(ctx, code)
	var notFoundError *apperr.Error
	if errors.As(err, &notFoundError) && notFoundError.StatusCode == http.StatusNotFound {
		return code, nil
	}
	if err != nil {
		return "", err
	}
	return s.GenerateCode(ctx)
}

func (s *Service) Create(ctx context.Context, tenantID string, dto CreateAppDTO) (*AppWithSecret, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	code, err := s.GenerateCode(ctx)
	if err != nil {
		return nil, err
	}

	plainKey, keyID, err := s.generator.Generate()
	if err != nil {
		return nil, err
	}
	keyHash, err := s.hasher.Hash(plainKey)
	if err != nil {
		return nil, err
	}

	app := &model.App{
		TenantID:    tenantID,
		Code:        code,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      model.AppStatusPending,
		Public:      dto.Public,
		KeyID:       keyID,
		KeyHash:     keyHash,
	}
	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}

	return &AppWithSecret{App: app, Key: plainKey}, nil
}

func (s *Service) Roll(ctx context.Context, tenantID, code string) (*AppWithSecret, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	if _, err := s.findByCodeInTenant(ctx, tenantID, code); err != nil {
		return nil, err
	}

	plainKey, keyID, err := s.generator.Generate()
	if err != nil {
		return nil, err
	}
	keyHash, err := s.hasher.Hash(plainKey)
	if err != nil {
		return nil, err
	}

	updates := &model.App{KeyID: keyID, KeyHash: keyHash}
	if err := s.appRepo.UpdateByCode(ctx, code, updates); err != nil {
		return nil, err
	}

	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &AppWithSecret{App: app, Key: plainKey}, nil
}

func (s *Service) UpdateByCode(ctx context.Context, tenantID, code string, dto UpdateAppDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	if _, err := s.findByCodeInTenant(ctx, tenantID, code); err != nil {
		return err
	}

	updates := &model.App{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Description != nil {
		updates.Description = *dto.Description
	}
	if err := s.appRepo.UpdateByCode(ctx, code, updates); err != nil {
		return err
	}

	if dto.Public != nil {
		return s.appRepo.UpdatePublicByCode(ctx, code, *dto.Public)
	}
	return nil
}

func (s *Service) UpdateAuthConfig(ctx context.Context, tenantID, code string, dto UpdateAppAuthConfigDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	if _, err := s.findByCodeInTenant(ctx, tenantID, code); err != nil {
		return err
	}
	return s.appRepo.UpdateAuthConfigByCode(ctx, code, dto.ToModel())
}

// ToggleActive flips an app between active and deactivated. A deactivated
// app fails key verification (see VerifyKey) so it can no longer ping in or
// authenticate.
func (s *Service) ToggleActive(ctx context.Context, tenantID, code string, active bool) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	app, err := s.findByCodeInTenant(ctx, tenantID, code)
	if err != nil {
		return err
	}
	if app.Status == model.AppStatusPending && active && !app.Public {
		return apperr.BadRequest("cannot activate a pending app")
	}
	status := model.AppStatusActive
	if !active {
		status = model.AppStatusDeactivated
	}
	return s.appRepo.UpdateByCode(ctx, code, &model.App{Status: status})
}

func (s *Service) DeleteByCode(ctx context.Context, tenantID, code string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	if _, err := s.findByCodeInTenant(ctx, tenantID, code); err != nil {
		return err
	}
	return s.appRepo.DeleteByCode(ctx, code)
}

// VerifyKey is the auth entrypoint itself — looked up by the key's own
// key_id before any tenant is known, called by guard.RequireAppKey. A
// deactivated app is rejected here so it can't authenticate at all.
func (s *Service) VerifyKey(ctx context.Context, key string) (*model.App, error) {
	keyID := s.generator.GetKeyID(key)
	paginated, err := s.appRepo.Filter(ctx, filter.New[model.App](filter.WhereEq("key_id", keyID)))
	if err != nil {
		return nil, apperr.Unauthorized("invalid app key")
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("app not found")
	}

	app := paginated.First()
	if app.Status == model.AppStatusDeactivated {
		return nil, apperr.Forbidden("app is disabled")
	}

	valid, err := s.hasher.Verify(key, app.KeyHash)
	if err != nil || !valid {
		return nil, apperr.Unauthorized("invalid app key")
	}
	return app, nil
}

// Ping is how a connected app confirms it's alive: it authenticates with its
// own key (guard.RequireAppKey), then hits this by its code. A pending app
// is promoted to active on its first successful ping.
func (s *Service) Ping(ctx context.Context, code string) error {
	identity, err := auth.CurrentAppKey(ctx)
	if err != nil {
		return err
	}
	if identity.Code != code {
		return apperr.Forbidden("app key does not match code")
	}

	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := &model.App{LastVerifiedAt: &now}
	if app.Status == model.AppStatusPending {
		updates.Status = model.AppStatusActive
	}
	return s.appRepo.UpdateByCode(ctx, code, updates)
}
