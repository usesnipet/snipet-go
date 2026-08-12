# User module (tenant staff)

> **Naming overlap — resolved.** `internal/module/user` and `model.User` already exist today for the *client-widget end-user* (visitor authenticating into a `Client` via OIDC/anonymous). This new `User` is a different actor (tenant staff logging into the dashboard) and both need to live in the same `internal/module`/`internal/model` packages — there's no `internal/ee/**` split anymore (see `plan-ee/boundary.md`). Decision: the existing client-widget module/entity gets renamed first, as a prerequisite, isolated commit — `internal/module/user` → `internal/module/clientuser`, `model.User` → `model.ClientUser` (and `internal/module/auth` → `internal/module/clientauth`, see `plan-ee/module/auth.md`). That frees `internal/module/user`/`model.User` for this doc's entity.

CRUD for the `User` entity (see `plan-ee/database.md`). Package: `internal/module/user`.

## Model fields

- ID (uuid)
- Name (string)
- Email (string, unique)
- PasswordHash (string, nullable — null when the user only has OAuth `Account`s)
- Picture (string, nullable)
- IsAdmin (bool) — platform-level flag, not tenant-scoped
- Challenges (jsonb array — `active_account`, `change_password`, etc.)
- CreatedAt / UpdatedAt

## Repository

`internal/repository/user.go`

```go
type IUserRepository interface {
    IRepository[model.User]
    FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository struct {
    *Repository[model.User]
}

func NewUserRepository(db *gorm.DB) IUserRepository
```

`FindByEmail` is included since email is the natural login/lookup key (same rationale as `client.FindByCode`); the auth module will call it once login is designed.

## Config — bootstrap admin

```go
// config.UserConfig (new)
type UserConfig struct {
    AdminName     string `env:"ADMIN_NAME, default=Admin"`
    AdminEmail    string `env:"ADMIN_EMAIL, default=admin@admin.com"`
    AdminPassword string `env:"ADMIN_PASSWORD"` // empty => generate a random one at boot
}
```

Separate from `config.TenantConfig` (`plan-ee/module/tenant.md`) and `config.LicenseConfig` (`plan-ee/licensing.md`) — each module owns its own config struct, no shared `EEConfig` grab-bag (there's no `ee` grouping anymore, see `plan-ee/boundary.md`).

## Service

```go
type Service struct {
    userRepo repository.IUserRepository
    config   config.UserConfig
}

func NewService(userRepo repository.IUserRepository, config config.UserConfig) *Service

func (s *Service) Init(ctx context.Context) error

func (s *Service) Filter(ctx, filter *filter.Options[model.User]) (*page.Paginated[model.User], error)
func (s *Service) FindByID(ctx, id string) (*model.User, error)
func (s *Service) FindByEmail(ctx, email string) (*model.User, error)
func (s *Service) Create(ctx, dto CreateUserDTO) (*model.User, error)
func (s *Service) UpdateByID(ctx, id string, dto UpdateUserDTO) error
func (s *Service) Delete(ctx, id string) error

func (s *Service) Me(ctx context.Context) (*model.User, error)
func (s *Service) UpdateMyPicture(ctx context.Context, dto UpdateProfilePictureDTO) error
```

`Create` hashes `dto.Password` into `PasswordHash` before insert (never accept a pre-hashed value from the client) and checks email uniqueness. Password *change* is a separate, non-CRUD flow — lives in `auth` (see `plan-ee/module/auth.md`), not here, same as `api-key`'s `roll`/`toggle` pattern for non-base operations.

`Me`/`UpdateMyPicture` resolve the current user id from `auth.GetPrincipal(ctx)` (bearer JWT), same pattern as `api-key.Service.Me`.

**`Init`** (called from bootstrap, before `tenant.Service.Init` — see `plan-ee/module/tenant.md`, which needs this user to already exist so it can add them as the single tenant's first `Member`): checks whether *any* `User` row exists (`Filter` + `.Total > 0`, same existence-check idiom as `tenant.Service.Init`). If none exists, creates one admin user:

- `Name`/`Email` from `config.AdminName`/`config.AdminEmail` (defaults `Admin` / `admin@admin.com`).
- `Password`: `config.AdminPassword` if set; otherwise generate a random one (reuse the same opaque-random-string approach as `auth.RefreshTokenService.GenerateRefreshToken`, or `client.GenerateCode`'s `crypto/rand` pattern).
- `IsAdmin: true`, `Challenges: []` (bootstrap account is considered already activated — it never goes through `auth.Register`, so there's no activation email to gate on).
- **If the password was generated** (not supplied via env), print it once to stdout at the end of `Init` — plainly, not through the structured `logger` (a generated one-time secret shouldn't risk ending up shipped to a remote log aggregator; a direct console print is the standard pattern for this, e.g. Rancher/Gitea do the same for their bootstrap admin). Never logged/printed again after this — only the hash persists.

## Handler

Routes under `/users`. `POST/PUT/DELETE /users*` are admin-only (apiKey/platform-admin middleware — provisioning users directly, distinct from self-service signup in `auth`); `/users/me*` are bearer-JWT, self-service:

| Method | Path | Handler | Auth |
|---|---|---|---|
| GET | `/users` | `filter` | admin |
| GET | `/users/{id}` | `findByID` | admin |
| POST | `/users` | `create` | admin |
| PUT | `/users/{id}` | `update` | admin |
| DELETE | `/users/{id}` | `delete` | admin |
| GET | `/users/me` | `me` | bearer (self) |
| PUT | `/users/me/picture` | `updatePicture` | bearer (self) |

`/users/me*` routes must be registered before `/users/{id}` so `me` isn't swallowed by the `{id}` wildcard.

## DTOs

```go
type UserResponse = model.User
type UsersPage = page.Paginated[model.User]

type CreateUserDTO struct {
    Name     string  `json:"name" validate:"required,max=255"`
    Email    string  `json:"email" validate:"required,email,max=255"`
    Password string  `json:"password" validate:"required,min=8"`
    Picture  *string `json:"picture" validate:"omitempty,max=255"`
    IsAdmin  bool    `json:"is_admin"`
}

type UpdateUserDTO struct {
    Name    *string `json:"name" validate:"omitempty,max=255"`
    Email   *string `json:"email" validate:"omitempty,email,max=255"`
    Picture *string `json:"picture" validate:"omitempty,max=255"`
    IsAdmin *bool   `json:"is_admin" validate:"omitempty"`
}

type FindUsersFilterDTO struct {
    Take *int `form:"take" validate:"omitempty,min=1"`
    Skip *int `form:"skip" validate:"omitempty,min=0"`
}

type UpdateProfilePictureDTO struct {
    Picture string `json:"picture" validate:"required,max=255"`
}
```

`Challenges` is server-managed (set on create, mutated by dedicated flows like account activation) — not exposed on `CreateUserDTO`/`UpdateUserDTO`.