package knowledge

import (
	"context"
	"errors"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

type Service struct {
	repository    repository.IKnowledgeRepository
	sourceManager *runtime.SourceManager
}

func NewService(repository repository.IKnowledgeRepository, sourceManager *runtime.SourceManager) *Service {
	return &Service{repository: repository, sourceManager: sourceManager}
}

func (s *Service) Filter(ctx context.Context, filter *filter.Options[model.Knowledge]) (*page.Paginated[model.Knowledge], error) {
	return s.repository.Filter(ctx, filter)
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.Knowledge, error) {
	return s.repository.FindByID(ctx, id)
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
	}

	if err := s.repository.Create(ctx, knowledge); err != nil {
		return nil, err
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
	return s.repository.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *Service) TestConnection(ctx context.Context, driver string, config util.JSONMap) error {
	_, err := s.sourceManager.Prepare(ctx, driver, config)
	if err != nil {
		if errors.Is(err, runtime.ErrSourceDriverNotFound) {
			return apperr.NotFound(err.Error())
		}
		return apperr.BadRequest(err.Error())
	}
	return nil
}
