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

func (s *Service) Filter(ctx context.Context, knowledgeID, indexID string, filter *filter.Options[model.IndexedKnowledgeItem]) (*page.Paginated[model.IndexedKnowledgeItem], error) {
	return s.repo.FilterInIndex(ctx, knowledgeID, indexID, filter)
}

func (s *Service) FindByID(ctx context.Context, knowledgeID, indexID, id string) (*model.IndexedKnowledgeItem, error) {
	return s.repo.FindByIDInIndex(ctx, knowledgeID, indexID, id)
}

func (s *Service) Create(ctx context.Context, knowledgeID, indexID string, dto CreateIndexedKnowledgeItemDTO) (*model.IndexedKnowledgeItem, error) {
	item := &model.IndexedKnowledgeItem{
		Status:          dto.Status,
		Version:         dto.Version,
		Hash:            dto.Hash,
		IndexedAt:       dto.IndexedAt,
		LastError:       dto.LastError,
		Metadata:        dto.Metadata,
		IndexID:         indexID,
		KnowledgeItemID: dto.KnowledgeItemID,
	}
	if err := s.repo.CreateInIndex(ctx, knowledgeID, indexID, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, knowledgeID, indexID, id string, dto UpdateIndexedKnowledgeItemDTO) error {
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
	return s.repo.UpdateInIndex(ctx, knowledgeID, indexID, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, knowledgeID, indexID, id string) error {
	return s.repo.DeleteInIndex(ctx, knowledgeID, indexID, id)
}
