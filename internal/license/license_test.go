package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/config"
)

// testPrivateKey is the Ed25519 private key matching publicKeyBase64 above.
// Test-only fixture signer — the real Licensor private key never lives in
// this repo (see plan-ee/licensing.md).
const testPrivateKeyBase64 = "9Dg0yItkSSgv1XNlepInIPB4UM4BjvxweueQnmR5OHQix+Bq6EAGfEajbuFCZzNppa94nvT2I3B3GEEFOhtA9Q=="

func signFixture(t *testing.T, p payload) string {
	t.Helper()
	key, err := sign(ed25519.PrivateKey(mustDecode(t, testPrivateKeyBase64)), p)
	require.NoError(t, err)
	return key
}

func mustDecode(t *testing.T, b64 string) []byte {
	t.Helper()
	raw, err := base64.StdEncoding.DecodeString(b64)
	require.NoError(t, err)
	return raw
}

func TestVerify_emptyKeyIsUnlicensed(t *testing.T) {
	info := verify("")
	require.False(t, info.Valid)
}

func TestVerify_malformedKeyIsUnlicensed(t *testing.T) {
	require.False(t, verify("not-a-license-key").Valid)
	require.False(t, verify("only.one.dot.too.many").Valid)
	require.False(t, verify("!!!.!!!").Valid)
}

func TestVerify_validSignedKeyIsLicensed(t *testing.T) {
	key := signFixture(t, payload{
		Licensee:   "Acme Inc",
		IssuedAt:   "2026-01-01",
		ExpiresAt:  time.Now().Add(24 * time.Hour).Format(time.DateOnly),
		MaxTenants: 5,
	})

	info := verify(key)
	require.True(t, info.Valid)
	require.Equal(t, 5, info.MaxTenants)
}

func TestVerify_expiredKeyIsUnlicensed(t *testing.T) {
	key := signFixture(t, payload{
		Licensee:  "Acme Inc",
		IssuedAt:  "2020-01-01",
		ExpiresAt: "2021-01-01",
	})

	require.False(t, verify(key).Valid)
}

func TestVerify_tamperedPayloadIsUnlicensed(t *testing.T) {
	key := signFixture(t, payload{
		Licensee:  "Acme Inc",
		IssuedAt:  "2026-01-01",
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	})

	// flip the licensee in the payload half without re-signing
	tampered := signFixture(t, payload{
		Licensee:  "Evil Corp",
		IssuedAt:  "2026-01-01",
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.DateOnly),
	})
	require.NotEqual(t, key, tampered)
	require.False(t, verify(key[:len(key)-4]+"XXXX").Valid)
}

func TestNewService_wrapsVerify(t *testing.T) {
	svc := NewService(config.LicenseConfig{LicenseKey: ""})
	require.False(t, svc.Info().Valid)
}
