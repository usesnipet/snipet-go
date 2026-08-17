# `internal/auth` + `internal/guard` + `internal/authz`

Three layers, three different jobs:

- **`internal/auth`** — JWT signing/verification, API key hashing, and the
  per-mechanism *identity* types + context helpers (below). Knows nothing
  about HTTP.
- **`internal/guard`** — composable `Gate`s that run those primitives against
  an incoming request and stash the resulting identity on the context.
  Knows nothing about tenants.
- **`internal/authz`** — tenant-membership checks (`RequireMember`,
  `RequireAdmin`) that services call *after* a guard has already
  authenticated the request. Knows nothing about HTTP or which guard ran.

A request authenticates via **exactly one** mechanism (unlike a
multi-principal design) — the three identity types below are mutually
exclusive on a given request's context.

## Mechanisms + identity types

| Mechanism | Credential | Identity type | Set by |
|---|---|---|---|
| Platform-staff JWT | `Authorization: Bearer` | `auth.UserIdentity` (`User`, `Memberships []*model.Member`) | `guard.RequireUser` |
| API key | `X-API-Key` | `auth.ApiKeyIdentity` (`APIKeyID`, `TenantID`) | `guard.RequireApiKey` |
| Client-widget JWT | `Authorization: Bearer` | `auth.ClientUserIdentity` (`UserID`, `ClientCode`) | `guard.RequireClientUser` |

```go
type UserIdentity struct {
    User        *model.User
    Memberships []*model.Member
}

type ApiKeyIdentity struct {
    APIKeyID string
    TenantID string
}

type ClientUserIdentity struct {
    UserID     string
    ClientCode string
}
```

## Reading who's authenticated

```go
identity, err := auth.CurrentUser(ctx)       // *auth.UserIdentity
identity, err := auth.CurrentApiKey(ctx)     // auth.ApiKeyIdentity
identity, err := auth.CurrentClientUser(ctx) // auth.ClientUserIdentity

auth.HasApiKey(ctx) // bool, no error — useful in code paths reachable by
                     // more than one mechanism (see session.Service.Run)
```

Each getter errors if that mechanism didn't authenticate the request — a
service only calls the one matching whatever guard its route uses.

`UserIdentity` carries `IsMemberOf(tenantID)`, `IsTenantAdmin(tenantID)`,
`IsPlatformAdmin()` — the primitives `internal/authz` builds on.

## Guards — require gates + `Or`

```go
requireUser := guard.RequireUser(userJWTService, userRepo, memberRepo)
requireApiKey := guard.RequireApiKey(apiKeyService, apiKeyCache)
requireClientUser := guard.RequireClientUser(clientJWTService)
anyClientAuth := guard.Or(requireClientUser, requireApiKey)

r.Use(requireUser.Handler())
r.Use(anyClientAuth.Handler())
```

`Or` tries each gate in order — the first one whose credential is present
decides the outcome:

- credential **absent** → try the next gate (`auth.ErrNotApplicable`)
- credential **present but invalid** → hard 401 (does not fall through)
- credential **valid** → sets that identity on the context, done

All gates are built once in `bootstrap.Bootstrap` (see
[bootstrap.md](./bootstrap.md)) and handed to modules as
`api.MiddlewareFunc` via `.Handler()` — a module's `handler.go` decides
which gate(s) to `r.Use(...)` per route group (see
[modules.md](./modules.md#handlergo)). A single module can mix gates
across route groups when some routes need different auth than others —
`agent`'s CRUD routes use `requireUser` while its `/run` route keeps
`requireApiKey` (see below).

## `internal/authz` — tenant membership, on top of `guard.RequireUser`

Every module whose admin/CRUD surface lives under
`/tenants/{tenant_id}/...` (`member`, `tenant`, `agent`, `api-key`,
`client`, `knowledge`, `knowledge-index`) is gated by `guard.RequireUser`
at the route level, then calls one of two package-level checks at the top
of each service method:

```go
func RequireMember(ctx context.Context, tenantID string) (*auth.UserIdentity, error)
func RequireAdmin(ctx context.Context, tenantID string) (*auth.UserIdentity, error)
```

`RequireMember` accepts any active role in the tenant; `RequireAdmin`
requires `model.RoleAdmin`. Both read `auth.CurrentUser(ctx)` and return
`apperr.Forbidden(...)` on failure — same shape either way, so a service
just does:

```go
func (s *Service) Create(ctx context.Context, tenantID string, dto CreateLLMDTO) (*model.LLM, error) {
    if _, err := authz.RequireMember(ctx, tenantID); err != nil {
        return nil, err
    }
    ...
}
```

Point lookups (`FindByID`, `Update`, `DeleteByID`) additionally
fetch-then-compare, since the generic repository's `FindByID` is id-only,
not tenant-scoped:

```go
found, err := s.repo.FindByID(ctx, id)
if err != nil {
    return nil, err
}
if found.TenantID != tenantID {
    return nil, apperr.NotFound("llm not found") // 404, not 403 — don't
}                                                  // leak cross-tenant existence
```

Listing goes through `filter.Merge` instead:

```go
s.repo.Filter(ctx, filter.Merge(opts, filter.New[model.LLM](filter.WhereEq("tenant_id", tenantID))))
```

## API-key-authenticated routes have no tenant on the URL

A handful of routes exist for narrow runtime/integration use, not
tenant-staff administration, and stay on `guard.RequireApiKey` instead of
moving under `/tenants/{tenant_id}/...`: `POST /agent/{id}/run`,
`GET /api-key/me`, and the client-widget-facing surface
(`GET /clients/{code}/public`, `GET /clients/{code}/agents`, everything
under `/client/{client_code}/session/...`). These derive tenant
server-side instead of from a URL segment or `authz`:

- `agent.Service.Run` reads `auth.CurrentApiKey(ctx)` (only when present —
  it's also reachable via a client-widget JWT through
  `session.Service.Run`, which has no `ApiKeyIdentity` to check) and
  compares `agent.TenantID` against the calling key's `TenantID`.
- `client.Service.GetAgents` resolves the `Client` by code first, then
  scopes by `client.TenantID` — no `authz` call, since there's no
  `auth.UserIdentity` on this request to check membership against.
- `session`/`execution` rows carry a denormalized `TenantID` (copied from
  the owning `Client`/`Agent` at write time) for future querying, but reads
  stay scoped by `client_code` exactly as before — no new gate.

## Failure path

Guards run *outside* the `HandlerFunc`/`api.Serve` flow (see
[api.md](./api.md)) — on an auth failure they write the error response
themselves and never call `next`. `authz` failures are the opposite: they
happen *inside* a service method, so they flow back through the normal
`*apperr.Error` → `api.Serve` → HTTP-status path like any other domain
error (see [errors.md](./errors.md)).
