package appuser

import (
	"context"

	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	userRepo repository.IAppUserRepository
	appRepo  repository.IAppRepository
}

func NewService(userRepo repository.IAppUserRepository, appRepo repository.IAppRepository) *Service {
	return &Service{userRepo: userRepo, appRepo: appRepo}
}

func (s *Service) ensureAppKeyAccess(ctx context.Context, appCode string) error {
	appIdentity, err := auth.CurrentAppKey(ctx)
	if err != nil {
		return err
	}
	if err = appIdentity.Is(appCode); err != nil {
		return err
	}
	return nil
}

func (s *Service) Create(ctx context.Context, appCode string, dto CreateAppUserDTO) error {
	if err := s.ensureAppKeyAccess(ctx, appCode); err != nil {
		return err
	}

	user := &model.AppUser{
		Name:     dto.Name,
		Metadata: dto.Metadata,
		Email:    &dto.Email,
		Picture:  dto.Picture,
	}
	if err := s.userRepo.CreateInApp(ctx, appCode, user, &dto.ExternalID); err != nil {
		return err
	}
	return nil
}

func (s *Service) FilterInApp(ctx context.Context, appCode string, filter *filter.Options[model.AppUser]) (*page.Paginated[model.AppUser], error) {
	if err := s.ensureAppKeyAccess(ctx, appCode); err != nil {
		return nil, err
	}

	return s.userRepo.FilterInApp(ctx, appCode, filter)
}

func (s *Service) Me(ctx context.Context, appCode string) (*model.AppUser, error) {
	identity, err := auth.CurrentAppUser(ctx)
	if err != nil {
		return nil, err
	}
	if err = identity.CanAccessApp(appCode); err != nil {
		return nil, err
	}
	return s.userRepo.FindByIDInApp(ctx, appCode, identity.UserID)
}
