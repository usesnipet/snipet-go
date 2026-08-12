# Account module

CRUD for the `Account` entity (see `plan-ee/database.md`) — OAuth provider links (google, github, etc.) for a `User`. Package: `internal/module/account`

## Model fields

- ID (uuid)
- UserID (uuid, FK to User)
- Provider (string — google, github, etc.)
- ExternalID (string)
- CreatedAt / UpdatedAt

Composite unique: `(provider, external_id)`.

## Repository

`internal/repository/account.go`

```go
type IAccountRepository interface {
    IRepository[model.Account]

    FindByProviderAndExternalID(ctx context.Context, provider, externalID string) (*model.Account, error)
}

type AccountRepository struct {
    *Repository[model.Account]
}

func NewAccountRepository(db *gorm.DB) IAccountRepository
```

`FindByProviderAndExternalID` is used by the `auth` module (see `plan-ee/module/auth.md`) on the OAuth callback to check whether the identity is already linked to a `User`.

## Service

```go
type Service struct {
    accountRepo repository.IAccountRepository
}

func NewService(accountRepo repository.IAccountRepository) *Service

func (s *Service) Filter(ctx, filter *filter.Options[model.Account]) (*page.Paginated[model.Account], error)
func (s *Service) FindByID(ctx, id string) (*model.Account, error)
func (s *Service) FindByProviderAndExternalID(ctx, provider, externalID string) (*model.Account, error)
func (s *Service) Create(ctx, dto CreateAccountDTO) (*model.Account, error)
func (s *Service) UpdateByID(ctx, id string, dto UpdateAccountDTO) error
func (s *Service) Delete(ctx, id string) error
```

`Create` must check `(provider, external_id)` uniqueness before insert and return `apperr.Conflict` on collision.

## Handler

Routes under `/accounts`, ID-based:

| Method | Path | Handler |
|---|---|---|
| GET | `/accounts` | `filter` |
| GET | `/accounts/{id}` | `findByID` |
| POST | `/accounts` | `create` |
| PUT | `/accounts/{id}` | `update` |
| DELETE | `/accounts/{id}` | `delete` |

## DTOs

```go
type AccountResponse = model.Account
type AccountsPage = page.Paginated[model.Account]

type CreateAccountDTO struct {
    UserID     string `json:"user_id" validate:"required,uuid"`
    Provider   string `json:"provider" validate:"required,max=255"`
    ExternalID string `json:"external_id" validate:"required,max=255"`
}

type UpdateAccountDTO struct {
    ExternalID *string `json:"external_id" validate:"omitempty,max=255"`
}

type FindAccountsFilterDTO struct {
    Take *int `form:"take" validate:"omitempty,min=1"`
    Skip *int `form:"skip" validate:"omitempty,min=0"`
}
```

`UserID` and `Provider` are immutable after creation — not present on `UpdateAccountDTO` (relinking to a different user/provider means delete + recreate).
