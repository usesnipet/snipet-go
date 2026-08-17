package llm

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/authz"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

type Service struct {
	repo       repository.ILLMRepository
	llmManager *manager.Driver[llm.Driver]
}

func NewService(repo repository.ILLMRepository, llmManager *manager.Driver[llm.Driver]) *Service {
	return &Service{
		repo:       repo,
		llmManager: llmManager,
	}
}

func (s *Service) Filter(ctx context.Context, tenantID string, opts *filter.Options[model.LLM]) (*page.Paginated[model.LLM], error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	return s.repo.Filter(ctx, filter.Merge(opts, filter.New[model.LLM](filter.WhereEq("tenant_id", tenantID))))
}

// findInTenant fetches by id then verifies the row belongs to tenantID —
// the generic repo FindByID is id-only, so cross-tenant existence is
// checked here (404, not 403, so it doesn't leak whether the id exists at
// all in another tenant).
func (s *Service) findInTenant(ctx context.Context, tenantID, id string) (*model.LLM, error) {
	found, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if found.TenantID != tenantID {
		return nil, apperr.NotFound("llm not found")
	}
	return found, nil
}

func (s *Service) FindByID(ctx context.Context, tenantID, id string) (*model.LLM, error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	return s.findInTenant(ctx, tenantID, id)
}

func (s *Service) Create(ctx context.Context, tenantID string, dto CreateLLMDTO) (*model.LLM, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return nil, err
	}

	if err := s.llmManager.ValidateConfigurationByKey(dto.Provider, dto.Configuration); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	llm := &model.LLM{
		TenantID:      tenantID,
		Name:          dto.Name,
		Provider:      dto.Provider,
		Configuration: dto.Configuration,
	}
	if err := s.repo.Create(ctx, llm); err != nil {
		return nil, err
	}
	return llm, nil
}

func (s *Service) Update(ctx context.Context, tenantID, id string, dto UpdateLLMDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}

	existing, err := s.findInTenant(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if dto.Provider != nil || dto.Configuration != nil {
		provider := existing.Provider
		configuration := existing.Configuration
		if dto.Provider != nil {
			provider = *dto.Provider
		}
		if dto.Configuration != nil {
			configuration = dto.Configuration
		}

		if err := s.llmManager.ValidateConfigurationByKey(provider, configuration); err != nil {
			return apperr.BadRequest(err.Error())
		}
	}

	updates := &model.LLM{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Provider != nil {
		updates.Provider = *dto.Provider
	}
	if dto.Configuration != nil {
		updates.Configuration = dto.Configuration
	}
	return s.repo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, tenantID, id string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleUser); err != nil {
		return err
	}
	if _, err := s.findInTenant(ctx, tenantID, id); err != nil {
		return err
	}
	return s.repo.DeleteByID(ctx, id)
}

func (s *Service) ListDrivers(ctx context.Context, tenantID string) ([]driver.Info, error) {
	if _, err := authz.RequireMember(ctx, tenantID); err != nil {
		return nil, err
	}
	return s.llmManager.ListDrivers(ctx)
}
