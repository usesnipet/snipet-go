package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair_roundTrip(t *testing.T) {
	pubB64, privB64, err := GenerateKeyPair()
	require.NoError(t, err)

	pub, err := base64.StdEncoding.DecodeString(pubB64)
	require.NoError(t, err)
	require.Len(t, pub, 32)

	priv, err := decodePrivateKey(privB64)
	require.NoError(t, err)
	require.Equal(t, pub, []byte(priv[32:]))

	payload := []byte(`{"licensee":"Acme","issued_at":"2026-01-01","expires_at":"` + time.Now().Add(24*time.Hour).Format(time.DateOnly) + `"}`)
	key, err := Issue(privB64, payload)
	require.NoError(t, err)
	require.False(t, verify(key).Valid, "a freshly generated pair must not verify against the baked-in public key")

	parts := strings.SplitN(key, ".", 2)
	require.Len(t, parts, 2)
	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	require.NoError(t, err)
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	require.True(t, ed25519.Verify(pub, body, sig))
}

func TestIssue_validPayloadIsVerifiable(t *testing.T) {
	payload := []byte(`{
		"licensee": "Acme Inc",
		"issued_at": "2026-01-01",
		"expires_at": "` + time.Now().Add(24*time.Hour).Format(time.DateOnly) + `",
		"max_tenants": 5
	}`)

	key, err := Issue(testPrivateKeyBase64, payload)
	require.NoError(t, err)

	info := verify(key)
	require.True(t, info.Valid)
	require.Equal(t, 5, info.MaxTenants)
}

func TestIssue_rejectsBadPayloads(t *testing.T) {
	future := time.Now().Add(24 * time.Hour).Format(time.DateOnly)
	cases := []struct {
		name string
		raw  string
	}{
		{name: "empty licensee", raw: `{"licensee":"","issued_at":"2026-01-01","expires_at":"` + future + `"}`},
		{name: "bad expires_at", raw: `{"licensee":"Acme","issued_at":"2026-01-01","expires_at":"not-a-date"}`},
		{name: "expired", raw: `{"licensee":"Acme","issued_at":"2020-01-01","expires_at":"2021-01-01"}`},
		{name: "unknown field", raw: `{"licensee":"Acme","issued_at":"2026-01-01","expires_at":"` + future + `","extra":true}`},
		{name: "not json", raw: `nope`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Issue(testPrivateKeyBase64, []byte(tc.raw))
			require.Error(t, err)
		})
	}
}

func TestIssue_rejectsBadPrivateKey(t *testing.T) {
	payload := []byte(`{"licensee":"Acme","issued_at":"2026-01-01","expires_at":"` + time.Now().Add(24*time.Hour).Format(time.DateOnly) + `"}`)
	_, err := Issue("not-base64!!!", payload)
	require.Error(t, err)
}
