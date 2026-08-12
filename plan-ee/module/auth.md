# Auth module (tenant staff)

Not CRUD — orchestrates `User`, `Account`, and `Token` (each already covered in their own module doc) for registration, login, OAuth, password recovery, and account activation. Package: `internal/module/auth`. This name is currently held by the client-widget auth module — renaming that one to `internal/module/clientauth` (see `plan-ee/module/user.md`'s naming-overlap note) is a prerequisite, isolated commit before this module can be built.

## Scope

- Register with email + password + name (account starts inactive)
- Log in with email + password (blocked until activated)
- Log in / link via external provider (google, github, etc.)
- Refresh access token
- Logout (revoke refresh token)
- Set/change password (logged in)
- Forgot password / reset password (logged out, via emailed token)
- Activate account (via emailed token)

Out of scope, handled elsewhere: `GET /users/me`, profile picture update (`plan-ee/module/user.md`); sending the actual emails (`plan-ee/module/email.md`).

## Dependencies

- `internal/auth.HashPassword` / `ComparePassword` — reused as-is, no `model.User` coupling.
- `internal/auth.RefreshTokenService` — `GenerateRefreshToken()` / `HashRefreshToken()`, reused as-is for opaque token generation (used for *all* token types now, not just refresh — name is legacy, behavior is generic).
- `repository.IUserRepository` — `FindByEmail`, `Create`, `UpdateByID`.
- `repository.IAccountRepository` — `FindByProviderAndExternalID`, `Create`.
- `repository.ITokenRepository` — `Create`, `FindByHashAndType`, `RevokeByID` (see `plan-ee/module/token.md` — replaces the earlier `RefreshToken`-only repo).
- `email.Service` — `SendTemplate` (see `plan-ee/module/email.md`).
- `*coreauth.JWTService[*UserClaims]` — see "Token issuance" below.

## Token issuance — genericize `internal/auth/jwt.go`

Today `JWTService`/`UserClaims` are hard-coupled to the client-widget flow (`ClientCode` baked into the claims, `GenerateToken(clientCode string, user *model.User)` typed to the client `model.User`). Decision: make the base generic, move `ClientCode` out into the client module's own claims type, and have this module (tenant staff) define its own (currently empty, ready for future fields):

```go
// internal/auth/jwt.go — generic base, shared
type Claims interface {
    jwt.Claims
}

type BaseClaims struct {
    jwt.RegisteredClaims
}

func NewBaseClaims(cfg config.AuthConfig, subject string) BaseClaims {
    now := time.Now()
    return BaseClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    cfg.JWTIssuer,
            Subject:   subject,
            Audience:  jwt.ClaimStrings{cfg.JWTAudience},
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWTExpiration)),
        },
    }
}

type JWTService[T Claims] struct {
    config     config.AuthConfig
    newClaims  func() T // factory so VerifyToken has something to unmarshal into
}

func NewJWTService[T Claims](config config.AuthConfig, newClaims func() T) *JWTService[T]

func (s *JWTService[T]) GenerateToken(claims T) (string, time.Time, error)
func (s *JWTService[T]) VerifyToken(tokenString string) (T, error)
```

```go
// internal/module/clientauth (renamed from internal/module/auth, see the
// package-name note above) — keeps ClientCode, now explicit
type UserClaims struct {
    coreauth.BaseClaims
    ClientCode string `json:"client_code"`
}
// jwtService := coreauth.NewJWTService[*UserClaims](cfg.Auth, func() *UserClaims { return &UserClaims{} })
```

```go
// internal/module/auth (this module, tenant staff) — no extra fields today,
// struct kept so tenant-scoped claims (e.g. an "active tenant" id, if that
// design is ever needed) can be added later without touching the generic
// JWTService.
type UserClaims struct {
    coreauth.BaseClaims
}
// jwtService := coreauth.NewJWTService[*UserClaims](cfg.Auth, func() *UserClaims { return &UserClaims{} })
```

**Blast radius of this change** (found via grep on `auth.UserClaims`/`auth.JWTService`/`Principal`): `internal/middleware/jwt.go`, `internal/middleware/any-auth.go`, `internal/bootstrap/bootstrap.go`, and `internal/auth/context.go`'s `Principal` all reference the concrete (non-generic) types today. `Principal.JWTClaims` will need to hold `jwt.Claims` (interface) instead of `*UserClaims` so both the (soon-to-be) `clientauth` and this module's `JWTService` instances can populate the same `Principal`/context machinery. Worth doing this refactor as its own isolated commit before building this module, since it touches existing client-widget auth code. See `plan-ee/module/middleware.md` for how the tenant-staff side actually consumes this (`Bearer`/`RequirePlatformAdmin` middleware, `coreauth.CurrentUserID` helper, moved into `internal/auth` since it's generic enough for any `Principal`) — it reuses `Principal` as-is rather than standing up a parallel context mechanism.

## One-time tokens (activation, reset) now go through `Token`, not a stateless JWT

Since `RefreshToken` became `Token` with a `type` column (`refresh`, `activate_account`, `reset_password` — see `plan-ee/module/token.md`), action tokens are DB-backed like refresh tokens, not signed JWTs. Same opaque-random-string + SHA-256-hash-for-storage approach as refresh tokens (reuse `RefreshTokenService.GenerateRefreshToken()`/`HashRefreshToken()`), just tagged with a different `Type` and shorter `ExpiresAt`. Benefit over the earlier stateless-JWT idea: revocable and single-use (mark `RevokedAt` immediately on consumption), and one code path (`ITokenRepository`) for all token kinds.

## Service

```go
type Service struct {
    userRepo    repository.IUserRepository
    accountRepo repository.IAccountRepository
    tokenRepo   repository.ITokenRepository
    jwtService  *coreauth.JWTService[*UserClaims]
    tokenGen    *auth.RefreshTokenService // reused from internal/auth
    emailSvc    *email.Service
    providerRegistry *ProviderRegistry // see "OAuth providers" below
}

func NewService(...) *Service

func (s *Service) Register(ctx, dto RegisterDTO) (*RegisterResponse, error)
func (s *Service) Login(ctx, dto LoginDTO) (*AuthenticateResponse, error)
func (s *Service) GetAuthorizationURL(ctx, provider ProviderName) (string, error)
func (s *Service) AuthenticateCallback(ctx, provider ProviderName, code string) (*AuthenticateResponse, error)
func (s *Service) Refresh(ctx, dto RefreshDTO) (*AuthenticateResponse, error)
func (s *Service) Logout(ctx, dto RefreshDTO) error
func (s *Service) SetPassword(ctx context.Context, dto SetPasswordDTO) error // current user from principal
func (s *Service) ForgotPassword(ctx, dto ForgotPasswordDTO) error
func (s *Service) ResetPassword(ctx, dto ResetPasswordDTO) error
func (s *Service) Activate(ctx, dto ActivateAccountDTO) error
func (s *Service) ResendActivation(ctx, dto ResendActivationDTO) error
```

Notes per method (decisions from this round in **bold**):

- `Register` — checks email uniqueness, hashes password, creates `User` with `Challenges: ["active_account"]`. **Does not issue access/refresh tokens.** Issues a `Token{Type: activate_account}`, emails the activation link, returns a minimal `RegisterResponse` (no tokens) instead of `AuthenticateResponse`.
- `Login` — `FindByEmail`, reject with generic "invalid credentials" if not found or `PasswordHash == nil` (OAuth-only account — don't leak which). **If `Challenges` still contains `active_account`, reject with a distinct "account not activated" error** (safe to be specific here since the caller already proved they know the password) before comparing further/issuing tokens.
- `AuthenticateCallback` — exchange `code` for provider identity; `FindByProviderAndExternalID`: found → load `User`, issue tokens (OAuth logins skip the activation gate — the provider already verified the email). Not found → `FindByEmail` to link to an existing account, else create a new `User` (`PasswordHash: nil`, `Challenges: []` — no activation needed) + `Account`, issue tokens.
- `Refresh`/`Logout` — same shape as existing `internal/module/auth`, just against `ITokenRepository` with `Type: refresh`: `FindByHashAndType` → check `RevokedAt`/`ExpiresAt` → `RevokeByID` → reissue on refresh; just `RevokeByID` on logout.
- `SetPassword` — **decided: this is a "set", not a "change".** No current-password check at all, regardless of whether `PasswordHash` was already set. (Simpler and consistent for both OAuth-only users setting their first password *and* existing users changing it; if you want current-password verification for the "already has a password" case later, that's an easy follow-up, not a blocker now.)
- `ForgotPassword` — always responds success regardless of whether the email exists (avoid user enumeration); if found, issues `Token{Type: reset_password}` and emails the link; no-ops silently otherwise.
- `ResetPassword` — `FindByHashAndType(hash, reset_password)`, check not expired/revoked, hash `NewPassword`, `RevokeByID` (single-use).
- `Activate` — `FindByHashAndType(hash, activate_account)`, check not expired/revoked, remove `active_account` from `Challenges`, `RevokeByID` (single-use).
- `ResendActivation` — same as `Register`'s token-issuing step, for a user who's already registered but lost the email; revokes any prior un-consumed `activate_account` token for that user first (avoid stacking valid tokens).

## OAuth providers

Existing `internal/module/clientauth/auth-provider` (post-rename, see the package-name note above) is designed for **per-`Client` BYO OIDC** (issuer/audience configured on `model.ClientConfig`). Tenant-staff social login is different: fixed platform-level app credentials (Snipet's own registered Google/GitHub OAuth apps), not configured per tenant. Don't reuse the `Registry`/`IProvider` abstraction as-is — a smaller, config-driven provider set is enough:

```go
type ProviderName string
const (
    ProviderGoogle ProviderName = "google"
    ProviderGithub ProviderName = "github"
)

type Identity struct {
    ExternalID string
    Email      string
    Name       string
    Picture    string
}

type IProvider interface {
    Name() ProviderName
    AuthorizationURL(state string) string
    Exchange(ctx context.Context, code string) (*Identity, error)
}
```

Credentials (`GoogleClientID/Secret`, `GithubClientID/Secret`, redirect URLs) go in `config.AuthConfig`, platform-wide.

## Handler

Routes under `/auth`:

| Method | Path | Handler | Auth |
|---|---|---|---|
| POST | `/auth/register` | `register` | public |
| POST | `/auth/login` | `login` | public |
| GET | `/auth/{provider}` | `authorizationURL` | public |
| GET | `/auth/{provider}/callback` | `callback` | public |
| POST | `/auth/refresh` | `refresh` | public (refresh token in body) |
| POST | `/auth/logout` | `logout` | public (refresh token in body) |
| PUT | `/auth/password` | `setPassword` | bearer |
| POST | `/auth/password/forgot` | `forgotPassword` | public |
| POST | `/auth/password/reset` | `resetPassword` | public |
| POST | `/auth/activate` | `activate` | public |
| POST | `/auth/activate/resend` | `resendActivation` | public (rate-limit by email, avoid abuse) |

## DTOs

```go
type UserResponse = model.User

type AuthenticateResponse struct {
    AccessToken           string       `json:"access_token"`
    AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
    RefreshToken          string       `json:"refresh_token"`
    RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
    User                  UserResponse `json:"user"`
}

type RegisterResponse struct {
    User UserResponse `json:"user"`
    // no tokens — login is blocked until the account is activated
}

type RegisterDTO struct {
    Name     string `json:"name" validate:"required,max=255"`
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginDTO struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required"`
}

type RefreshDTO struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

type SetPasswordDTO struct {
    NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ForgotPasswordDTO struct {
    Email string `json:"email" validate:"required,email,max=255"`
}

type ResetPasswordDTO struct {
    Token       string `json:"token" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ActivateAccountDTO struct {
    Token string `json:"token" validate:"required"`
}

type ResendActivationDTO struct {
    Email string `json:"email" validate:"required,email,max=255"`
}

type ProviderCallbackQueryDTO struct {
    Code  string `form:"code" validate:"required"`
    State string `form:"state" validate:"required"`
}
```

## Remaining open question

None blocking base implementation. One design nit for later: `Login`'s "account not activated" error currently exposes that the email/password pair is valid but inactive — intentional per this round's answer (item 2's spirit: keep flows simple), but worth a rate-limit/lockout pass alongside the rest of the brute-force protections whenever those get designed.
