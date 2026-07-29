package knowledgeindex

import (
	"context"
	"fmt"
	"slices"

	"github.com/riverqueue/river"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
	kdriver "github.com/usesnipet/snipet/pkg/driver/knowledge"
)

type SyncIndexArgs struct {
	KnowledgeID string `json:"knowledge_id"`
	IndexID     string `json:"index_id"`
}

func (SyncIndexArgs) Kind() string {
	return "sync_index"
}

type SyncIndexResult struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Deleted int64 `json:"deleted"`
	Failed  int64 `json:"failed"`
}

type SyncIndexWorker struct {
	river.WorkerDefaults[SyncIndexArgs]

	sourceManager            *runtime.DriverManager[kdriver.ISourceDriver]
	indexManager             *runtime.DriverManager[kdriver.IIndexDriver]
	knowledgeRepo            repository.IKnowledgeRepository
	knowledgeItemRepo        repository.IKnowledgeItemRepository
	indexRepo                repository.IKnowledgeIndexRepository
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository
	logger                   *logger.Logger
}

func NewSyncIndexWorker(
	indexManager *runtime.DriverManager[kdriver.IIndexDriver],
	sourceManager *runtime.DriverManager[kdriver.ISourceDriver],
	knowledgeRepo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	indexRepo repository.IKnowledgeIndexRepository,
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository,
	logger *logger.Logger,
) *SyncIndexWorker {
	return &SyncIndexWorker{
		indexManager:             indexManager,
		sourceManager:            sourceManager,
		knowledgeRepo:            knowledgeRepo,
		knowledgeItemRepo:        knowledgeItemRepo,
		indexRepo:                indexRepo,
		indexedKnowledgeItemRepo: indexedKnowledgeItemRepo,
		logger:                   logger,
	}
}

func (s *SyncIndexWorker) Work(ctx context.Context, job *river.Job[SyncIndexArgs]) error {
	knowledge, err := s.knowledgeRepo.FindByID(ctx, job.Args.KnowledgeID)
	if err != nil {
		return err
	}

	index, err := s.indexRepo.FindByIDInKnowledge(ctx, job.Args.KnowledgeID, job.Args.IndexID)
	if err != nil {
		return err
	}

	indexDriver, err := s.indexManager.Prepare(ctx, index.Driver, index.Configuration)
	if err != nil {
		return err
	}

	sourceDriver, err := s.sourceManager.Prepare(ctx, knowledge.Driver, knowledge.Configuration)
	if err != nil {
		return err
	}

	writer, err := indexDriver.Writer(index.Configuration)
	if err != nil {
		return err
	}
	defer writer.Close()

	toCreate, err := s.indexedKnowledgeItemRepo.FindToCreateInIndex(ctx, job.Args.KnowledgeID, job.Args.IndexID)
	if err != nil {
		return err
	}
	toUpdate, err := s.indexedKnowledgeItemRepo.FindToUpdateInIndex(ctx, job.Args.KnowledgeID, job.Args.IndexID)
	if err != nil {
		return err
	}
	toDelete, err := s.indexedKnowledgeItemRepo.FindToDeleteInIndex(ctx, job.Args.KnowledgeID, job.Args.IndexID)
	if err != nil {
		return err
	}

	var indexedToCreate []*model.IndexedKnowledgeItem
	for _, item := range toCreate {
		indexedToCreate = append(indexedToCreate, &model.IndexedKnowledgeItem{
			IndexID:         job.Args.IndexID,
			KnowledgeItemID: &item.ID,
			Hash:            item.Hash,
			Metadata:        item.Metadata,
			Status:          model.IndexedStatusPending,
		})
	}
	if len(indexedToCreate) > 0 {
		if err := s.indexedKnowledgeItemRepo.CreateManyInIndex(ctx, job.Args.KnowledgeID, job.Args.IndexID, indexedToCreate); err != nil {
			return err
		}
	}

	toUpdateIDs := util.Map(toUpdate, func(item model.IndexedKnowledgeItem) string {
		return item.ID
	})
	toDeleteIDs := util.Map(toDelete, func(item model.IndexedKnowledgeItem) string {
		return item.ID
	})
	pendingIDs := append(toUpdateIDs, toDeleteIDs...)
	if len(pendingIDs) > 0 {
		err = s.indexedKnowledgeItemRepo.UpdateStatusesByIDsInIndex(
			ctx,
			job.Args.KnowledgeID,
			job.Args.IndexID,
			pendingIDs,
			model.IndexedStatusPending,
		)
		if err != nil {
			return err
		}
	}

	if len(toDeleteIDs) > 0 {
		if err := writer.DeleteMany(ctx, toDeleteIDs); err != nil {
			return err
		}
		for _, id := range toDeleteIDs {
			if err := s.indexedKnowledgeItemRepo.DeleteInIndex(ctx, job.Args.KnowledgeID, job.Args.IndexID, id); err != nil {
				return err
			}
		}
	}

	for _, item := range indexedToCreate {
		if err := s.indexItem(ctx, knowledge, index, writer, sourceDriver, *item); err != nil {
			s.logger.Errorf("failed to index item %s: %v", item.ID, err)
		}
	}

	for _, item := range toUpdate {
		if err := s.indexItem(ctx, knowledge, index, writer, sourceDriver, item); err != nil {
			s.logger.Errorf("failed to reindex item %s: %v", item.ID, err)
		}
	}

	return nil
}

func (s *SyncIndexWorker) indexItem(
	ctx context.Context,
	knowledge *model.Knowledge,
	index *model.KnowledgeIndex,
	writer kdriver.IKnowledgeIndexWriter,
	sourceDriver kdriver.ISourceDriver,
	item model.IndexedKnowledgeItem,
) error {
	updateStatus := func(status model.IndexStatus, hash string, reason *string, errMessage *string) error {
		update := &model.IndexedKnowledgeItem{
			Status:    status,
			Reason:    reason,
			LastError: errMessage,
		}
		if hash != "" {
			update.Hash = hash
		}
		return s.indexedKnowledgeItemRepo.UpdateInIndex(
			ctx,
			knowledge.ID,
			index.ID,
			item.ID,
			update,
		)
	}

	if item.KnowledgeItemID == nil {
		reason := "knowledge item reference is missing"
		return updateStatus(model.IndexedStatusError, "", &reason, nil)
	}

	knowledgeItem, err := s.knowledgeItemRepo.FindByIDInKnowledge(ctx, knowledge.ID, *item.KnowledgeItemID)
	if err != nil {
		errMessage := fmt.Sprintf("failed to load knowledge item: %v", err)
		return updateStatus(model.IndexedStatusError, "", nil, &errMessage)
	}

	if err := updateStatus(model.IndexedStatusSyncing, "", nil, nil); err != nil {
		return err
	}

	reader, err := sourceDriver.Reader(ctx, knowledge.Configuration, knowledgeItem.ExternalID)
	if err != nil {
		errMessage := fmt.Sprintf("failed to get content: %v", err)
		return updateStatus(model.IndexedStatusError, "", nil, &errMessage)
	}
	defer reader.Close()

	if !slices.Contains(writer.SupportedKinds(), reader.Kind()) {
		reason := fmt.Sprintf("item kind (%s) is not supported by index", reader.Kind())
		return updateStatus(model.IndexedStatusSkipped, "", &reason, nil)
	}

	content, err := reader.Open(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("failed to open content: %v", err)
		return updateStatus(model.IndexedStatusError, "", nil, &errMessage)
	}

	if err := writer.Index(ctx, item.ID, reader.Kind(), content, reader.Attributes()); err != nil {
		errMessage := fmt.Sprintf("failed to index item %s: %v", item.ID, err)
		return updateStatus(model.IndexedStatusError, "", nil, &errMessage)
	}

	return updateStatus(model.IndexedStatusIndexed, knowledgeItem.Hash, nil, nil)
}
