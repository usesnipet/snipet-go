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
	userRepo repository.IUserRepository
}

func NewService(userRepo repository.IUserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) generateAnonymousName() string {
	return fmt.Sprintf("Anonymous %s", uuid.New().String()[:8])
}

func (s *Service) CreateAnonymous(ctx context.Context, clientCode string, dto CreateAnonymousClientUserDTO) error {
	name := s.generateAnonymousName()
	if dto.Name != nil {
		name = *dto.Name
	}
	user := &model.User{
		Name:     name,
		Metadata: dto.Metadata,
	}
	if err := s.userRepo.CreateInClient(ctx, clientCode, user, nil); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateAuthenticated(ctx context.Context, clientCode string, dto CreateAuthenticatedClientUserDTO) error {
	user := &model.User{
		Name:     dto.Name,
		Metadata: dto.Metadata,
		Email:    &dto.Email,
		Picture:  dto.Picture,
	}
	if err := s.userRepo.CreateInClient(ctx, clientCode, user, &dto.ExternalID); err != nil {
		return err
	}
	return nil
}

func (s *Service) FilterInClient(ctx context.Context, clientCode string, filter *filter.Options[model.User]) (*page.Paginated[model.User], error) {
	return s.userRepo.FilterInClient(ctx, clientCode, filter)
}
