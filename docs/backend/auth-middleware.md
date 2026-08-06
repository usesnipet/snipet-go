# `internal/auth` + `internal/middleware`

Two auth mechanisms, both landing on the same `auth.Principal` abstraction
so a handler doesn't need to know which one authenticated the caller.

- **`internal/auth`** — the mechanism-specific primitives (JWT
  signing/verification, API key hashing/generation, `Principal`).
- **`internal/middleware`** — chi middleware that runs those primitives
  against a request and puts a `Principal` on its context.

## The two mechanisms

- **JWT** (`internal/auth/jwt.go`) — short-lived user sessions.
  `JWTService.GenerateToken` signs a `UserClaims` (embeds
  `jwt.RegisteredClaims` + `ClientCode`) with HS256; `VerifyToken` parses
  and validates a `"Bearer <token>"` header. Paired with
  `internal/auth/refresh_token.go` for refreshing an expired access token
  without re-authenticating.
- **API key** (`internal/auth/apikey-generator.go` +
  `internal/module/api-key`) — long-lived machine credentials. Keys are
  generated once, shown to the user once, and stored **hashed**
  (`internal/auth/hasher.go`) — verifying a request means hashing the
  presented key and comparing, never storing/comparing plaintext.

## `Principal` — the common result of authenticating

```go
type PrincipalType string
const (
    PrincipalTypeAPIKey PrincipalType = "api_key"
    PrincipalTypeJWT    PrincipalType = "jwt"
)

type Principal struct {
    Type      PrincipalType
    APIKeyID  *string
    JWTClaims *UserClaims
}
```

Stashed on the request context (`auth.SetPrincipal`/`auth.GetPrincipal`,
`context.go`) by whichever middleware ran. A handler/service that needs to
know who's calling reads `auth.GetPrincipal(ctx)` and checks `Type` rather
than assuming which auth mode protected the route.

## Middleware

```go
apiKeyMiddleware := middleware.APIKeyMiddleware(apiKeyService, apiKeyCache)
anyAuthMiddleware := middleware.AnyAuth(jwtService, apiKeyService, apiKeyCache)
jwtMiddleware := middleware.JWT(jwtService)
```

- **`APIKeyMiddleware`** — requires `X-API-Key`. Checks a short-TTL
  in-memory cache first (`internal/infra/cache`, see
  [infra.md](./infra.md)) before hitting `apiKeyService.VerifyAPIKey` (a
  hash + DB lookup), and populates the cache on a hit — every API-key
  protected request would otherwise hash+query on every call.
- **`JWT`** — requires an `Authorization: Bearer <token>` header.
- **`AnyAuth`** — accepts either header, preferring JWT when both are
  present; used by routes reachable from both the admin (JWT) and
  programmatic (API key) surfaces (e.g. `client`, `session`, `user`).

All three are built once in `bootstrap.Bootstrap` (see
[bootstrap.md](./bootstrap.md)) and handed to modules as
`api.MiddlewareFunc` — a module's `handler.go` decides which middleware(s)
to `r.Use(...)` per route group, it doesn't build its own (see
[modules.md](./modules.md#handlergo)).

## Failure path

Middleware runs *outside* the `HandlerFunc`/`api.Serve` flow (see
[api.md](./api.md)) — on an auth failure it has to write the error response
itself (`api.WriteError(w, http.StatusUnauthorized, ...)`) and never call
`next`, rather than returning an error for `Serve` to translate.
