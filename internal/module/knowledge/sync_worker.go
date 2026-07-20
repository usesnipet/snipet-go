package knowledge

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	knowledgeindex "github.com/usesnipet/snipet/internal/module/knowledge-index"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/driver"
)

type SyncKnowledgeArgs struct {
	KnowledgeID string `json:"knowledge_id"`
	Force       bool   `json:"force"`
}

func (SyncKnowledgeArgs) Kind() string {
	return "knowledge-sync"
}

type SyncResult struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Deleted int64 `json:"deleted"`
	Failed  int64 `json:"failed"`
}

type SyncWorker struct {
	river.WorkerDefaults[SyncKnowledgeArgs]

	txManager          repository.ITxManager
	sourceManager      *driver.Manager[driver.IKnowledgeSource]
	knowledgeRepo      repository.IKnowledgeRepository
	knowledgeItemRepo  repository.IKnowledgeItemRepository
	knowledgeIndexRepo repository.IKnowledgeIndexRepository
	batchSize          int
	logger             *logger.Logger
}

func NewSyncWorker(
	txManager repository.ITxManager,
	sourceManager *driver.Manager[driver.IKnowledgeSource],
	knowledgeRepo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	knowledgeIndexRepo repository.IKnowledgeIndexRepository,
	batchSize int,
	logger *logger.Logger,
) *SyncWorker {
	return &SyncWorker{
		txManager:          txManager,
		sourceManager:      sourceManager,
		knowledgeRepo:      knowledgeRepo,
		knowledgeItemRepo:  knowledgeItemRepo,
		knowledgeIndexRepo: knowledgeIndexRepo,
		batchSize:          batchSize,
		logger:             logger,
	}
}

func (s *SyncWorker) Work(ctx context.Context, job *river.Job[SyncKnowledgeArgs]) error {
	knowledgeID := job.Args.KnowledgeID
	force := job.Args.Force

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
		if err = s.knowledgeRepo.UpdateByID(ctx, knowledgeID, kn); err != nil {
			s.logger.Errorf("failed to update knowledge sync status for knowledge %s: %s", knowledgeID, err.Error())
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

	indexed, err := s.knowledgeIndexRepo.FilterInKnowledge(
		ctx,
		knowledgeID,
		filter.Default[model.KnowledgeIndex](),
	)
	if err != nil {
		s.logger.Errorf("failed to get indexed knowledge items for knowledge %s: %s", knowledgeID, err.Error())
		return err
	}

	s.logger.Verbosef(
		"knowledge sync completed for knowledge %s: c=%d, u=%d, d=%d, f=%d",
		knowledgeID, result.Created, result.Updated, result.Deleted, result.Failed,
	)

	return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		var err error
		for _, index := range indexed.Data {
			_, err = queue.PushFromContext(
				ctx,
				knowledgeindex.SyncIndexArgs{IndexID: index.ID, KnowledgeID: knowledgeID},
				nil,
			)
			if err != nil {
				s.logger.Errorf("failed to push sync index job for index %s: %s", index.ID, err.Error())
				break
			}
		}
		return err
	})
}
