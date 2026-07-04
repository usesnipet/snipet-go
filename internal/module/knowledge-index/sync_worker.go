package knowledgeindex

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/repository"
)

type SyncIndexArgs struct {
	IndexID string `json:"index_id"`
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

	indexRepo                repository.IKnowledgeIndexRepository
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository
	logger                   *logger.Logger
}

func NewSyncIndexWorker(
	indexRepo repository.IKnowledgeIndexRepository,
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository,
	logger *logger.Logger,
) *SyncIndexWorker {
	return &SyncIndexWorker{
		indexRepo:                indexRepo,
		indexedKnowledgeItemRepo: indexedKnowledgeItemRepo,
		logger:                   logger,
	}
}

func (s *SyncIndexWorker) Work(ctx context.Context, job *river.Job[SyncIndexArgs]) error {
	return nil
}
