package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/filter"
	"github.com/usesnipet/snipet/internal/license"
	"github.com/usesnipet/snipet/internal/model"
	"github.com/usesnipet/snipet/internal/repository"
)

type Service struct {
	tenantRepo repository.ITenantRepository
	memberRepo repository.IMemberRepository
	userRepo   repository.IUserRepository
	txManager  repository.ITxManager
	config     config.TenantConfig
	license    *license.Service
}

func NewService(
	tenantRepo repository.ITenantRepository,
	memberRepo repository.IMemberRepository,
	userRepo repository.IUserRepository,
	txManager repository.ITxManager,
	config config.TenantConfig,
	license *license.Service,
) *Service {
	return &Service{
		tenantRepo: tenantRepo,
		memberRepo: memberRepo,
		userRepo:   userRepo,
		txManager:  txManager,
		config:     config,
		license:    license,
	}
}

// Init creates the bootstrap Tenant if none exists yet, and makes sure the
// bootstrap admin user (user.Service.Init) is an admin Member of it — an
// unlicensed instance never reaches Create, so this is the only place that
// membership gets attached. Idempotent, always runs regardless of license
// state: every instance needs exactly this one Tenant to exist. Must run
// after user.Service.Init in bootstrap ordering.
func (s *Service) Init(ctx context.Context, adminEmail string) error {
	existing, err := s.tenantRepo.Filter(ctx, filter.New[model.Tenant](filter.Take(1)))
	if err != nil {
		return err
	}

	var tenant *model.Tenant
	if existing.IsEmpty() {
		tenant = &model.Tenant{Name: s.config.TenantName, Slug: s.config.TenantSlug}
		if err := s.tenantRepo.Create(ctx, tenant); err != nil {
			return err
		}
	} else {
		tenant = existing.First()
	}

	admin, err := s.userRepo.FindByEmail(ctx, adminEmail)
	if err != nil {
		var appErr *apperr.Error
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}

	if _, err := s.memberRepo.FindByUserAndTenant(ctx, admin.ID, tenant.ID); err == nil {
		return nil
	} else {
		var appErr *apperr.Error
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return s.memberRepo.Create(ctx, &model.Member{
		UserID:   admin.ID,
		TenantID: tenant.ID,
		Role:     model.RoleAdmin,
		IsActive: true,
	})
}

// FindByID requires the caller to be a member of the tenant (any role)
func (s *Service) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	identity, err := auth.CurrentIdentity(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.IsMemberOf(id) {
		return nil, apperr.Forbidden("not a member of this tenant")
	}

	return s.tenantRepo.FindByID(ctx, id)
}

// FindBySlug requires the caller to be a member of the tenant (any role)
func (s *Service) FindBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	identity, err := auth.CurrentIdentity(ctx)
	if err != nil {
		return nil, err
	}

	found, err := s.tenantRepo.Filter(ctx, filter.New[model.Tenant](
		filter.WhereEq("slug", slug),
		filter.Take(1),
	))
	if err != nil {
		return nil, err
	}
	if found.IsEmpty() {
		return nil, apperr.NotFound("tenant not found")
	}
	tenant := found.First()

	if !identity.IsMemberOf(tenant.ID) {
		return nil, apperr.Forbidden("not a member of this tenant")
	}

	return tenant, nil
}

// FindMine returns the tenants the current user is a member of.
func (s *Service) FindMine(ctx context.Context) (*TenantsPage, error) {
	identity, err := auth.CurrentIdentity(ctx)
	if err != nil {
		return nil, err
	}

	tenantIDs := make([]any, 0, len(identity.Memberships))
	for _, m := range identity.Memberships {
		tenantIDs = append(tenantIDs, m.TenantID)
	}
	if len(tenantIDs) == 0 {
		return &TenantsPage{Data: []model.Tenant{}}, nil
	}

	return s.tenantRepo.Filter(ctx, filter.New[model.Tenant](
		filter.WhereIn("id", tenantIDs...),
		filter.Take(len(tenantIDs)),
	))
}

// Create is gated by the license: unlicensed (or expired/invalid) instances
// may only ever have the one bootstrap Tenant; a licensed instance with a
// MaxTenants cap is rejected once that cap is reached. The creator becomes
// the new tenant's own admin automatically.
func (s *Service) Create(ctx context.Context, dto CreateTenantDTO) (*model.Tenant, error) {
	identity, err := auth.CurrentIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := s.tenantRepo.Filter(ctx, filter.New[model.Tenant](filter.WhereEq("slug", dto.Slug)))
	if err != nil {
		return nil, err
	}
	if existing.IsNotEmpty() {
		return nil, apperr.Conflict("slug already in use")
	}

	count, err := s.tenantRepo.Filter(ctx, filter.New[model.Tenant](filter.Take(1)))
	if err != nil {
		return nil, err
	}

	info := s.license.Info()
	if !info.Valid && count.Total > 0 {
		return nil, apperr.Forbidden("multi-tenant use requires a Snipet Enterprise License")
	}
	if info.Valid && info.MaxTenants > 0 && int(count.Total) >= info.MaxTenants {
		return nil, apperr.Forbidden("tenant limit reached for this license")
	}

	created := &model.Tenant{Name: dto.Name, Slug: dto.Slug, Icon: dto.Icon}
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.tenantRepo.Create(ctx, created); err != nil {
			return err
		}
		return s.memberRepo.Create(ctx, &model.Member{
			UserID:   identity.User.ID,
			TenantID: created.ID,
			Role:     model.RoleAdmin,
			IsActive: true,
		})
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) UpdateByID(ctx context.Context, id string, dto UpdateTenantDTO) error {
	if err := s.requireCanManage(ctx, id); err != nil {
		return err
	}

	updates := &model.Tenant{}
	if dto.Name != nil {
		updates.Name = *dto.Name
	}
	if dto.Slug != nil {
		updates.Slug = *dto.Slug
	}
	if dto.Icon != nil {
		updates.Icon = dto.Icon
	}
	return s.tenantRepo.UpdateByID(ctx, id, updates)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	if err := s.requireCanManage(ctx, id); err != nil {
		return err
	}
	return s.tenantRepo.DeleteByID(ctx, id)
}

func (s *Service) requireCanManage(ctx context.Context, tenantID string) error {
	identity, err := auth.CurrentIdentity(ctx)
	if err != nil {
		return err
	}

	if identity.IsPlatformAdmin() || identity.IsTenantAdmin(tenantID) {
		return nil
	}

	return apperr.Forbidden("not allowed to manage this tenant")
}

// Leave blocks if the caller is the tenant's last active admin.
func (s *Service) Leave(ctx context.Context, tenantID string) error {
	identity, err := auth.CurrentIdentity(ctx)
	if err != nil {
		return err
	}
	member, ok := identity.MembershipOf(tenantID)
	if !ok {
		return apperr.Forbidden("not a member of this tenant")
	}

	if member.Role == model.RoleAdmin {
		admins, err := s.memberRepo.Filter(ctx, filter.New[model.Member](
			filter.WhereEq("tenant_id", tenantID),
			filter.WhereEq("role", model.RoleAdmin),
			filter.WhereEq("is_active", true),
		))
		if err != nil {
			return err
		}
		if admins.Total <= 1 {
			return apperr.Conflict("you are the last admin of the tenant and cannot leave")
		}
	}

	return s.memberRepo.DeleteByID(ctx, member.ID)
}
