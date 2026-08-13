// Package license verifies offline, Ed25519-signed license keys that gate
// multi-tenant use. See plan-ee/licensing.md.
package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/usesnipet/snipet/config"
)

// publicKey verifies keys signed by the Licensor's private key. Not
// configurable — a self-hoster cannot point verification at their own key.
const publicKeyBase64 = "IsfgauhABnxGo27hQmczaaWveJ709iNwdxhBBTobQPU="

type payload struct {
	Licensee   string `json:"licensee"`
	IssuedAt   string `json:"issued_at"`
	ExpiresAt  string `json:"expires_at"`
	MaxTenants int    `json:"max_tenants"`
}

type Info struct {
	Valid      bool
	MaxTenants int // 0 = unlimited
}

type Service struct {
	info Info
}

// NewService verifies cfg.LicenseKey once and caches the result for the
// process lifetime — no per-request re-verification.
func NewService(cfg config.LicenseConfig) *Service {
	return &Service{info: verify(cfg.LicenseKey)}
}

func (s *Service) Info() Info {
	return s.info
}

// verify never errors — an empty, malformed, or expired key just means
// "behave as unlicensed"; it must never crash boot, since a self-hosted
// community instance has no key at all by design.
func verify(key string) Info {
	unlicensed := Info{Valid: false}
	if key == "" {
		return unlicensed
	}

	// Key = base64url(payload_json) + "." + base64url(ed25519_signature).
	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return unlicensed
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return unlicensed
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return unlicensed
	}

	pubKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil || len(pubKey) != ed25519.PublicKeySize {
		return unlicensed
	}
	if !ed25519.Verify(pubKey, payloadBytes, signature) {
		return unlicensed
	}

	var p payload
	if err := json.Unmarshal(payloadBytes, &p); err != nil {
		return unlicensed
	}

	expiresAt, err := time.Parse(time.DateOnly, p.ExpiresAt)
	if err != nil || time.Now().After(expiresAt) {
		return unlicensed
	}

	return Info{Valid: true, MaxTenants: p.MaxTenants}
}
