package knowledgeindex

import (
	"context"
	"errors"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/authz"
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
	knowledgeRepo            repository.IKnowledgeRepository
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository
	indexManager             *manager.Driver[kdriver.IIndexDriver]
	pool                     queue.IPool
	syncWorker               *SyncIndexWorker
	txManager                repository.ITxManager
}

func NewService(
	repo repository.IKnowledgeIndexRepository,
	knowledgeRepo repository.IKnowledgeRepository,
	indexedKnowledgeItemRepo repository.IIndexedKnowledgeItemRepository,
	indexManager *manager.Driver[kdriver.IIndexDriver],
	pool queue.IPool,
	syncWorker *SyncIndexWorker,
	txManager repository.ITxManager,
) *Service {
	return &Service{
		repo:                     repo,
		knowledgeRepo:            knowledgeRepo,
		indexedKnowledgeItemRepo: indexedKnowledgeItemRepo,
		indexManager:             indexManager,
		pool:                     pool,
		syncWorker:               syncWorker,
		txManager:                txManager,
	}
}

// knowledgeInTenant verifies knowledgeID belongs to tenantID before any
// operation on its indexes delegates to the knowledge_id-scoped repo
// methods below — those don't know about tenants at all, so the boundary
// is enforced here, one hop up.
func (s *Service) knowledgeInTenant(ctx context.Context, tenantID, knowledgeID string) (*model.Knowledge, error) {
	knowledge, err := s.knowledgeRepo.FindByID(ctx, knowledgeID)
	if err != nil {
		return nil, err
	}
	if knowledge.TenantID != tenantID {
		return nil, apperr.NotFound("knowledge not found")
	}
	return knowledge, nil
}

func (s *Service) Filter(ctx context.Context, tenantID, knowledgeID string, opts *filter.Options[model.KnowledgeIndex]) (*page.Paginated[model.KnowledgeIndex], error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	if _, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID); err != nil {
		return nil, err
	}
	return s.repo.FilterInKnowledge(ctx, knowledgeID, opts)
}

func (s *Service) FilterItems(ctx context.Context, tenantID, knowledgeID, indexID string, opts *filter.Options[model.IndexedKnowledgeItem]) (*page.Paginated[model.IndexedKnowledgeItem], error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	if _, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID); err != nil {
		return nil, err
	}
	return s.indexedKnowledgeItemRepo.FilterInIndex(ctx, knowledgeID, indexID, opts)
}

func (s *Service) FindByID(ctx context.Context, tenantID, knowledgeID, id string) (*model.KnowledgeIndex, error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	if _, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID); err != nil {
		return nil, err
	}
	return s.repo.FindByIDInKnowledge(ctx, knowledgeID, id)
}

func (s *Service) Create(ctx context.Context, tenantID, knowledgeID string, dto CreateKnowledgeIndexDTO) (*model.KnowledgeIndex, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return nil, err
	}
	knowledge, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID)
	if err != nil {
		return nil, err
	}

	if err := s.TestConnection(ctx, tenantID, dto.Driver, dto.Configuration); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	index := &model.KnowledgeIndex{
		TenantID:      knowledge.TenantID,
		Name:          dto.Name,
		Driver:        dto.Driver,
		Configuration: dto.Configuration,
		KnowledgeID:   knowledgeID,
	}
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		return s.repo.CreateInKnowledge(ctx, knowledgeID, index)
	})
	return index, err
}

func (s *Service) Update(ctx context.Context, tenantID, knowledgeID, id string, dto UpdateKnowledgeIndexDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}
	if _, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID); err != nil {
		return err
	}
	if dto.Name != nil {
		return s.repo.UpdateNameInKnowledge(ctx, knowledgeID, id, *dto.Name)
	}
	return nil
}

func (s *Service) DeleteByID(ctx context.Context, tenantID, knowledgeID, id string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}
	if _, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID); err != nil {
		return err
	}
	return s.repo.DeleteInKnowledge(ctx, knowledgeID, id)
}

func (s *Service) ListDrivers(ctx context.Context, tenantID string) (*DriversDTO, error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}

	indexDrivers, err := s.indexManager.ListDrivers(ctx)
	if err != nil {
		return nil, apperr.InternalServerError(err.Error())
	}

	return &DriversDTO{
		IndexDrivers: indexDrivers,
	}, nil
}

func (s *Service) TestConnection(ctx context.Context, tenantID string, key string, config jsonx.JSONMap) error {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return err
	}

	_, err := s.indexManager.Prepare(ctx, key, config)
	if err != nil {
		if errors.Is(err, driver.ErrDriverNotFound) {
			return apperr.NotFound(err.Error())
		}
		return apperr.BadRequest(err.Error())
	}
	return nil
}

func (s *Service) Sync(ctx context.Context, tenantID, knowledgeID, indexID string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}
	if _, err := s.knowledgeInTenant(ctx, tenantID, knowledgeID); err != nil {
		return err
	}
	if _, err := s.repo.FindByIDInKnowledge(ctx, knowledgeID, indexID); err != nil {
		return err
	}

	err := s.pool.Submit(ctx, func(ctx context.Context) error {
		return s.syncWorker.Sync(ctx, knowledgeID, indexID)
	})
	if err != nil {
		return apperr.InternalServerError(err.Error())
	}
	return nil
}
