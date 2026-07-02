package indexedknowledgeitem

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	repo repository.IIndexedKnowledgeItemRepository
}

func NewService(repo repository.IIndexedKnowledgeItemRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Filter(ctx context.Context) (*page.Paginated[model.IndexedKnowledgeItem], error) {
	return s.repo.Filter(ctx, filter.Default[model.IndexedKnowledgeItem]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.IndexedKnowledgeItem, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateIndexedKnowledgeItemDTO) (*model.IndexedKnowledgeItem, error) {
	item := &model.IndexedKnowledgeItem{
		Status:          dto.Status,
		Version:         dto.Version,
		Hash:            dto.Hash,
		IndexedAt:       dto.IndexedAt,
		LastError:       dto.LastError,
		Metadata:        dto.Metadata,
		IndexID:         dto.IndexID,
		KnowledgeItemID: dto.KnowledgeItemID,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateIndexedKnowledgeItemDTO) error {
	updates := &model.IndexedKnowledgeItem{}
	if dto.Status != nil {
		updates.Status = *dto.Status
	}
	if dto.Version != nil {
		updates.Version = *dto.Version
	}
	if dto.Hash != nil {
		updates.Hash = *dto.Hash
	}
	if dto.IndexedAt != nil {
		updates.IndexedAt = dto.IndexedAt
	}
	if dto.LastError != nil {
		updates.LastError = *dto.LastError
	}
	if dto.Metadata != nil {
		updates.Metadata = *dto.Metadata
	}
	return s.repo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repo.DeleteByID(ctx, id)
}
