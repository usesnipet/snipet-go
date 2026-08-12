# TenantInvitation module

Sub-resource of `Tenant` — every route is nested under `/tenants/{tenant_id}/invitations`, not top-level. Package: `internal/module/tenantinvitation`.

## Model fields

- ID (uuid)
- TenantID (uuid, FK to Tenant)
- Email (string)
- Token (string, unique) — opaque secret sent in the email link, proof of receipt
- Role (`model.Role`, see `plan-ee/module/member.md`) — role the invitee gets on accept
- Status (`model.InvitationStatus` — typed string, see below)
- ExpiresAt (timestamp)
- CreatedAt / UpdatedAt

```go
type InvitationStatus string

const (
    InvitationStatusPending  InvitationStatus = "pending"
    InvitationStatusAccepted InvitationStatus = "accepted"
    InvitationStatusDeclined InvitationStatus = "declined"
)
```

`expired` is **not** a stored status — it's `Status == pending && ExpiresAt < now`, computed at query time (see `Filter` below). `cancelled` isn't stored either — canceling an invite is a hard delete (see `Delete`), it never got a real response so there's nothing worth keeping.

## Repository

`internal/repository/tenant-invitation.go`

```go
type ITenantInvitationRepository interface {
    IRepository[model.TenantInvitation]

    FindPendingByEmailAndTenant(ctx context.Context, email, tenantID string) (*model.TenantInvitation, error)
}

type TenantInvitationRepository struct {
    *Repository[model.TenantInvitation]
}

func NewTenantInvitationRepository(db *gorm.DB) ITenantInvitationRepository
```

`FindPendingByEmailAndTenant` backs the duplicate-invite guard in `Create` (see below). Accept/decline/cancel/list all resolve a specific invitation via the generic `FindByID` — the service then checks `invitation.TenantID == tenantIDFromRoute` itself (404, not 403, if it doesn't match — avoids leaking "this id exists but belongs to another tenant").

## Service

```go
type Service struct {
    invitationRepo repository.ITenantInvitationRepository
    memberRepo     repository.IMemberRepository
    userRepo       repository.IUserRepository
    tenantRepo     repository.ITenantRepository
    txManager      repository.ITxManager
    emailSvc       *email.Service
}

func NewService(...) *Service

func (s *Service) Filter(ctx context.Context, tenantID string, filter InvitationFilter) (*page.Paginated[model.TenantInvitation], error)
func (s *Service) Create(ctx context.Context, tenantID string, dto CreateTenantInvitationDTO) (*model.TenantInvitation, error)
func (s *Service) Accept(ctx context.Context, tenantID, invitationID string, dto AcceptInvitationDTO) error
func (s *Service) Decline(ctx context.Context, tenantID, invitationID string, dto DeclineInvitationDTO) error
func (s *Service) Delete(ctx context.Context, tenantID, invitationID string) error
```

Notes per method:

- **`Create`** — requires `authz.CanManageTenant(ctx, tenantID)` (tenant admin or platform admin — same check as `tenant.Service.UpdateByID`, see "Shared authorization helper" in `plan-ee/module/tenant.md`). Rejects with `apperr.Conflict` if `FindPendingByEmailAndTenant` already returns a row (don't stack duplicate invites — caller cancels the old one first, or waits for it to expire). Generates `Token` (opaque random, reuse `RefreshTokenService.GenerateRefreshToken()`-style generator), `Status: pending`, `ExpiresAt` defaulted (e.g. now + 7 days) if not supplied. Sends the invite email via `emailSvc.SendTemplate` with the tenant name + accept link (`TemplateTenantInvitation`, see `plan-ee/module/email.md`).
- **`Filter`** — requires `authz.CanManageTenant(ctx, tenantID)` (invite lists include email addresses — admin-only, not visible to regular members). `InvitationFilter.Status` accepts `pending`/`accepted`/`declined`/`expired`. Translated to the query as: `pending` → `status = pending AND expires_at >= now()`; `expired` → `status = pending AND expires_at < now()`; `accepted`/`declined` → plain `status = ...` (already terminal, expiry irrelevant).
- **`Accept`**/**`Decline`** — bearer auth, any authenticated user (this is *their* invitation, not a tenant-admin action). Loads the invitation by ID, checks `TenantID` matches the route, `Status == pending`, not expired, and — the actual security check — `dto.Token` matches the stored `Token` (constant-time compare) **and** the caller's own email (from `userRepo.FindByID(currentUserID)`, not client input) matches `invitation.Email`. Both conditions matter: the token proves they received *this* email, the email match proves they're logged in as *that* account (not just someone who intercepted the link). `Accept`, wrapped in `txManager.WithTransaction`: creates `Member{UserID: currentUserID, TenantID, Role: invitation.Role, IsActive: true}` (or reactivates/updates role if a `Member` row already exists — e.g. they left and got reinvited) and sets `Status: accepted`. `Decline` just sets `Status: declined`. Accept/decline is a `POST` (not `GET`) specifically so an email client's link-prefetcher can't silently consume it.
- **`Delete`** (cancel) — requires `authz.CanManageTenant(ctx, tenantID)`. Hard-deletes the row regardless of current `Status` (admin cleanup) — this is the only place invitations actually disappear; accepted/declined ones otherwise stay for history/audit in `Filter`.

## Handler

Routes nested under `/tenants/{tenant_id}/invitations` (mounted from the `tenant` handler, or a separate router mounted at that prefix — implementation detail):

| Method | Path | Handler | Auth |
|---|---|---|---|
| GET | `/tenants/{tenant_id}/invitations` | `filter` | bearer + `authz.CanManageTenant` |
| POST | `/tenants/{tenant_id}/invitations` | `create` | bearer + `authz.CanManageTenant` |
| POST | `/tenants/{tenant_id}/invitations/{id}/accept` | `accept` | bearer (invitee, email+token checked in service) |
| POST | `/tenants/{tenant_id}/invitations/{id}/decline` | `decline` | bearer (invitee, email+token checked in service) |
| DELETE | `/tenants/{tenant_id}/invitations/{id}` | `delete` | bearer + `authz.CanManageTenant` |

## DTOs

```go
type TenantInvitationResponse = model.TenantInvitation
type TenantInvitationsPage = page.Paginated[model.TenantInvitation]

type CreateTenantInvitationDTO struct {
    Email     string     `json:"email" validate:"required,email,max=255"`
    Role      model.Role `json:"role" validate:"required,oneof=admin user"`
    ExpiresAt *time.Time `json:"expires_at" validate:"omitempty"`
}

type AcceptInvitationDTO struct {
    Token string `json:"token" validate:"required"`
}

type DeclineInvitationDTO struct {
    Token string `json:"token" validate:"required"`
}

type FindTenantInvitationsFilterDTO struct {
    Status *model.InvitationStatus `form:"status" validate:"omitempty,oneof=pending accepted declined expired"`
    Take   *int `form:"take" validate:"omitempty,min=1"`
    Skip   *int `form:"skip" validate:"omitempty,min=0"`
}
```

`TenantID` is no longer on `CreateTenantInvitationDTO` — it comes from the route (`{tenant_id}`), not the body, now that this is a proper sub-resource. `Token` is never present on `CreateTenantInvitationDTO` — server-generated only, only ever sent back out via email, never in an API response body (`TenantInvitationResponse`/`model.TenantInvitation` should exclude it from JSON, e.g. `json:"-"`, same principle as `PasswordHash` on `User`).
