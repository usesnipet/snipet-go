# `internal/auth` + `internal/middleware`

Three auth mechanisms. A request can carry **more than one** at once
(e.g. `Or(RequireClientJWT, RequireAPIKey)` with both headers present) —
every method that succeeds is stored in a principals list on the context.

- **`internal/auth`** — JWT signing/verification, claim types, API key
  hashing/generation, `Principal` / `Principals`, context helpers.
- **`internal/middleware`** — composable gates that run those primitives
  and append successful authentications to the list.

## Mechanisms + claim types

Both JWT flavours live in `internal/auth` (no type parameter at call
sites):

| Mechanism | Credential | Principal type | Claims |
|---|---|---|---|
| Client JWT | `Authorization: Bearer` | `PrincipalTypeClientJWT` | `auth.ClientUserClaims` (`ClientCode`) |
| Platform JWT | `Authorization: Bearer` | `PrincipalTypePlatformJWT` | `auth.PlatformUserClaims` |
| API key | `X-API-Key` | `PrincipalTypeAPIKey` | — (`APIKeyID`) |

```go
type Principal struct {
    Type      PrincipalType
    APIKeyID  *string
    JWTClaims jwt.Claims
}

type Principals []*Principal
```

## Reading who's authenticated

```go
auth.HasPrincipal(ctx, auth.PrincipalTypeAPIKey)

id, err := auth.APIKeyID(ctx)
claims, err := auth.ClientJWTClaims(ctx)   // *auth.ClientUserClaims
claims, err := auth.PlatformJWTClaims(ctx) // *auth.PlatformUserClaims
uid, err := auth.ClientUserID(ctx)
uid, err := auth.PlatformUserID(ctx)
```

Helpers look up the matching entry in the principals list — they do not
require a type argument.

## Middleware — require gates + `Or`

```go
requireAPIKey := middleware.RequireAPIKey(apiKeyService, apiKeyCache)
requireClientJWT := middleware.RequireClientJWT(clientJWTService)
requirePlatformJWT := middleware.RequirePlatformJWT(platformJWTService)
anyClientAuth := middleware.Or(requireClientJWT, requireAPIKey)

r.Use(requireAPIKey.Handler())
r.Use(requirePlatformJWT.Handler())
r.Use(anyClientAuth.Handler())
```

`Or` runs **every** gate:

- credential **absent** → skip (`auth.ErrNotApplicable`)
- credential **present but invalid** → hard 401 (does not fall through)
- credential **valid** → append that `Principal` to the list

At least one success is required. So a request with both a valid client
JWT and a valid API key ends up with two principals on the context.

All gates are built once in `bootstrap.Bootstrap` (see
[bootstrap.md](./bootstrap.md)) and handed to modules as
`api.MiddlewareFunc` via `.Handler()` — a module's `handler.go` decides
which middleware(s) to `r.Use(...)` per route group (see
[modules.md](./modules.md#handlergo)).

## Failure path

Middleware runs *outside* the `HandlerFunc`/`api.Serve` flow (see
[api.md](./api.md)) — on an auth failure it writes the error response
itself (`api.WriteError(w, http.StatusUnauthorized, ...)`) and never
calls `next`.
