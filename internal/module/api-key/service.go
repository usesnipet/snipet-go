package apikey

import (
	"context"
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
	logger     *logger.Logger
	repository repository.IApiKeyRepository
	generator  *auth.APIKeyGenerator
	hasher     *auth.KeyHasher
}

func NewService(
	logger *logger.Logger,
	repository repository.IApiKeyRepository,
	generator *auth.APIKeyGenerator,
	hasher *auth.KeyHasher,
) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
		generator:  generator,
		hasher:     hasher,
	}
}

// VerifyAPIKey is the auth entrypoint itself — looked up by the key's own
// key_id before any tenant is known, called by guard.RequireApiKey. Not
// tenant-scoped; the returned APIKey's TenantID is what the caller uses to
// populate auth.ApiKeyIdentity.
func (s *Service) VerifyAPIKey(ctx context.Context, apiKey string) (*model.APIKey, error) {
	keyID := s.generator.GetKeyID(apiKey)
	paginatedApiKeys, err := s.repository.Filter(
		ctx,
		filter.New[model.APIKey](
			filter.WhereEq("key_id", keyID),
		),
	)
	if err != nil {
		return nil, apperr.Unauthorized("invalid api key")
	}
	if paginatedApiKeys.IsEmpty() {
		return nil, apperr.NotFound("api key not found")
	}

	key := paginatedApiKeys.First()
	if !key.Active {
		return nil, apperr.Forbidden("api key is disabled")
	}
	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, apperr.Forbidden("api key is expired")
	}

	valid, err := s.hasher.Verify(apiKey, key.Key)
	if err != nil || !valid {
		return nil, apperr.Unauthorized("invalid api key")
	}
	return key, nil
}

func (s *Service) Filter(ctx context.Context, tenantID string, opts *filter.Options[model.APIKey]) (*page.Paginated[model.APIKey], error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	return s.repository.Filter(ctx, filter.Merge(opts, filter.New[model.APIKey](filter.WhereEq("tenant_id", tenantID))))
}

// findInTenant fetches by id then verifies the row belongs to tenantID.
func (s *Service) findInTenant(ctx context.Context, tenantID, id string) (*model.APIKey, error) {
	found, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if found.TenantID != tenantID {
		return nil, apperr.NotFound("api key not found")
	}
	return found, nil
}

func (s *Service) FindByID(ctx context.Context, tenantID, id string) (*model.APIKey, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	return s.findInTenant(ctx, tenantID, id)
}

// Me returns the API key that authenticated the current request — inherently
// tenant-agnostic at the call site (there's no tenant_id in the URL on this
// route), it just reads whichever key's identity guard.RequireApiKey set.
func (s *Service) Me(ctx context.Context) (*model.APIKey, error) {
	identity, err := auth.CurrentApiKey(ctx)
	if err != nil {
		return nil, err
	}
	return s.repository.FindByID(ctx, identity.APIKeyID)
}

func (s *Service) Create(ctx context.Context, tenantID string, dto CreateAPIKeyDTO) (*APIKeyWithSecret, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
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

	apiKey := &model.APIKey{
		TenantID:  tenantID,
		Name:      dto.Name,
		KeyID:     keyID,
		Key:       keyHash,
		Active:    true,
		ExpiresAt: dto.ExpiresAt,
	}
	if err := s.repository.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &APIKeyWithSecret{
		APIKey: apiKey,
		Key:    plainKey,
	}, nil
}

func (s *Service) Roll(ctx context.Context, tenantID, id string) (*APIKeyWithSecret, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	if _, err := s.findInTenant(ctx, tenantID, id); err != nil {
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

	updates := &model.APIKey{
		KeyID: keyID,
		Key:   keyHash,
	}
	if err := s.repository.UpdateByID(ctx, id, updates); err != nil {
		return nil, err
	}

	apiKey, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &APIKeyWithSecret{
		APIKey: apiKey,
		Key:    plainKey,
	}, nil
}

func (s *Service) UpdateExpiration(ctx context.Context, tenantID, id string, dto UpdateExpirationDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	return s.repository.UpdateExpiration(ctx, tenantID, id, dto.ExpiresAt)
}

func (s *Service) ToggleActive(ctx context.Context, tenantID, id string, active bool) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	return s.repository.ToggleActive(ctx, tenantID, id, active)
}

func (s *Service) Delete(ctx context.Context, tenantID, id string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	if _, err := s.findInTenant(ctx, tenantID, id); err != nil {
		return err
	}
	return s.repository.DeleteByID(ctx, id)
}
