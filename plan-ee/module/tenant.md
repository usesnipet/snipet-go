# Tenant module

CRUD for the `Tenant` entity (see `plan-ee/database.md`), plus single/multi-tenant mode bootstrapping and membership-aware create/update/delete/leave. Package: `internal/module/tenant`.

## Model fields

- ID (uuid)
- Name (string)
- Slug (string, unique)
- Icon (string, nullable)
- CreatedAt / UpdatedAt

## Config — single vs multi tenant

Whether more than one `Tenant` may exist is **not** an operator-trusted config flag — it's derived from a signed license key at boot. See `plan-ee/licensing.md` for the key format/verification mechanism; this module only consumes the result via `license.Service.Info()`.

```go
// config.TenantConfig (new)
type TenantConfig struct {
    TenantName string `env:"SINGLE_TENANT_NAME, default=Snipet"` // bootstrap tenant's name
    TenantSlug string `env:"SINGLE_TENANT_SLUG, default=default"` // bootstrap tenant's slug
}
```

`TenantName`/`TenantSlug` name the one `Tenant` created at bootstrap (`Service.Init` below) — needed regardless of license state, every instance has at least one `Tenant`. The license key itself lives in its own config (`config.LicenseConfig.LicenseKey`, see `plan-ee/licensing.md`), not here — `tenant.Service` depends on `*license.Service`, not on the raw key.

Same shape as the existing `config.AppConfig.InheritClient*` fields — `Service.Init` below mirrors `client.Service.Init`.

## Repository

`internal/repository/tenant.go`

```go
type ITenantRepository interface {
    IRepository[model.Tenant]
}

type TenantRepository struct {
    *Repository[model.Tenant]
}

func NewTenantRepository(db *gorm.DB) ITenantRepository
```

Generic CRUD only — existence checks use `Filter` + `.Total` (e.g. `Filter(ctx, filter.New[model.Tenant](filter.Take(1))).Total > 0`), same as other modules; no dedicated `Count`/`Exists` method needed.

## Service

```go
type Service struct {
    tenantRepo repository.ITenantRepository
    memberRepo repository.IMemberRepository // for auto-membership, permission checks, "my tenants"
    userRepo   repository.IUserRepository   // for the is_admin check
    txManager  repository.ITxManager
    config     config.TenantConfig
    license    *license.Service // see plan-ee/licensing.md — Info() returns cached {Valid, MaxTenants}
}

func NewService(tenantRepo repository.ITenantRepository, memberRepo repository.IMemberRepository, userRepo repository.IUserRepository, txManager repository.ITxManager, config config.TenantConfig, license *license.Service) *Service

func (s *Service) Init(ctx context.Context) error

func (s *Service) Filter(ctx, filter *filter.Options[model.Tenant]) (*page.Paginated[model.Tenant], error) // platform-admin only, enforced at handler
func (s *Service) FindByID(ctx, id string) (*model.Tenant, error)
func (s *Service) FindMine(ctx context.Context) (*page.Paginated[model.Tenant], error)
func (s *Service) Create(ctx context.Context, dto CreateTenantDTO) (*model.Tenant, error)
func (s *Service) UpdateByID(ctx context.Context, id string, dto UpdateTenantDTO) error
func (s *Service) DeleteByID(ctx context.Context, id string) error
func (s *Service) Leave(ctx context.Context, id string) error
```

**`Init`** (called from bootstrap, same place `clientService.Init` is called today) — always runs, regardless of license state: if no `Tenant` row exists yet, creates one with `Name: config.TenantName`, `Slug: config.TenantSlug`. Does *not* attach any member here — `Init` must also ensure the bootstrapped admin user (`plan-ee/module/user.md`'s `user.Service.Init`) is a `Member` of this tenant, since an unlicensed instance never reaches `Create` to do it. Simplest: `tenant.Service.Init` runs after `user.Service.Init` in bootstrap ordering, looks up the admin user by the configured admin email, and creates the `Member{Role: admin}` row if missing (idempotent — check via `memberRepo.FindByUserAndTenant` first).

**`Create`** — resolves current user id from `auth.GetPrincipal(ctx)` (bearer-only route). Checks `s.license.Info()` (see `plan-ee/licensing.md`): if unlicensed (or expired/invalid) and a `Tenant` already exists (it always will, post-`Init`), reject with `apperr.Forbidden("multi-tenant use requires a Snipet Enterprise License")`; if licensed with `MaxTenants > 0` and the current tenant count already reached it, reject with `apperr.Forbidden("tenant limit reached for this license")`. Otherwise, wrapped in `txManager.WithTransaction`: create the `Tenant`, then create `Member{UserID: currentUserID, TenantID: tenant.ID, Role: model.RoleAdmin, IsActive: true}` — the creator becomes the tenant's own admin automatically.

**`FindMine`** — `memberRepo.Filter` joined/filtered by `user_id = currentUserID`, mapped to their tenants (or `tenantRepo.Filter` with `WhereIn("id", memberTenantIDs)` — implementation detail, either join works).

**`UpdateByID`/`DeleteByID`** — both require `authz.CanManageTenant(ctx, tenantID)` (see "Shared authorization helper" below) to return `true`, else `apperr.Forbidden`. `DeleteByID` cascades to `Member`/`TenantInvitation` rows at the DB level (FK `ON DELETE CASCADE`, same convention as `refresh_tokens` → `users` today).

**`Leave`** — resolves current user id from principal, `memberRepo.FindByUserAndTenant(currentUserID, tenantID)` (404 if not a member). **Blocks if the caller is the tenant's last active admin:** if `member.Role == model.RoleAdmin`, count other active admins via `memberRepo.Filter(WhereEq("tenant_id", tenantID), WhereEq("role", model.RoleAdmin), WhereEq("is_active", true))` — if `.Total <= 1` (only this member), reject with `apperr.Conflict("cannot leave: last admin of the tenant")`. Otherwise `memberRepo.DeleteByID(member.ID)`.

## Shared authorization helper

`member` (`plan-ee/module/member.md`) and `tenant-invitation` (`plan-ee/module/tenant-invitation.md`) need the exact same "is this user allowed to manage this tenant" check. Rather than duplicate the two-lookup logic in three services, factor it as a stateless package-level function — not a service dependency, so it doesn't force `tenant.Service` to depend on `member.Service` or vice versa (each of the three services already holds `userRepo`+`memberRepo` directly):

```go
// internal/authz/authz.go
func CanManageTenant(
    ctx context.Context,
    userRepo repository.IUserRepository,
    memberRepo repository.IMemberRepository,
    currentUserID, tenantID string,
) (bool, error) {
    user, err := userRepo.FindByID(ctx, currentUserID)
    if err != nil {
        return false, err
    }
    if user.IsAdmin {
        return true, nil
    }
    member, err := memberRepo.FindByUserAndTenant(ctx, currentUserID, tenantID)
    if err != nil {
        if apperr.IsNotFound(err) {
            return false, nil
        }
        return false, err
    }
    return member.Role == model.RoleAdmin && member.IsActive, nil
}
```

`GET /tenants` platform-admin check (below) only needs the `user.IsAdmin` half — call `userRepo.FindByID` directly there rather than this helper (no `tenantID` to check against).

## Handler

Routes under `/tenants`:

| Method | Path | Handler | Auth |
|---|---|---|---|
| GET | `/tenants` | `filter` | bearer + platform-admin (`user.is_admin`) |
| GET | `/tenants/me` | `findMine` | bearer |
| GET | `/tenants/{id}` | `findByID` | bearer + member of `{id}` (or platform-admin) |
| POST | `/tenants` | `create` | bearer (blocked without a valid Multi-Tenant Use license, see `plan-ee/licensing.md`) |
| PUT | `/tenants/{id}` | `update` | bearer + `canManage` |
| DELETE | `/tenants/{id}` | `delete` | bearer + `canManage` |
| POST | `/tenants/{id}/leave` | `leave` | bearer + member of `{id}` |

`/tenants/me` registered before `/tenants/{id}` (same wildcard-ordering note as `/users/me`). `GET /tenants` stacks `mw.Bearer(jwtService)` + `mw.RequirePlatformAdmin(userRepo)` (see `plan-ee/module/middleware.md`) — the rest of the routes use just `mw.Bearer`, with `authz.CanManageTenant` (this doc) checked inline in the service where a `{tenant_id}` and finer-grained result are needed.

Also mounted here, as sub-routers: `plan-ee/module/member.md` (`/tenants/{tenant_id}/members/*`) and `plan-ee/module/tenant-invitation.md` (`/tenants/{tenant_id}/invitations/*`).

## DTOs

```go
type TenantResponse = model.Tenant
type TenantsPage = page.Paginated[model.Tenant]

type CreateTenantDTO struct {
    Name string `json:"name" validate:"required,max=255"`
    Slug string `json:"slug" validate:"required,max=255,alphanum_dash"` // lowercase, [a-z0-9-]
    Icon *string `json:"icon" validate:"omitempty,max=255"`
}

type UpdateTenantDTO struct {
    Name *string `json:"name" validate:"omitempty,max=255"`
    Slug *string `json:"slug" validate:"omitempty,max=255,alphanum_dash"`
    Icon *string `json:"icon" validate:"omitempty,max=255"`
}

type FindTenantsFilterDTO struct {
    Take *int `form:"take" validate:"omitempty,min=1"`
    Skip *int `form:"skip" validate:"omitempty,min=0"`
}
```

`Create` still checks slug uniqueness (`Filter` by `slug` before insert, same pattern as `client.GenerateCode`/`FindByCode`) and returns `apperr.Conflict` on collision.
