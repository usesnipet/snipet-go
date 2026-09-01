# `internal/auth` + `internal/guard`

Two layers, two different jobs:

- **`internal/auth`** — JWT signing/verification, API-key hashing, and the
  per-mechanism *identity* types + context helpers (below). Knows nothing
  about HTTP.
- **`internal/guard`** — composable `Gate`s that run those primitives
  against an incoming request and stash the resulting identity on the
  context.

The app is **single-tenant**: everything an authenticated admin can see
belongs to that one operator, so there is no membership/ownership layer on
top of the guards — a route is either open, gated by admin basic auth, or
gated by a per-App credential. A request authenticates via **exactly one**
mechanism; the identity types below are mutually exclusive on a given
request's context.

## Mechanisms + identity types

| Mechanism | Credential | Identity type | Set by |
|---|---|---|---|
| Admin basic auth | `Authorization: Basic` | `auth.BasicIdentity` (`Username`) | `guard.RequireBasicAuth` |
| API key | `X-API-Key` | `auth.ApiKeyIdentity` (`APIKeyID`) | `guard.RequireApiKey` |
| App key | `X-App-Key` | `auth.AppKeyIdentity` (`AppID`, `Code`) | `guard.RequireAppKey` |
| App end-user JWT | `Authorization: Bearer` | `auth.AppUserIdentity` (`UserID`, `AppCode`) | `guard.RequireAppUser` |

```go
type BasicIdentity struct {
    Username string
}

type ApiKeyIdentity struct {
    APIKeyID string
}

type AppKeyIdentity struct {
    AppID string
    Code  string
}

type AppUserIdentity struct {
    UserID  string
    AppCode string
}
```

Admin basic auth is the whole platform-staff story: a single
username/password pair from `config.AuthConfig`
(`AUTH_BASIC_AUTH_USERNAME` / `AUTH_BASIC_AUTH_PASSWORD`), compared in
constant time. There is no user table to look up.

## Reading who's authenticated

```go
identity, err := auth.CurrentBasic(ctx)   // auth.BasicIdentity
identity, err := auth.CurrentApiKey(ctx)  // auth.ApiKeyIdentity
identity, err := auth.CurrentAppKey(ctx)  // auth.AppKeyIdentity
identity, err := auth.CurrentAppUser(ctx) // auth.AppUserIdentity

auth.HasApiKey(ctx)  // bool, no error — useful in code paths reachable by
auth.HasAppKey(ctx)  // more than one mechanism (see session.Service.Run)
auth.HasAppUser(ctx)
```

Each getter errors if that mechanism didn't authenticate the request — a
service only calls the one matching whatever guard its route uses.

`AppKeyIdentity.Is(codeOrID)` and `AppUserIdentity.CanAccessApp(code)`
return `apperr.Forbidden(...)` when the credential is valid but points at a
different App than the route — the only per-request authorization check the
codebase still makes.

## Guards — gates + `api.Or`

Every guard constructor returns an `api.Gate`
(`func(*http.Request) (context.Context, error)`):

```go
requireBasicAuth := guard.RequireBasicAuth(cfg.Auth.BasicAuthUsername, cfg.Auth.BasicAuthPassword)
requireAPIKey    := guard.RequireApiKey(apiKeyService, apiKeyCache)
requireAppKey    := guard.RequireAppKey(appService, appKeyCache)
requireAppUser   := guard.RequireAppUser(appUserJWTService)

r.Use(requireBasicAuth.Handler())
r.Use(api.Or(requireAppKey, requireAppUser).Handler())
```

`api.Or` (in [`internal/api`](./api.md)) tries each gate in order — the
first one whose credential is present decides the outcome:

- credential **absent** → try the next gate (`auth.ErrNotApplicable`)
- credential **present but invalid** → hard 401 (does not fall through)
- credential **valid** → sets that identity on the context, done

All gates are built once in `bootstrap.Bootstrap` (see
[bootstrap.md](./bootstrap.md)) and handed to modules as raw `api.Gate`s —
a module's `handler.go` decides which gate(s) to `r.Use(...)` per route
group via `.Handler()` (see [modules.md](./modules.md#handlergo)). A single
module can mix gates across route groups when some routes need different
auth than others — `agent`'s CRUD routes use `requireBasicAuth` while its
`/run` route uses `requireApiKey`; `session`'s routes accept either an App
key or an App end-user JWT via `api.Or`.

## Per-App credentials don't need admin auth

The App-facing surface (`app`'s own read routes, `appuser`, `session`) is
authenticated by a credential scoped to a single App rather than by admin
basic auth:

- `RequireAppKey` verifies an `X-App-Key` against the `app` module and
  caches the result for a minute (same shape as `RequireApiKey`).
- `RequireAppUser` verifies a bearer JWT minted by `appauth` for an App
  end-user.
- Services on these routes call `identity.Is(...)` / `CanAccessApp(...)`
  to confirm the credential matches the App in the URL, then scope every
  query by that App — `session.Service.resolveApp` is the canonical
  example.

## Failure path

Guards run *outside* the `HandlerFunc`/`api.Serve` flow (see
[api.md](./api.md)) — on an auth failure they write the error response
themselves and never call `next`. A `Forbidden` returned from *inside* a
service method is the opposite: it flows back through the normal
`*apperr.Error` → `api.Serve` → HTTP-status path like any other domain
error (see [errors.md](./errors.md)).
