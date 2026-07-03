package knowledge

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/logger"
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
	logger            *logger.Logger
}

func NewSyncService(
	sourceManager *runtime.SourceManager,
	knowledgeRepo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	batchSize int,
	logger *logger.Logger,
) *SyncService {
	return &SyncService{
		sourceManager:     sourceManager,
		knowledgeRepo:     knowledgeRepo,
		knowledgeItemRepo: knowledgeItemRepo,
		batchSize:         batchSize,
		logger:            logger,
	}
}

func (s *SyncService) Sync(ctx context.Context, knowledgeID string) (*SyncResult, error) {
	knowledge, err := s.knowledgeRepo.FindByID(ctx, knowledgeID)
	if err != nil {
		return nil, err
	}

	if knowledge.SyncStatus == model.SyncStatusInProgress {
		return nil, apperr.Conflict("knowledge sync already in progress")
	}
	err = s.knowledgeRepo.UpdateByID(ctx, knowledgeID, &model.Knowledge{SyncStatus: model.SyncStatusInProgress})
	if err != nil {
		return nil, err
	}
	defer func() {
		var status model.SyncStatus
		var errMessage string
		now := time.Now()
		if err == nil {
			status = model.SyncStatusSuccess
		} else {
			status = model.SyncStatusFailed
			errMessage = err.Error()
		}

		kn := &model.Knowledge{SyncStatus: status, LastSyncedAt: &now, SyncError: &errMessage}
		if err = s.knowledgeRepo.UpdateByID(ctx, knowledgeID, kn); err != nil {
			s.logger.Errorf("failed to update knowledge sync status for knowledge %s: %s", knowledgeID, err.Error())
		}
		s.logger.Debugf(
			"knowledge sync status updated for knowledge %s: status=%s, last_synced_at=%s, sync_error=%s",
			knowledgeID, status, now.Format(time.RFC3339), errMessage,
		)
	}()

	sourceDriver, err := s.sourceManager.Prepare(ctx, knowledge.Driver, knowledge.Configuration)
	if err != nil {
		return nil, err
	}

	existingHashes, err := s.knowledgeItemRepo.HashesByExternalIDInKnowledge(ctx, knowledgeID)
	if err != nil {
		return nil, err
	}

	it, err := sourceDriver.Iterator(ctx, knowledge.Configuration)
	if err != nil {
		return nil, err
	}
	defer it.Close()

	result := &SyncResult{}
	seen := make(map[string]struct{}, len(existingHashes))

	batch := make([]model.KnowledgeItem, 0, s.batchSize)
	var pendingCreated, pendingUpdated int64

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := s.knowledgeItemRepo.UpsertMany(ctx, batch, s.batchSize); err != nil {
			result.Failed += int64(len(batch))
		} else {
			result.Created += pendingCreated
			result.Updated += pendingUpdated
		}
		batch = batch[:0]
		pendingCreated = 0
		pendingUpdated = 0
	}

	for it.Next(ctx) {
		item := it.Item()
		hash := it.GetHash()
		seen[item.ID] = struct{}{}

		previousHash, exists := existingHashes[item.ID]
		if exists && previousHash == hash {
			continue
		}

		if exists {
			pendingUpdated++
		} else {
			pendingCreated++
		}

		batch = append(batch, model.KnowledgeItem{
			KnowledgeID:  knowledgeID,
			ExternalID:   item.ID,
			Name:         item.Name,
			Hash:         hash,
			Metadata:     item.Metadata,
			LastModified: item.LastModified,
		})

		if len(batch) >= s.batchSize {
			flush()
		}
	}
	if err := it.Err(); err != nil {
		return nil, err
	}
	flush()

	deletedExternalIDs := make([]string, 0)
	for externalID := range existingHashes {
		if _, ok := seen[externalID]; !ok {
			deletedExternalIDs = append(deletedExternalIDs, externalID)
		}
	}
	if len(deletedExternalIDs) > 0 {
		deleted, err := s.knowledgeItemRepo.DeleteByExternalIDsInKnowledge(ctx, knowledgeID, deletedExternalIDs)
		if err != nil {
			return nil, err
		}
		result.Deleted = deleted
	}

	return result, nil
}
