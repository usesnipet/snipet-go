package knowledge

import (
	"context"

	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
)

type SyncResult struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Deleted int64 `json:"deleted"`
	Failed  int64 `json:"failed"`
}

type SyncService struct {
	sourceManager     *runtime.SourceManager
	knowledgeRepo     repository.IKnowledgeRepository
	knowledgeItemRepo repository.IKnowledgeItemRepository
	batchSize         int
}

func NewSyncService(
	sourceManager *runtime.SourceManager,
	knowledgeRepo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	batchSize int,
) *SyncService {
	return &SyncService{
		sourceManager:     sourceManager,
		knowledgeRepo:     knowledgeRepo,
		knowledgeItemRepo: knowledgeItemRepo,
		batchSize:         batchSize,
	}
}

func (s *SyncService) Sync(ctx context.Context, knowledgeID string) (*SyncResult, error) {
	knowledge, err := s.knowledgeRepo.FindByID(ctx, knowledgeID)
	if err != nil {
		return nil, err
	}

	sourceDriver, err := s.sourceManager.Prepare(ctx, knowledge.Driver, knowledge.Configuration)
	if err != nil {
		return nil, err
	}

	it, err := sourceDriver.Iterator(ctx, knowledge.Configuration)
	if err != nil {
		return nil, err
	}
	defer it.Close()

	result := &SyncResult{}
	items := []model.KnowledgeItem{}

	for it.Next(ctx) {
		item := it.Item()
		hash := it.GetHash()
		items = append(items, model.KnowledgeItem{
			KnowledgeID:  knowledgeID,
			ExternalID:   item.ID,
			Name:         item.Name,
			Hash:         hash,
			Metadata:     item.Metadata,
			LastModified: item.LastModified,
		})
		if len(items) > s.batchSize {
			if err := s.knowledgeItemRepo.UpsertMany(ctx, items); err != nil {
				result.Failed += int64(s.batchSize)
			}
			items = []model.KnowledgeItem{}
		}
	}
	if err := it.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
