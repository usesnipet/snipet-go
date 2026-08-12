# Token module

CRUD for the `Token` entity (see `plan-ee/database.md`) — replaces the earlier `RefreshToken`-only design. Same shape as a refresh token, plus a `type` discriminator so the same table backs refresh tokens *and* one-time action tokens (account activation, password reset). Package: `internal/module/token`.

## Model fields

- ID (uuid)
- Type (string — `refresh`, `activate_account`, `reset_password`)
- Hash (string)
- UserID (uuid, FK to User)
- ExpiresAt (timestamp)
- RevokedAt (timestamp, nullable)
- Metadata (jsonb)
- CreatedAt

## Repository

`internal/repository/token.go`

```go
type ITokenRepository interface {
    IRepository[model.Token]

    FindByHashAndType(ctx context.Context, hash string, tokenType model.TokenType) (*model.Token, error)
    RevokeByID(ctx context.Context, id string) error
}

type TokenRepository struct {
    *Repository[model.Token]
}

func NewTokenRepository(db *gorm.DB) ITokenRepository
```

`FindByHashAndType` mirrors the old `FindByHash` but also filters on `type`, so a stolen/leaked activation token can't accidentally be replayed as a refresh token (and vice versa) even if the hash collided. `model.TokenType` is a typed string constant (`refresh`, `activate_account`, `reset_password`) — see `plan-ee/module/auth.md` for how each is issued/consumed.

## Service

```go
type Service struct {
    tokenRepo repository.ITokenRepository
}

func NewService(tokenRepo repository.ITokenRepository) *Service

func (s *Service) Filter(ctx, filter *filter.Options[model.Token]) (*page.Paginated[model.Token], error)
func (s *Service) FindByID(ctx, id string) (*model.Token, error)
func (s *Service) FindByHashAndType(ctx, hash string, tokenType model.TokenType) (*model.Token, error)
func (s *Service) Create(ctx, dto CreateTokenDTO) (*model.Token, error)
func (s *Service) UpdateByID(ctx, id string, dto UpdateTokenDTO) error
func (s *Service) RevokeByID(ctx, id string) error
func (s *Service) Delete(ctx, id string) error
```

Issuing/validating/consuming tokens as part of login, refresh, activation, or password-reset flows is `auth`-module business logic (it calls `Create`/`FindByHashAndType`/`RevokeByID` directly) — this module only exposes persistence.

## Handler

Routes under `/tokens`, ID-based, admin/audit-only (same rationale as the old `refresh-tokens` route — the `auth` module issues/consumes tokens directly via the service, not through this public route):

| Method | Path | Handler |
|---|---|---|
| GET | `/tokens` | `filter` |
| GET | `/tokens/{id}` | `findByID` |
| POST | `/tokens` | `create` |
| PUT | `/tokens/{id}` | `update` |
| DELETE | `/tokens/{id}` | `delete` |

## DTOs

```go
type TokenResponse = model.Token
type TokensPage = page.Paginated[model.Token]

type CreateTokenDTO struct {
    UserID    string        `json:"user_id" validate:"required,uuid"`
    Type      model.TokenType `json:"type" validate:"required,oneof=refresh activate_account reset_password"`
    Hash      string        `json:"hash" validate:"required,max=255"`
    ExpiresAt time.Time     `json:"expires_at" validate:"required"`
    Metadata  jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type UpdateTokenDTO struct {
    RevokedAt *time.Time     `json:"revoked_at" validate:"omitempty"`
    Metadata  *jsonx.JSONMap `json:"metadata" validate:"omitempty"`
}

type FindTokensFilterDTO struct {
    Take *int `form:"take" validate:"omitempty,min=1"`
    Skip *int `form:"skip" validate:"omitempty,min=0"`
}
```

`Hash`, `UserID`, and `Type` are immutable after creation. Revoking/consuming a token is done via `UpdateByID`/`RevokeByID` setting `RevokedAt` — action tokens (`activate_account`, `reset_password`) get revoked immediately on successful use so they're single-use; `refresh` tokens get revoked on rotation/logout, same as before.
