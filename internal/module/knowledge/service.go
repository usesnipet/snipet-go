package knowledge

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
	txManager         repository.ITxManager
	repo              repository.IKnowledgeRepository
	knowledgeItemRepo repository.IKnowledgeItemRepository
	sourceManager     *manager.Driver[kdriver.ISourceDriver]
	pool              queue.IPool
	syncWorker        *SyncWorker
}

func NewService(
	txManager repository.ITxManager,
	repo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	sourceManager *manager.Driver[kdriver.ISourceDriver],
	pool queue.IPool,
	syncWorker *SyncWorker,
) *Service {
	return &Service{
		txManager:         txManager,
		repo:              repo,
		knowledgeItemRepo: knowledgeItemRepo,
		sourceManager:     sourceManager,
		pool:              pool,
		syncWorker:        syncWorker,
	}
}

func (s *Service) Filter(ctx context.Context, opts *filter.Options[model.Knowledge]) (*page.Paginated[model.Knowledge], error) {
	return s.repo.Filter(ctx, opts)
}

func (s *Service) FilterItems(
	ctx context.Context,
	knowledgeID string,
	opts *filter.Options[model.KnowledgeItem],
) (*page.Paginated[model.KnowledgeItem], error) {
	return s.knowledgeItemRepo.FilterInKnowledge(ctx, knowledgeID, opts)
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Knowledge, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateKnowledgeDTO) (*model.Knowledge, error) {
	if err := s.TestConnection(ctx, dto.Driver, dto.Configuration); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	knowledge := &model.Knowledge{
		Name:          dto.Name,
		Description:   dto.Description,
		Driver:        dto.Driver,
		Configuration: dto.Configuration,
		SyncStatus:    model.SyncStatusPending,
	}

	err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		return s.repo.Create(ctx, knowledge)
	})
	if err != nil {
		return nil, err
	}

	if err := s.Sync(ctx, knowledge.ID, true); err != nil {
		return knowledge, err
	}
	return knowledge, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateKnowledgeDTO) error {
	updates := &model.Knowledge{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Description != nil {
		updates.Description = *dto.Description
	}
	return s.repo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *Service) ListDrivers(ctx context.Context) (*DriversDTO, error) {
	sourceDrivers, err := s.sourceManager.ListDrivers(ctx)
	if err != nil {
		return nil, apperr.InternalServerError(err.Error())
	}

	return &DriversDTO{
		SourceDrivers: sourceDrivers,
	}, nil
}

func (s *Service) TestConnection(ctx context.Context, key string, config jsonx.JSONMap) error {
	_, err := s.sourceManager.Prepare(ctx, key, config)
	if err != nil {
		if errors.Is(err, driver.ErrDriverNotFound) {
			return apperr.NotFound(err.Error())
		}
		return apperr.BadRequest(err.Error())
	}
	return nil
}

func (s *Service) Sync(ctx context.Context, knowledgeID string, force bool) error {
	err := s.pool.Submit(ctx, func(ctx context.Context) error {
		return s.syncWorker.Sync(ctx, knowledgeID, force)
	})
	if err != nil {
		return apperr.InternalServerError(err.Error())
	}
	return nil
}
