# Member module

Join between `User` and `Tenant`, carrying the per-tenant `role` (see `plan-ee/database.md`). Sub-resource of `Tenant` — the public routes are nested under `/tenants/{tenant_id}/members`, not top-level. Package: `internal/module/member`.

## Model fields

- ID (uuid)
- UserID (uuid, FK to User)
- TenantID (uuid, FK to Tenant)
- Role (`model.Role` — typed string, see below)
- IsActive (bool)
- CreatedAt / UpdatedAt

Composite unique: `(user_id, tenant_id)`.

```go
type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
)
```

Only two roles for now (per this round's answer) — revisit once more granular permissions are needed (bitmask, custom roles, etc.); `Role` being a typed string rather than a raw column keeps that migration cheap later.

## Repository

`internal/repository/member.go`

```go
type IMemberRepository interface {
    IRepository[model.Member]

    FindByUserAndTenant(ctx context.Context, userID, tenantID string) (*model.Member, error)
    FilterWithUser(ctx context.Context, filter *filter.Options[model.Member]) (*page.Paginated[MemberWithUser], error)
}

type MemberRepository struct {
    *Repository[model.Member]
}

func NewMemberRepository(db *gorm.DB) IMemberRepository
```

`FindByUserAndTenant` backs the `(user_id, tenant_id)` unique lookup — used by `tenant`'s `canManage` permission check, `Leave`, and the auto-membership check in `tenant.Service.Init` (see `plan-ee/module/tenant.md`). `FilterWithUser` is the same `Filter` query with `Preload("User")` (or a join, GORM's choice) so the list endpoint can return the member's name/email/picture without the frontend making N follow-up calls.

## Service

```go
type Service struct {
    memberRepo repository.IMemberRepository
    userRepo   repository.IUserRepository
}

func NewService(memberRepo repository.IMemberRepository, userRepo repository.IUserRepository) *Service

func (s *Service) FindByTenant(ctx context.Context, tenantID string, dto FindMembersFilterDTO) (*page.Paginated[MemberResponse], error)
func (s *Service) FindByUserAndTenant(ctx, userID, tenantID string) (*model.Member, error)
func (s *Service) UpdateRole(ctx context.Context, tenantID, memberID string, dto UpdateMemberDTO) error
func (s *Service) Delete(ctx context.Context, tenantID, memberID string) error

// internal only — no public route; called by tenant.Service.Create and
// tenant-invitation.Service.Accept, not exposed as its own endpoint.
func (s *Service) create(ctx context.Context, member *model.Member) (*model.Member, error)
```

`FindByTenant` requires the caller to *be a member* of `tenantID` (any role) — this is a "see your team" screen, not an admin-only one, unlike `tenant-invitation.Filter`. `dto.IncludeUser` (see DTOs) picks `FilterWithUser` vs plain `Filter`. `UpdateRole`/`Delete` require `authz.CanManageTenant(ctx, tenantID)` (tenant admin or platform admin — see "Shared authorization helper" in `plan-ee/module/tenant.md`) and 404 if the target member's `TenantID` doesn't match the route's `{tenant_id}` (same not-found-not-forbidden reasoning as `tenant-invitation`).

There's no public `Create`/`POST /members` anymore — `Member` rows only come from `tenant.Service.Create` (auto-membership for the tenant creator) or `tenant-invitation.Service.Accept` (invite flow). Both call the unexported `create` directly; it still does the `(user_id, tenant_id)` uniqueness check internally.

**Last-admin guard — same invariant as `tenant.Service.Leave`** (see `plan-ee/module/tenant.md`), enforced here too since `UpdateRole`/`Delete` are the other two doors that can strip a tenant of its last admin:

- `UpdateRole` — if the target member's current `Role == model.RoleAdmin` and the update would change it (`dto.Role != nil && *dto.Role != model.RoleAdmin`) or deactivate it (`dto.IsActive != nil && *dto.IsActive == false`), count other active admins first: `memberRepo.Filter(WhereEq("tenant_id", member.TenantID), WhereEq("role", model.RoleAdmin), WhereEq("is_active", true))`. If `.Total <= 1`, reject with `apperr.Conflict("cannot change role: last admin of the tenant")` before applying the update.
- `Delete` — same check: if the member being deleted is `Role == model.RoleAdmin && IsActive`, and it's the only one (`.Total <= 1`), reject with `apperr.Conflict("cannot remove: last admin of the tenant")`.

Both count queries are the same shape as `tenant.Service.Leave`'s — kept as inline logic in each service (not factored into a shared helper) since `tenant.Service` depends on `repository.IMemberRepository` directly rather than on `member.Service`, matching how other modules in this codebase compose (e.g. `client.Service` takes `agentRepo` directly, not `agent.Service`).

## Handler

Routes nested under `/tenants/{tenant_id}/members` (mounted from the `tenant` handler, or a separate router mounted at that prefix — implementation detail):

| Method | Path | Handler | Auth |
|---|---|---|---|
| GET | `/tenants/{tenant_id}/members` | `filter` | bearer + member of `tenant_id` |
| PUT | `/tenants/{tenant_id}/members/{id}` | `updateRole` | bearer + `authz.CanManageTenant` |
| DELETE | `/tenants/{tenant_id}/members/{id}` | `delete` | bearer + `authz.CanManageTenant` |

## DTOs

```go
type MemberResponse struct {
    *model.Member
    User *model.User `json:"user,omitempty"` // populated only when include_user=true
}
type MembersPage = page.Paginated[MemberResponse]

type UpdateMemberDTO struct {
    Role     *model.Role `json:"role" validate:"omitempty,oneof=admin user"`
    IsActive *bool       `json:"is_active" validate:"omitempty"`
}

type FindMembersFilterDTO struct {
    IncludeUser *bool `form:"include_user" validate:"omitempty"`
    Take        *int  `form:"take" validate:"omitempty,min=1"`
    Skip        *int  `form:"skip" validate:"omitempty,min=0"`
}
```

`UserID` and `TenantID` are immutable — there's no DTO field for either since they're never client-settable (route param for `TenantID`, never re-assignable for `UserID`; moving a membership to a different user means delete + re-invite).
