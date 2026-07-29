package llm

import (
	"context"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/driver/llm"
)

type Service struct {
	repo       repository.ILLMRepository
	llmManager *runtime.DriverManager[llm.Driver]
}

func NewService(repo repository.ILLMRepository, llmManager *runtime.DriverManager[llm.Driver]) *Service {
	return &Service{
		repo:       repo,
		llmManager: llmManager,
	}
}

func (s *Service) Filter(ctx context.Context) (*page.Paginated[model.LLM], error) {
	return s.repo.Filter(ctx, filter.Default[model.LLM]())
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.LLM, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateLLMDTO) (*model.LLM, error) {
	if err := s.llmManager.ValidateConfigurationByKey(dto.Provider, dto.Configuration); err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	llm := &model.LLM{
		Name:          dto.Name,
		Provider:      dto.Provider,
		Configuration: dto.Configuration,
	}
	if err := s.repo.Create(ctx, llm); err != nil {
		return nil, err
	}
	return llm, nil
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateLLMDTO) error {
	if dto.Provider != nil || dto.Configuration != nil {
		existing, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}

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

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *Service) ListDrivers(ctx context.Context) ([]driver.Info, error) {
	return s.llmManager.ListDrivers(ctx)
}
