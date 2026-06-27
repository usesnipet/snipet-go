package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	cUserRepo  repository.IUserRepository
	clientRepo repository.IClientRepository
}

func NewService(cUserRepo repository.IUserRepository, clientRepo repository.IClientRepository) *Service {
	return &Service{cUserRepo: cUserRepo, clientRepo: clientRepo}
}

func (s *Service) generateAnonymousName() string {
	return fmt.Sprintf("Anonymous %s", uuid.New().String()[:8])
}

func (s *Service) CreateAnonymous(ctx context.Context, clientCode string, dto CreateAnonymousClientUserDTO) error {
	name := s.generateAnonymousName()
	if dto.Name != nil {
		name = *dto.Name
	}
	cUser := &model.User{
		Name:     name,
		Metadata: dto.Metadata,
	}
	if err := s.cUserRepo.CreateInClient(ctx, clientCode, cUser, nil); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateAuthenticated(ctx context.Context, clientCode string, dto CreateAuthenticatedClientUserDTO) error {
	cUser := &model.User{
		Name:     dto.Name,
		Metadata: dto.Metadata,
	}
	if err := s.cUserRepo.CreateInClient(ctx, clientCode, cUser, &dto.ExternalID); err != nil {
		return err
	}
	return nil
}

func (s *Service) FilterInClient(ctx context.Context, clientCode string) (*page.Paginated[model.User], error) {
	return s.cUserRepo.FilterInClient(ctx, clientCode, filter.Default[model.User]())
}
