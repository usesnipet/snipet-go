# Licensing — single vs multi tenant via license key

Legal backing: `LICENSE.md` §5 ("Multi-Tenant Use — Licensed Capability"). Architectural backing: `plan-ee/boundary.md` — there is no separate enterprise codebase or build; every instance runs the same code, and this doc's mechanism is the only thing that decides whether a second `Tenant` may be created.

## Mechanism — offline-verifiable license key

No phone-home requirement — self-hosted instances may run fully air-gapped. Verification is a local signature check at boot, not a network call.

- **Key format:** base64 of `{payload_json}.{signature}`, where `payload_json` is:
  ```json
  {"licensee": "Acme Inc", "issued_at": "2026-01-01", "expires_at": "2027-01-01", "max_tenants": 0}
  ```
  `max_tenants: 0` (or absent) means unlimited; a positive number caps it below unlimited, if a tiered pricing plan is ever sold. `signature` is an Ed25519 signature over `payload_json`, produced by the Licensor's private key at license-issuance time.
- **Verification key:** the corresponding Ed25519 public key is baked into the binary as a `const` in `internal/license` — not configurable, so a self-hoster cannot point verification at their own key.
- **Config:**
  ```go
  // config.LicenseConfig
  type LicenseConfig struct {
      LicenseKey string `env:"LICENSE_KEY"` // empty => unlicensed, single-tenant only
  }
  ```
- **Validation:** parsed and signature-checked once at boot by `internal/license.Service`, fail-soft (an empty, invalid, or expired key just means "behave as unlicensed" — it must never crash boot, since a self-hosted community instance has no key at all by design). Result cached in memory for the process lifetime:
  ```go
  // internal/license/license.go
  type Info struct {
      Valid      bool
      MaxTenants int // 0 = unlimited
  }

  type Service struct {
      info Info // computed once in NewService, from config.LicenseConfig
  }

  func NewService(cfg config.LicenseConfig) *Service
  func (s *Service) Info() Info
  ```
  No per-request re-verification, no filesystem/network access after boot.

## Enforcement point

`tenant.Service.Create` (`plan-ee/module/tenant.md`) is the only place this is checked — the single spot that already handles creating a `Tenant`:

```go
func (s *Service) Create(ctx context.Context, dto CreateTenantDTO) (*model.Tenant, error) {
    lic := s.license.Info()
    count := /* current tenant count, via tenantRepo.Filter + .Total */
    if !lic.Valid && count > 0 {
        return nil, apperr.Forbidden("multi-tenant use requires a Snipet Enterprise License")
    }
    if lic.Valid && lic.MaxTenants > 0 && count >= lic.MaxTenants {
        return nil, apperr.Forbidden("tenant limit reached for this license")
    }
    // unlicensed + zero tenants yet: allowed — this is the bootstrap single tenant
    // ...existing create logic
}
```

`tenant.Service.Init` (bootstrap of the first `Tenant`, run from `internal/bootstrap`) is unaffected by license state — it always runs, licensed or not, since every instance needs exactly that one `Tenant` to exist.

## Why this is enough, and what it doesn't try to do

Self-hosters have source access (`LICENSE.md` §2b grants copy/modify rights) and can delete the check above from a local build. This is expected and explicitly addressed in `LICENSE.md` §5: doing so does not grant the right to run Multi-Tenant Use, it just makes running it unlicensed a breach of contract instead of a technical impossibility. This mirrors how other source-available, commercially-restricted projects (n8n, Chatwoot, Cal.com, etc.) handle the same tension — protection is primarily legal (§5 + §7 termination-on-breach), the runtime check is friction for the common case, not DRM. Do not invest in obfuscation, anti-tamper, or telemetry-based enforcement here — that fight isn't worth winning technically when the license already wins it contractually.

## Open questions (not blocking single-tenant implementation)

- License issuance/signing tooling (a small internal CLI to mint keys) — not designed yet, needed before the first paying multi-tenant customer, not before.
- Whether `max_tenants` tiers are ever actually sold, or licenses are simply valid/unlimited vs absent/single — affects whether `MaxTenants` is worth keeping in the payload now or added later (adding a field to a signed payload later doesn't invalidate already-issued keys, so safe to defer).
