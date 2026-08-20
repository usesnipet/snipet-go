package client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/authz"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	clientRepo repository.IClientRepository
	agentRepo  repository.IAgentRepository
	logger     *logger.Logger
}

func NewService(clientRepo repository.IClientRepository, agentRepo repository.IAgentRepository, logger *logger.Logger) *Service {
	return &Service{clientRepo: clientRepo, agentRepo: agentRepo, logger: logger}
}

func (s *Service) Filter(ctx context.Context, tenantID string, opts *filter.Options[model.App]) (*page.Paginated[model.App], error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	return s.clientRepo.Filter(ctx, filter.Merge(opts, filter.New[model.App](filter.WhereEq("tenant_id", tenantID))))
}

// FindByCode is the tenant-agnostic lookup used by the public/widget-facing
// surface (FindPublicByCode, GetAgents) and by Init's bootstrap flow, where
// there's no tenant-staff caller to scope against.
func (s *Service) FindByCode(ctx context.Context, code string) (*model.App, error) {
	paginated, err := s.clientRepo.Filter(ctx, filter.New[model.App](filter.WhereEq("code", code)))
	if err != nil {
		return nil, err
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("client not found")
	}
	return paginated.First(), nil
}

// findByCodeInTenant is the admin-path lookup — verifies the client belongs
// to tenantID (404, not 403, to avoid confirming the code exists elsewhere).
func (s *Service) findByCodeInTenant(ctx context.Context, tenantID, code string) (*model.App, error) {
	found, err := s.clientRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if found.TenantID != tenantID {
		return nil, apperr.NotFound("client not found")
	}
	return found, nil
}

func (s *Service) FindByCodeInTenant(ctx context.Context, tenantID, code string) (*model.App, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	return s.findByCodeInTenant(ctx, tenantID, code)
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

func (s *Service) Create(ctx context.Context, tenantID string, dto CreateClientDTO) (*model.App, error) {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return nil, err
	}
	code, err := s.GenerateCode(ctx)
	if err != nil {
		return nil, err
	}
	return s.createWithCode(ctx, tenantID, dto, code)
}

func (s *Service) createWithCode(ctx context.Context, tenantID string, dto CreateClientDTO, code string) (*model.App, error) {
	if dto.Config.Webhook.URL != "" {
		dto.Config.Webhook.Enabled = true
	}
	if dto.Config.OIDC.Issuer != "" || dto.Config.OIDC.Audience != "" {
		dto.Config.OIDC.Enabled = true
	}

	client := &model.App{
		TenantID: tenantID,
		Code:     code,
		Name:     dto.Name,
		Config:   dto.Config,
	}
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Service) UpdateByCode(ctx context.Context, tenantID, code string, dto UpdateClientDTO) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	if _, err := s.findByCodeInTenant(ctx, tenantID, code); err != nil {
		return err
	}

	updates := &model.App{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Config != nil {
		updates.Config = *dto.Config
	}
	return s.clientRepo.UpdateByCode(ctx, code, updates)
}

func (s *Service) DeleteByCode(ctx context.Context, tenantID, code string) error {
	if _, err := authz.RequireTenantRole(ctx, tenantID, model.RoleAdmin); err != nil {
		return err
	}
	if _, err := s.findByCodeInTenant(ctx, tenantID, code); err != nil {
		return err
	}
	return s.clientRepo.DeleteByCode(ctx, code)
}

// GetAgents is reached via anyAuthMiddleware (API key or client-widget-user
// JWT), not tenant-staff bearer auth — no tenantID param, per decision 3.
// Still scopes results to the client's own tenant (fixes a pre-existing bug
// where this returned every agent in the system regardless of client) —
// this is not "agents linked to this client" (no such link table exists
// yet), just the minimum fix to stop leaking other tenants' agents.
func (s *Service) GetAgents(
	ctx context.Context,
	clientCode string,
	opts *filter.Options[model.Agent],
) (*page.Paginated[model.Agent], error) {
	client, err := s.FindByCode(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	return s.agentRepo.Filter(ctx, filter.Merge(opts, filter.New[model.Agent](filter.WhereEq("tenant_id", client.TenantID))))
}
