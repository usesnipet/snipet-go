package client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	clientRepo repository.IClientRepository
	logger     *logger.Logger
}

func NewService(clientRepo repository.IClientRepository, logger *logger.Logger) *Service {
	return &Service{clientRepo: clientRepo, logger: logger}
}

func (s *Service) Init(ctx context.Context, cfg *config.AppConfig) error {
	if !cfg.InheritClient {
		return nil
	}

	client, err := s.FindByCode(ctx, cfg.InheritClientCode)
	var notFoundError *apperr.Error
	if errors.As(err, &notFoundError) && notFoundError.StatusCode == http.StatusNotFound {
		s.logger.Infof("creating inherit client: %s with name %s", cfg.InheritClientCode, cfg.InheritClientName)
		var clientConfig model.ClientConfig
		clientConfig.Anonymous.Enabled = true
		_, err = s.CreateWithCode(ctx, CreateClientDTO{
			Name:   cfg.InheritClientName,
			Config: clientConfig,
		}, cfg.InheritClientCode)
		return err
	}
	if err != nil {
		return err
	}

	if client.Name != cfg.InheritClientName {
		s.logger.Infof("inherit client name update: %s -> %s", client.Name, cfg.InheritClientName)
		return s.UpdateByCode(
			ctx,
			cfg.InheritClientCode,
			UpdateClientDTO{Name: &cfg.InheritClientName},
		)
	}
	return nil
}

func (s *Service) Filter(ctx context.Context, filter *filter.Options[model.Client]) (*page.Paginated[model.Client], error) {
	return s.clientRepo.Filter(ctx, filter)
}

func (s *Service) FindByCode(ctx context.Context, code string) (*model.Client, error) {
	paginated, err := s.clientRepo.Filter(ctx, filter.New[model.Client](filter.WhereEq("code", code)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("client not found")
	}
	return paginated.First(), nil
}

func (s *Service) FindPublicByCode(ctx context.Context, code string) (*ClientPublicDTO, error) {
	client, err := s.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &ClientPublicDTO{
		Code:           client.Code,
		Name:           client.Name,
		AllowAnonymous: client.Config.Anonymous.Enabled,
	}, nil
}

func (s *Service) GenerateCode(ctx context.Context) (string, error) {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}
	code := base64.RawURLEncoding.EncodeToString(randomBytes)[:10] // first 10 characters

	// check if code is already in use
	_, err := s.FindByCode(ctx, code)
	var notFoundError *apperr.Error
	if errors.As(err, &notFoundError) && notFoundError.StatusCode == http.StatusNotFound {
		return code, nil
	}
	if err != nil {
		return "", err
	}
	return s.GenerateCode(ctx)
}

func (s *Service) Create(ctx context.Context, dto CreateClientDTO) (*model.Client, error) {
	code, err := s.GenerateCode(ctx)
	if err != nil {
		return nil, err
	}
	return s.CreateWithCode(ctx, dto, code)
}

func (s *Service) CreateWithCode(ctx context.Context, dto CreateClientDTO, code string) (*model.Client, error) {
	if dto.Config.Webhook.URL != "" {
		dto.Config.Webhook.Enabled = true
	}
	if dto.Config.OIDC.Issuer != "" || dto.Config.OIDC.Audience != "" {
		dto.Config.OIDC.Enabled = true
	}

	client := &model.Client{
		Code:   code,
		Name:   dto.Name,
		Config: dto.Config,
	}
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) UpdateByCode(ctx context.Context, code string, dto UpdateClientDTO) error {
	updates := &model.Client{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Config != nil {
		updates.Config = *dto.Config
	}
	return s.clientRepo.UpdateByCode(ctx, code, updates)
}

func (s *Service) DeleteByCode(ctx context.Context, code string) error {
	return s.clientRepo.DeleteByCode(ctx, code)
}
