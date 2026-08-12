# Middleware (tenant staff)

Not CRUD, no entity — two `api.MiddlewareFunc`s the tenant-staff routes (`auth`, `user`, `tenant`, `member`, `tenant-invitation`) compose in. They live in the existing `internal/middleware` package, alongside `JWT`/`AnyAuth` — same file layout as any other middleware, no separate directory (there's no `ee` split, see `plan-ee/boundary.md`).

**Wiring note:** since there's no build-tag boundary to respect anymore, `Bearer`/`RequirePlatformAdmin` are constructed in `internal/bootstrap/bootstrap.go` exactly like `JWT`/`AnyAuth` are today — no indirection needed.

## Depends on

- `*auth.JWTService[*platformauth.UserClaims]` (see `plan-ee/module/auth.md`'s "Token issuance" section) — `internal/auth` (the generic base, `Principal`/`GetPrincipal`/`JWTService`) is imported plainly as `auth`; `internal/module/auth` (this doc's consumer, tenant staff) is aliased `platformauth` to avoid the collision, same convention as `plan-ee/module/auth.md`'s `coreauth` alias used the other way around from inside that module.
- The generic `internal/auth.Principal`/`SetPrincipal`/`GetPrincipal` — **reused as-is, not duplicated.** Per `auth.md`'s JWT genericization, `Principal.JWTClaims` becomes `jwt.Claims` (interface) instead of the concrete client `*UserClaims`, specifically so both the client-widget `JWTService` and this `JWTService[*platformauth.UserClaims]` can populate the same context mechanism. No separate `Principal`/context package needed — that would just be the same 20 lines twice for no isolation benefit, since setting/reading a context value is stateless either way.
- `repository.IUserRepository` — only for `RequirePlatformAdmin`.

## `Bearer` — authenticate the tenant staff user's JWT

```go
func Bearer(jwtService *auth.JWTService[*platformauth.UserClaims]) api.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if !strings.HasPrefix(token, "Bearer ") {
                api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
                return
            }
            claims, err := jwtService.VerifyToken(token)
            if err != nil {
                api.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
                return
            }
            principal := auth.NewPrincipal(auth.PrincipalTypeJWT, nil, claims)
            ctx := auth.SetPrincipal(r.Context(), principal)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

Same shape as `internal/middleware.JWT`/`jwtAuth`, just parameterized on the `platformauth.JWTService` instance instead of the client one — no API-key branch (tenant staff has no API-key concept today, unlike `AnyAuth`).

**Reading the current user id** — every tenant-staff service that needs "who's calling" (`user.Me`, `tenant.Create`, `tenant.Leave`, `member.UpdateRole`'s caller check, etc.) does the same three lines, so it's worth one small helper instead of repeating them. Lives in `internal/auth` itself (not this module) since it's generic enough for any `Principal`, not tenant-staff-specific:

```go
// internal/auth/context.go
func CurrentUserID(ctx context.Context) (string, error) {
    principal, ok := GetPrincipal(ctx)
    if !ok || principal.GetType() != PrincipalTypeJWT {
        return "", apperr.Unauthorized("unauthorized")
    }
    return principal.GetJWTClaims().GetSubject() // jwt.Claims.GetSubject() — no cast to *platformauth.UserClaims needed for just the id
}
```

`jwt.Claims.GetSubject()` is part of the `jwt-go` interface itself (`GetSubject() (string, error)`), so reading just the user id needs no type assertion back to `*platformauth.UserClaims` — only code that needs a tenant-staff-specific claim (none exist yet, see `auth.md`) would assert `principal.GetJWTClaims().(*platformauth.UserClaims)`.

## `RequirePlatformAdmin` — require `user.is_admin`

```go
func RequirePlatformAdmin(userRepo repository.IUserRepository) api.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID, err := auth.CurrentUserID(r.Context())
            if err != nil {
                api.WriteError(w, http.StatusUnauthorized, err)
                return
            }
            user, err := userRepo.FindByID(r.Context(), userID)
            if err != nil {
                api.WriteError(w, http.StatusUnauthorized, err)
                return
            }
            if !user.IsAdmin {
                api.WriteError(w, http.StatusForbidden, errors.New("forbidden: platform admin required"))
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Must run after `Bearer`** in the chain — it reads the principal `Bearer` sets, it doesn't authenticate the token itself. `IsAdmin` isn't in the JWT claims (kept minimal per `auth.md`), so this is a DB hit on every request to a platform-admin route — same trade-off `apiKeyAuth` already makes for API keys (there it's mitigated with a short-TTL cache; not bothering with that here since platform-admin routes are low-traffic — revisit if that stops being true).

Usage (e.g. `tenant.md`'s `GET /tenants`):

```go
r.Route("/tenants", func(r chi.Router) {
    r.Group(func(r chi.Router) {
        r.Use(mw.Bearer(jwtService))
        r.Use(mw.RequirePlatformAdmin(userRepo))
        r.Get("/", serve(h.filter))
    })
    // ...other /tenants routes, just mw.Bearer(jwtService), no admin requirement...
})
```

## Not covered here

The tenant-scoped "is this user admin **of this tenant** (or platform admin)" check (`authz.CanManageTenant`, see `plan-ee/module/tenant.md`) stays a service-level call, not a third middleware — it needs the `{tenant_id}` route param *and* feeds into methods that also do other things with the result (e.g. `tenant.DeleteByID` needs the boolean, not just a 403 short-circuit), so it's colocated with the business logic that uses it rather than extracted into the request pipeline.
