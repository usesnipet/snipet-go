package knowledgeindex

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	repo                     repository.IKnowledgeIndexRepository
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository
}

func NewService(
	repo repository.IKnowledgeIndexRepository,
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository,
) *Service {
	return &Service{
		repo:                     repo,
		indexedKnowledgeItemRepo: indexedKnowledgeItemRepo,
	}
}

func (s *Service) Filter(ctx context.Context, knowledgeID string, filter *filter.Options[model.KnowledgeIndex]) (*page.Paginated[model.KnowledgeIndex], error) {
	return s.repo.FilterInKnowledge(ctx, knowledgeID, filter)
}

func (s *Service) FilterItems(ctx context.Context, knowledgeID, indexID string, filter *filter.Options[model.IndexedKnowledgeItem]) (*page.Paginated[model.IndexedKnowledgeItem], error) {
	return s.indexedKnowledgeItemRepo.FilterInIndex(ctx, knowledgeID, indexID, filter)
}

func (s *Service) FindByID(ctx context.Context, knowledgeID, id string) (*model.KnowledgeIndex, error) {
	return s.repo.FindByIDInKnowledge(ctx, knowledgeID, id)
}

func (s *Service) Create(ctx context.Context, knowledgeID string, dto CreateKnowledgeIndexDTO) (*model.KnowledgeIndex, error) {
	index := &model.KnowledgeIndex{
		Name:          dto.Name,
		Driver:        dto.Driver,
		Configuration: dto.Configuration,
		KnowledgeID:   knowledgeID,
	}
	if err := s.repo.CreateInKnowledge(ctx, knowledgeID, index); err != nil {
		return nil, err
	}
	return index, nil
}

func (s *Service) Update(ctx context.Context, knowledgeID, id string, dto UpdateKnowledgeIndexDTO) error {
	if dto.Name != nil {
		return s.repo.UpdateNameInKnowledge(ctx, knowledgeID, id, *dto.Name)
	}
	return nil
}

func (s *Service) DeleteByID(ctx context.Context, knowledgeID, id string) error {
	return s.repo.DeleteInKnowledge(ctx, knowledgeID, id)
}
