package knowledge

import (
	"context"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	kdriver "github.com/usesnipet/snipet/pkg/driver/knowledge"
)

type SyncResult struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Deleted int64 `json:"deleted"`
	Failed  int64 `json:"failed"`
}

// IndexSyncFunc runs an index sync for a knowledge source after items are synced.
type IndexSyncFunc func(ctx context.Context, knowledgeID, indexID string) error

type SyncWorker struct {
	sourceManager      *manager.Driver[kdriver.ISourceDriver]
	knowledgeRepo      repository.IKnowledgeRepository
	knowledgeItemRepo  repository.IKnowledgeItemRepository
	knowledgeIndexRepo repository.IKnowledgeIndexRepository
	indexSync          IndexSyncFunc
	batchSize          int
	logger             *logger.Logger
}

func NewSyncWorker(
	sourceManager *manager.Driver[kdriver.ISourceDriver],
	knowledgeRepo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	knowledgeIndexRepo repository.IKnowledgeIndexRepository,
	indexSync IndexSyncFunc,
	batchSize int,
	logger *logger.Logger,
) *SyncWorker {
	return &SyncWorker{
		sourceManager:      sourceManager,
		knowledgeRepo:      knowledgeRepo,
		knowledgeItemRepo:  knowledgeItemRepo,
		knowledgeIndexRepo: knowledgeIndexRepo,
		indexSync:          indexSync,
		batchSize:          batchSize,
		logger:             logger,
	}
}

func (s *SyncWorker) Sync(ctx context.Context, knowledgeID string, force bool) (err error) {
	if force {
		s.logger.Verbosef("syncing knowledge %s with force", knowledgeID)
	} else {
		s.logger.Verbosef("syncing knowledge %s", knowledgeID)
	}

	knowledge, err := s.knowledgeRepo.FindByID(ctx, knowledgeID)
	if err != nil {
		return err
	}

	if knowledge.SyncStatus == model.SyncStatusInProgress {
		return apperr.Conflict("knowledge sync already in progress")
	}
	err = s.knowledgeRepo.UpdateByID(ctx, knowledgeID, &model.Knowledge{SyncStatus: model.SyncStatusInProgress})
	if err != nil {
		return err
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
		if updateErr := s.knowledgeRepo.UpdateByID(ctx, knowledgeID, kn); updateErr != nil {
			s.logger.Errorf("failed to update knowledge sync status for knowledge %s: %s", knowledgeID, updateErr.Error())
		}
		s.logger.Verbosef(
			"knowledge sync status updated for knowledge %s: status=%s, last_synced_at=%s, sync_error=%s",
			knowledgeID, status, now.Format(time.RFC3339), errMessage,
		)
	}()

	sourceDriver, err := s.sourceManager.Prepare(ctx, knowledge.Driver, knowledge.Configuration)
	if err != nil {
		return err
	}

	existingHashes, err := s.knowledgeItemRepo.HashesByExternalIDInKnowledge(ctx, knowledgeID)
	if err != nil {
		return err
	}

	it, err := sourceDriver.Iterator(ctx, knowledge.Configuration)
	if err != nil {
		return err
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
		if exists && previousHash == hash && !force {
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
			Kind:         item.Kind,
			Attributes:   item.Attributes,
			Metadata:     item.Metadata,
			LastModified: item.LastModified,
		})

		if len(batch) >= s.batchSize {
			flush()
		}
	}
	if err := it.Err(); err != nil {
		return err
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
			return err
		}
		result.Deleted = deleted
	}

	s.logger.Verbosef(
		"knowledge sync completed for knowledge %s: c=%d, u=%d, d=%d, f=%d",
		knowledgeID, result.Created, result.Updated, result.Deleted, result.Failed,
	)

	if s.indexSync == nil {
		return nil
	}

	indexed, err := s.knowledgeIndexRepo.FilterInKnowledge(
		ctx,
		knowledgeID,
		filter.Default[model.KnowledgeIndex](),
	)
	if err != nil {
		s.logger.Errorf("failed to get indexed knowledge items for knowledge %s: %s", knowledgeID, err.Error())
		return err
	}

	for _, index := range indexed.Data {
		if syncErr := s.indexSync(ctx, knowledgeID, index.ID); syncErr != nil {
			s.logger.Errorf("failed to sync index %s for knowledge %s: %s", index.ID, knowledgeID, syncErr.Error())
		}
	}
	return nil
}
