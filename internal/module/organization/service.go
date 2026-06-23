package organization

import (
	"context"

	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/infra/database"
	"github.com/usesnipet/snipet/internal/model"
)

type Service struct {
	repository IRepository
}

func (s *Service) FindBy(ctx context.Context) (*database.Paginated[model.Organization], error) {
	return s.repository.FindBy(ctx, filter.New[model.Organization]())
}

func (s *Service) Create(ctx context.Context, dto CreateOrganizationDTO) error {
	return s.repository.Create(ctx, &model.Organization{
		Name: dto.Name,
		Slug: dto.Slug,
	})
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
}
