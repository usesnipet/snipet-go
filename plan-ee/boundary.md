# No code boundary — single codebase, license-key gated

## What this doc used to say

This doc originally proposed a compile-time split: `internal/ee/**` for tenant/user/auth/member/etc., built into the binary only via `go build -tags ee`, kept out of the default "community" build entirely. That assumption doesn't hold: `Tenant`, `User`, `Account`, `Token`, `Member`, `TenantInvitation` are needed by every deployment, licensed or not — a single self-hosted instance still needs a `User` to log in, a `Token` for refresh, exactly one `Tenant` to exist. There is no code to exclude from a "community" build; the only thing that differs between a licensed and unlicensed instance is whether a *second* `Tenant` may be created. That's a usage limit, not a feature set, so it can't be enforced by choosing which files get compiled.

## Current model

- **No `internal/ee/**` directory, no build tags, no `registerEE` hook.** Every module — `tenant`, `user`, `auth`, `account`, `token`, `member`, `tenant-invitation`, plus their middleware — lives under `internal/module/**`, `internal/repository/**`, `internal/middleware/**`, exactly like every other module in this codebase (`agent`, `api-key`, `client`, `knowledge`, etc.). One binary, one build, no flags.
- **What actually gates single vs. multi tenant is a runtime license key**, checked inside `tenant.Service.Create` when a second `Tenant` would be created. See `plan-ee/licensing.md` for the mechanism. `LICENSE.md` §5 ("Multi-Tenant Use") is the legal backing for this — the license restricts the *right to operate* the capability, independent of where the code lives.

## Naming — freeing up `user`/`auth` for the new modules

Dropping the `internal/ee/**` prefix removes the free disambiguator the earlier plan relied on. `internal/module/user` and `internal/module/auth` are already taken today by the client-widget flow (a visitor authenticating into a `Client` via OIDC/anonymous) — a different actor from the new tenant-staff `User`/`auth`. Decision: rename the existing client-widget code first, as its own prerequisite commit, before building anything from this plan:

| Today | Renamed to |
|---|---|
| `internal/module/user` | `internal/module/clientuser` |
| `internal/module/auth` | `internal/module/clientauth` |
| `model.User` (client-widget) | `model.ClientUser` |

This frees `internal/module/user`, `internal/module/auth`, and `model.User` for the tenant-staff entities documented in `plan-ee/module/user.md` and `plan-ee/module/auth.md`. `internal/auth` (the low-level package — `Principal`, `JWTService[T]`, `BaseClaims`, `CurrentUserID`) is unaffected by the rename; it's already generic and shared by both, imported as `coreauth` from within `clientauth`/`auth` where a plain `auth` import would otherwise collide with `internal/auth` itself.

## Where each module lives now

No indirection to wire up — `internal/bootstrap/bootstrap.go` constructs and mounts `tenant`, `user`, `auth`, `account`, `token`, `member`, `tenant-invitation` the same way it already constructs every other module's handler, in the same unconditional code path.
