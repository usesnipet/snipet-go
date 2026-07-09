package knowledge

import (
	"context"
	"errors"
	"log"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/queue"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

type Service struct {
	txManager         repository.ITxManager
	repo              repository.IKnowledgeRepository
	knowledgeItemRepo repository.IKnowledgeItemRepository
	sourceManager     *runtime.SourceManager
	riverClient       queue.IJobQueue
}

func NewService(
	txManager repository.ITxManager,
	repo repository.IKnowledgeRepository,
	knowledgeItemRepo repository.IKnowledgeItemRepository,
	sourceManager *runtime.SourceManager,
	riverClient queue.IJobQueue,
) *Service {
	return &Service{
		txManager:         txManager,
		repo:              repo,
		knowledgeItemRepo: knowledgeItemRepo,
		sourceManager:     sourceManager,
		riverClient:       riverClient,
	}
}

func (s *Service) Filter(ctx context.Context, filter *filter.Options[model.Knowledge]) (*page.Paginated[model.Knowledge], error) {
	return s.repo.Filter(ctx, filter)
}

func (s *Service) FilterItems(
	ctx context.Context,
	knowledgeID string,
	filter *filter.Options[model.KnowledgeItem],
) (*page.Paginated[model.KnowledgeItem], error) {
	return s.knowledgeItemRepo.FilterInKnowledge(ctx, knowledgeID, filter)
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Knowledge, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateKnowledgeDTO) (*model.Knowledge, int64, error) {
	if err := s.TestConnection(ctx, dto.Driver, dto.Configuration); err != nil {
		return nil, 0, apperr.BadRequest(err.Error())
	}

	knowledge := &model.Knowledge{
		Name:          dto.Name,
		Description:   dto.Description,
		Driver:        dto.Driver,
		Configuration: dto.Configuration,
	}

	var jobID int64
	err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.repo.Create(ctx, knowledge); err != nil {
			return err
		}
		jID, err := s.Sync(ctx, knowledge.ID, false)
		jobID = jID
		return err
	})

	return knowledge, jobID, err
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

func (s *Service) TestConnection(ctx context.Context, driver string, config util.JSONMap) error {
	_, err := s.sourceManager.Prepare(ctx, driver, config)
	if err != nil {
		if errors.Is(err, runtime.ErrDriverNotFound) {
			return apperr.NotFound(err.Error())
		}
		return apperr.BadRequest(err.Error())
	}
	return nil
}

func (s *Service) Sync(ctx context.Context, knowledgeID string, force bool) (int64, error) {
	_, err := s.FindByID(ctx, knowledgeID)
	if err != nil {
		return 0, err
	}
	log.Printf("syncing knowledge %s with force: %v", knowledgeID, force)
	jobID, err := s.riverClient.Push(ctx, SyncKnowledgeArgs{KnowledgeID: knowledgeID, Force: force}, nil)

	if err != nil {
		return 0, apperr.InternalServerError(err.Error())
	}
	return jobID, nil
}

func (s *Service) GetSyncStatus(ctx context.Context, jobID int64) (string, error) {
	job, err := s.riverClient.JobGet(ctx, jobID)
	if err != nil {
		return "", apperr.InternalServerError(err.Error())
	}
	return string(job.State), nil
}
