package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/page"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	appRepo   repository.IAppRepository
	agentRepo repository.IAgentRepository
	txManager repository.ITxManager
	generator *auth.APIKeyGenerator
	hasher    *auth.KeyHasher
	logger    *logger.Logger
}

func NewService(
	appRepo repository.IAppRepository,
	agentRepo repository.IAgentRepository,
	txManager repository.ITxManager,
	generator *auth.APIKeyGenerator,
	hasher *auth.KeyHasher,
	logger *logger.Logger,
) *Service {
	return &Service{
		appRepo:   appRepo,
		agentRepo: agentRepo,
		txManager: txManager,
		generator: generator,
		hasher:    hasher,
		logger:    logger,
	}
}

func (s *Service) Filter(ctx context.Context, opts *filter.Options[model.App]) (*page.Paginated[model.App], error) {
	return s.appRepo.Filter(ctx, opts)
}

// FindByCode is the lookup used both internally (e.g. Ping) and by the
// admin CRUD surface.
func (s *Service) FindByCode(ctx context.Context, code string) (*model.App, error) {
	return s.appRepo.FindByCode(ctx, code)
}

// FindPublicByCode is the unauthenticated lookup for an app's public info
// (e.g. a frontend-only widget looking itself up by code). Only apps marked
// Public are exposed; a non-public app 404s the same as one that doesn't
// exist, so callers can't use this to probe for private app codes.
func (s *Service) FindPublicByCode(ctx context.Context, code string) (*PublicAppDTO, error) {
	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if !app.Public {
		return nil, apperr.NotFound("app is not public")
	}
	return &PublicAppDTO{
		Code:        app.Code,
		Name:        app.Name,
		Description: app.Description,
		AppToAgents: app.AppToAgents,
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

func (s *Service) Create(ctx context.Context, dto CreateAppDTO) (*AppWithSecret, error) {
	code, err := s.GenerateCode(ctx)
	if err != nil {
		return nil, err
	}

	plainKey, keyID, err := s.generator.Generate()
	if err != nil {
		return nil, err
	}
	keyHash, err := s.hasher.Hash(plainKey)
	if err != nil {
		return nil, err
	}

	if err := s.validateAgentIDs(ctx, dto.AgentIDs); err != nil {
		return nil, err
	}

	app := &model.App{
		Code:        code,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      model.AppStatusPending,
		Public:      dto.Public,
		KeyID:       keyID,
		KeyHash:     keyHash,
	}
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.appRepo.Create(ctx, app); err != nil {
			return err
		}
		if len(dto.AgentIDs) > 0 {
			return s.appRepo.ReplaceAgents(ctx, app.ID, dto.AgentIDs)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(dto.AgentIDs) > 0 {
		created, err := s.appRepo.FindByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		app = created
	}

	return &AppWithSecret{App: app, Key: plainKey}, nil
}

func (s *Service) Roll(ctx context.Context, code string) (*AppWithSecret, error) {
	plainKey, keyID, err := s.generator.Generate()
	if err != nil {
		return nil, err
	}
	keyHash, err := s.hasher.Hash(plainKey)
	if err != nil {
		return nil, err
	}

	updates := &model.App{KeyID: keyID, KeyHash: keyHash}
	if err := s.appRepo.UpdateByCode(ctx, code, updates); err != nil {
		return nil, err
	}

	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &AppWithSecret{App: app, Key: plainKey}, nil
}

func (s *Service) UpdateByCode(ctx context.Context, code string, dto UpdateAppDTO) error {
	if err := s.validateAgentIDs(ctx, dto.AgentIDs); err != nil {
		return err
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		updates := &model.App{}
		if dto.Name != nil {
			updates.Name = *dto.Name
		}
		if dto.Description != nil {
			updates.Description = *dto.Description
		}
		if err := s.appRepo.UpdateByCode(ctx, code, updates); err != nil {
			return err
		}

		if dto.Public != nil {
			if err := s.appRepo.UpdatePublicByCode(ctx, code, *dto.Public); err != nil {
				return err
			}
		}

		if dto.AgentIDs != nil {
			return s.linkAgents(ctx, code, dto.AgentIDs)
		}
		return nil
	})
}

// LinkAgents replaces the whole set of agents linked to the app identified by
// code. It backs PUT /apps/{code}/agents.
func (s *Service) LinkAgents(ctx context.Context, code string, agentIDs []string) error {
	if err := s.validateAgentIDs(ctx, agentIDs); err != nil {
		return err
	}
	return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		return s.linkAgents(ctx, code, agentIDs)
	})
}

// linkAgents resolves the app by code and swaps its agent links. Callers are
// responsible for validating agentIDs and opening a transaction.
func (s *Service) linkAgents(ctx context.Context, code string, agentIDs []string) error {
	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}
	return s.appRepo.ReplaceAgents(ctx, app.ID, agentIDs)
}

// validateAgentIDs rejects duplicates and any id that doesn't resolve to an
// existing agent, so a bad link fails with 400 rather than an FK error. A nil
// or empty slice is valid (it clears the links).
func (s *Service) validateAgentIDs(ctx context.Context, agentIDs []string) error {
	if len(agentIDs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(agentIDs))
	ids := make([]any, 0, len(agentIDs))
	for _, id := range agentIDs {
		if _, dup := seen[id]; dup {
			return apperr.BadRequest("duplicate agent id: " + id)
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	found, err := s.agentRepo.Filter(ctx, filter.New[model.Agent](filter.WhereIn("id", ids...)))
	if err != nil {
		return err
	}
	if found.Count() != len(agentIDs) {
		return apperr.BadRequest("one or more agent ids do not exist")
	}
	return nil
}

func (s *Service) UpdateAuthConfig(ctx context.Context, code string, dto UpdateAppAuthConfigDTO) error {
	return s.appRepo.UpdateAuthConfigByCode(ctx, code, dto.ToModel())
}

// ToggleActive flips an app between active and deactivated. A deactivated
// app fails key verification (see VerifyKey) so it can no longer ping in or
// authenticate.
func (s *Service) ToggleActive(ctx context.Context, code string, active bool) error {
	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}
	if app.Status == model.AppStatusPending && active && !app.Public {
		return apperr.BadRequest("cannot activate a pending app")
	}
	status := model.AppStatusActive
	if !active {
		status = model.AppStatusDeactivated
	}
	return s.appRepo.UpdateByCode(ctx, code, &model.App{Status: status})
}

func (s *Service) DeleteByCode(ctx context.Context, code string) error {
	return s.appRepo.DeleteByCode(ctx, code)
}

// VerifyKey is the auth entrypoint itself — looked up by the key's own
// key_id, called by guard.RequireAppKey. A deactivated app is rejected
// here so it can't authenticate at all.
func (s *Service) VerifyKey(ctx context.Context, key string) (*model.App, error) {
	keyID := s.generator.GetKeyID(key)
	paginated, err := s.appRepo.Filter(ctx, filter.New[model.App](filter.WhereEq("key_id", keyID)))
	if err != nil {
		return nil, apperr.Unauthorized("invalid app key")
	}
	if paginated.IsEmpty() {
		return nil, apperr.NotFound("app not found")
	}

	app := paginated.First()
	if app.Status == model.AppStatusDeactivated {
		return nil, apperr.Forbidden("app is disabled")
	}

	valid, err := s.hasher.Verify(key, app.KeyHash)
	if err != nil || !valid {
		return nil, apperr.Unauthorized("invalid app key")
	}
	return app, nil
}

// Ping is how a connected app confirms it's alive: it authenticates with its
// own key (guard.RequireAppKey), then hits this by its code. A pending app
// is promoted to active on its first successful ping.
func (s *Service) Ping(ctx context.Context, code string) error {
	identity, err := auth.CurrentAppKey(ctx)
	if err != nil {
		return err
	}
	if identity.Code != code {
		return apperr.Forbidden("app key does not match code")
	}

	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := &model.App{LastVerifiedAt: &now}
	if app.Status == model.AppStatusPending {
		updates.Status = model.AppStatusActive
	}
	return s.appRepo.UpdateByCode(ctx, code, updates)
}
