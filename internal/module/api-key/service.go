package apikey

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	logger     *logger.Logger
	repository repository.IApiKeyRepository
	generator  *APIKeyGenerator
	hasher     *KeyHasher
}

func NewService(logger *logger.Logger, repository repository.IApiKeyRepository) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
		generator:  NewAPIKeyGenerator(),
		hasher:     NewKeyHasher(),
	}
}

func (s *Service) Init(ctx context.Context) error {
	apiKeys, err := s.repository.FilterBy(ctx, filter.Default[model.APIKey]())
	if err != nil {
		return err
	}
	if apiKeys.IsEmpty() {
		s.logger.Infof("no api keys found, creating root api key")
		created, err := s.Create(ctx, CreateAPIKeyDTO{Name: "Root"})
		if err != nil {
			s.logger.Errorf("failed to create root api key: %v", err)
			return err
		}
		s.logger.Infof("root api key created: %s", created.Key)
	}
	return nil
}

func (s *Service) VerifyAPIKey(ctx context.Context, apiKey string) error {
	keyID := s.generator.GetKeyID(apiKey)
	paginatedApiKeys, err := s.repository.FilterBy(
		ctx,
		filter.New[model.APIKey](
			filter.WhereEq("key_id", keyID),
		),
	)
	if err != nil {
		return apperr.Unauthorized("invalid api key")
	}
	if paginatedApiKeys.IsEmpty() {
		return apperr.NotFound("api key not found")
	}

	key := paginatedApiKeys.First()
	if !key.Active {
		return apperr.Forbidden("api key is disabled")
	}
	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return apperr.Forbidden("api key is expired")
	}

	valid, err := s.hasher.Verify(apiKey, key.Key)
	if err != nil || !valid {
		return apperr.Unauthorized("invalid api key")
	}
	return nil
}

func (s *Service) FilterBy(ctx context.Context) (*page.Paginated[model.APIKey], error) {
	return s.repository.FilterBy(ctx, filter.Default[model.APIKey]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.APIKey, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateAPIKeyDTO) (*APIKeyWithSecret, error) {
	plainKey, keyID, err := s.generator.Generate()
	if err != nil {
		return nil, err
	}

	keyHash, err := s.hasher.Hash(plainKey)
	if err != nil {
		return nil, err
	}

	apiKey := &model.APIKey{
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

func (s *Service) Roll(ctx context.Context, id string) (*APIKeyWithSecret, error) {
	if _, err := s.repository.FindByID(ctx, id); err != nil {
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

func (s *Service) UpdateExpiration(ctx context.Context, id string, dto UpdateExpirationDTO) error {
	return s.repository.UpdateExpiration(ctx, id, dto.ExpiresAt)
}

func (s *Service) ToggleActive(ctx context.Context, id string, active bool) error {
	return s.repository.ToggleActive(ctx, id, active)
}
