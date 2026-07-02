package knowledgeitem

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	repo repository.IKnowledgeItemRepository
}

func NewService(repo repository.IKnowledgeItemRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Filter(ctx context.Context, knowledgeID string, filter *filter.Options[model.KnowledgeItem]) (*page.Paginated[model.KnowledgeItem], error) {
	return s.repo.FilterInKnowledge(ctx, knowledgeID, filter)
}

func (s *Service) FindByID(ctx context.Context, knowledgeID, id string) (*model.KnowledgeItem, error) {
	return s.repo.FindByIDInKnowledge(ctx, knowledgeID, id)
}

func (s *Service) Create(ctx context.Context, knowledgeID string, dto CreateKnowledgeItemDTO) (*model.KnowledgeItem, error) {
	item := &model.KnowledgeItem{
		ExternalID:   dto.ExternalID,
		Name:         dto.Name,
		Hash:         dto.Hash,
		Metadata:     dto.Metadata,
		LastModified: dto.LastModified,
		KnowledgeID:  knowledgeID,
	}
	if err := s.repo.CreateInKnowledge(ctx, knowledgeID, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, knowledgeID, id string, dto UpdateKnowledgeItemDTO) error {
	updates := &model.KnowledgeItem{}
	if dto.ExternalID != nil {
		updates.ExternalID = *dto.ExternalID
	}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Hash != nil {
		updates.Hash = *dto.Hash
	}
	if dto.Metadata != nil {
		updates.Metadata = *dto.Metadata
	}
	if dto.LastModified != nil {
		updates.LastModified = dto.LastModified
	}
	return s.repo.UpdateInKnowledge(ctx, knowledgeID, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, knowledgeID, id string) error {
	return s.repo.DeleteInKnowledge(ctx, knowledgeID, id)
}
