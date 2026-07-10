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

	sourceManager            *runtime.SourceManager
	indexManager             *runtime.IndexManager
	knowledgeRepo            repository.IKnowledgeRepository
	indexRepo                repository.IKnowledgeIndexRepository
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository
	logger                   *logger.Logger
}

func NewSyncIndexWorker(
	indexManager *runtime.IndexManager,
	sourceManager *runtime.SourceManager,
	knowledgeRepo repository.IKnowledgeRepository,
	indexRepo repository.IKnowledgeIndexRepository,
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository,
	logger *logger.Logger,
) *SyncIndexWorker {
	return &SyncIndexWorker{
		indexManager:             indexManager,
		sourceManager:            sourceManager,
		knowledgeRepo:            knowledgeRepo,
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
	if err := s.indexedKnowledgeItemRepo.CreateManyInIndex(ctx, job.Args.KnowledgeID, job.Args.IndexID, indexedToCreate); err != nil {
		return err
	}

	toUpdateIDs := util.Map(toUpdate, func(item model.IndexedKnowledgeItem) string {
		return item.ID
	})
	toDeleteIDs := util.Map(toDelete, func(item model.IndexedKnowledgeItem) string {
		return item.ID
	})
	err = s.indexedKnowledgeItemRepo.UpdateStatusesByIDsInIndex(
		ctx,
		job.Args.KnowledgeID,
		job.Args.IndexID,
		append(toUpdateIDs, toDeleteIDs...),
		model.IndexedStatusPending,
	)
	if err != nil {
		return err
	}

	err = writer.DeleteMany(ctx, toDeleteIDs)
	if err != nil {
		return err
	}

	for _, item := range indexedToCreate {
		s.indexItem(
			ctx,
			job.Args.KnowledgeID,
			job.Args.IndexID,
			writer,
			sourceDriver,
			*item,
			false,
		)
	}

	for _, item := range toUpdate {
		s.indexItem(
			ctx,
			job.Args.KnowledgeID,
			job.Args.IndexID,
			writer,
			sourceDriver,
			item,
			true,
		)
	}

	return nil
}

func (s *SyncIndexWorker) indexItem(
	ctx context.Context,
	knowledgeID string,
	indexID string,
	writer runtime.IIndexWriter,
	sourceDriver runtime.ISourceDriver,
	item model.IndexedKnowledgeItem,
	deleteBeforeIndex bool,
) error {
	updateStatus := func(status model.IndexStatus, reason *string, errMessage *string) error {
		return s.indexedKnowledgeItemRepo.UpdateInIndex(
			ctx,
			knowledgeID,
			indexID,
			item.ID,
			&model.IndexedKnowledgeItem{
				Status:    status,
				Reason:    reason,
				LastError: errMessage,
			},
		)
	}

	err := updateStatus(model.IndexedStatusSyncing, nil, nil)
	if err != nil {
		return err
	}

	if deleteBeforeIndex {
		err := writer.DeleteMany(ctx, []string{item.ID})
		if err != nil {
			errMessage := fmt.Sprintf("failed to delete item before indexing %s: %v", item.ID, err)
			return updateStatus(model.IndexedStatusError, nil, &errMessage)
		}
	}

	content, err := sourceDriver.GetContent(ctx, *item.KnowledgeItemID)
	if err != nil {
		errMessage := fmt.Sprintf("failed to get content: %v", err)
		return updateStatus(model.IndexedStatusError, nil, &errMessage)
	}

	if !slices.Contains(writer.SupportedKinds(), content.Kind()) {
		reason := fmt.Sprintf("item kind (%s) is not supported by index", content.Kind())
		return updateStatus(model.IndexedStatusSkipped, &reason, nil)
	}

	err = writer.Index(ctx, content)
	if err != nil {
		errMessage := fmt.Sprintf("failed to index item %s: %v", item.ID, err)
		return updateStatus(model.IndexedStatusError, nil, &errMessage)
	}

	return updateStatus(model.IndexedStatusIndexed, nil, nil)
}
