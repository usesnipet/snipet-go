package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	userRepo repository.IUserRepository
	config   config.UserConfig
}

func NewService(userRepo repository.IUserRepository, config config.UserConfig) *Service {
	return &Service{userRepo: userRepo, config: config}
}

// Init creates the bootstrap admin user if no User exists yet. Called once
// from bootstrap, before any tenant is created.
func (s *Service) Init(ctx context.Context) error {
	existing, err := s.userRepo.Filter(ctx, filter.Default[model.User]())
	if err != nil {
		return err
	}
	if existing.IsNotEmpty() {
		return nil
	}

	password := s.config.AdminPassword
	generated := password == ""
	if generated {
		token, err := auth.NewTokenService().GenerateToken()
		if err != nil {
			return err
		}
		password = token
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	admin := &model.User{
		Name:         s.config.AdminName,
		Email:        s.config.AdminEmail,
		PasswordHash: &passwordHash,
		IsAdmin:      true,
		Challenges:   []model.Challenge{},
	}
	if err := s.userRepo.Create(ctx, admin); err != nil {
		return err
	}

	if generated {
		fmt.Println("generated admin password (shown once):", password)
	}
	return nil
}

func (s *Service) Filter(ctx context.Context, opts *filter.Options[model.User]) (*page.Paginated[model.User], error) {
	return s.userRepo.Filter(ctx, opts)
}

func (s *Service) FindByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *Service) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *Service) Create(ctx context.Context, dto CreateUserDTO) (*model.User, error) {
	if _, err := s.userRepo.FindByEmail(ctx, dto.Email); err == nil {
		return nil, apperr.Conflict("email already in use")
	} else {
		var appErr *apperr.Error
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	passwordHash, err := auth.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	created := &model.User{
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: &passwordHash,
		Picture:      dto.Picture,
		IsAdmin:      dto.IsAdmin,
		Challenges:   []model.Challenge{},
	}
	if err := s.userRepo.Create(ctx, created); err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) UpdateByID(ctx context.Context, id string, dto UpdateUserDTO) error {
	updates := &model.User{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Email != nil {
		updates.Email = *dto.Email
	}
	if dto.Picture != nil {
		updates.Picture = dto.Picture
	}
	if dto.IsAdmin != nil {
		updates.IsAdmin = *dto.IsAdmin
	}
	return s.userRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.userRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.userRepo.DeleteByID(ctx, id)
}

func (s *Service) Me(ctx context.Context) (*model.User, error) {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return identity.User, nil
}

func (s *Service) UpdateMyPicture(ctx context.Context, dto UpdateProfilePictureDTO) error {
	identity, err := auth.CurrentUser(ctx)
	if err != nil {
		return err
	}
	return s.userRepo.UpdateByID(ctx, identity.User.ID, &model.User{Picture: &dto.Picture})
}
