package apikey

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
)

type Service struct {
	repository IRepository
	generator  *APIKeyGenerator
	hasher     *KeyHasher
}

func (s *Service) FindBy(ctx context.Context) (*database.Paginated[model.APIKey], error) {
	return s.repository.FindBy(ctx, filter.Default[model.APIKey]())
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

func NewService(repository IRepository) *Service {
	return &Service{
		repository: repository,
		generator:  NewAPIKeyGenerator(),
		hasher:     NewKeyHasher(),
	}
}
