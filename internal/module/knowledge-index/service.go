package knowledgeindex

import (
	"context"
	"errors"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	kdriver "github.com/usesnipet/snipet/pkg/driver/knowledge"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type Service struct {
	repo                     repository.IKnowledgeIndexRepository
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository
	indexManager             *manager.Driver[kdriver.IIndexDriver]
	pool                     queue.IPool
	syncWorker               *SyncIndexWorker
	txManager                repository.ITxManager
}

func NewService(
	repo repository.IKnowledgeIndexRepository,
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository,
	indexManager *manager.Driver[kdriver.IIndexDriver],
	pool queue.IPool,
	syncWorker *SyncIndexWorker,
	txManager repository.ITxManager,
) *Service {
	return &Service{
		repo:                     repo,
		indexedKnowledgeItemRepo: indexedKnowledgeItemRepo,
		indexManager:             indexManager,
		pool:                     pool,
		syncWorker:               syncWorker,
		txManager:                txManager,
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
	if err := s.TestConnection(ctx, dto.Driver, dto.Configuration); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	index := &model.KnowledgeIndex{
		Name:          dto.Name,
		Driver:        dto.Driver,
		Configuration: dto.Configuration,
		KnowledgeID:   knowledgeID,
	}
	err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		return s.repo.CreateInKnowledge(ctx, knowledgeID, index)
	})
	return index, err
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

func (s *Service) ListDrivers(ctx context.Context) (*DriversDTO, error) {
	indexDrivers, err := s.indexManager.ListDrivers(ctx)
	if err != nil {
		return nil, apperr.InternalServerError(err.Error())
	}

	return &DriversDTO{
		IndexDrivers: indexDrivers,
	}, nil
}

func (s *Service) TestConnection(ctx context.Context, key string, config jsonx.JSONMap) error {
	_, err := s.indexManager.Prepare(ctx, key, config)
	if err != nil {
		if errors.Is(err, driver.ErrDriverNotFound) {
			return apperr.NotFound(err.Error())
		}
		return apperr.BadRequest(err.Error())
	}
	return nil
}

func (s *Service) Sync(ctx context.Context, knowledgeID, indexID string) error {
	_, err := s.FindByID(ctx, knowledgeID, indexID)
	if err != nil {
		return err
	}

	err = s.pool.Submit(ctx, func(ctx context.Context) error {
		return s.syncWorker.Sync(ctx, knowledgeID, indexID)
	})
	if err != nil {
		return apperr.InternalServerError(err.Error())
	}
	return nil
}
