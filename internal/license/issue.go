package license

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GenerateKeyPair returns a new Ed25519 pair as Std-base64 strings (the same
// encoding as publicKeyBase64 / the private-key files the CLI writes).
func GenerateKeyPair() (pubB64, privB64 string, err error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pub), base64.StdEncoding.EncodeToString(priv), nil
}

// Issue validates payloadJSON, signs it with the Std-base64 private key, and
// returns a LICENSE_KEY string the runtime verifier accepts.
func Issue(privateKeyB64 string, payloadJSON []byte) (string, error) {
	p, err := parsePayload(payloadJSON)
	if err != nil {
		return "", err
	}
	priv, err := decodePrivateKey(privateKeyB64)
	if err != nil {
		return "", err
	}
	return sign(priv, p)
}

func parsePayload(raw []byte) (payload, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var p payload
	if err := dec.Decode(&p); err != nil {
		return payload{}, fmt.Errorf("invalid payload: %w", err)
	}
	if dec.More() {
		return payload{}, fmt.Errorf("invalid payload: trailing data")
	}
	if strings.TrimSpace(p.Licensee) == "" {
		return payload{}, fmt.Errorf("invalid payload: licensee is required")
	}
	if _, err := time.Parse(time.DateOnly, p.IssuedAt); err != nil {
		return payload{}, fmt.Errorf("invalid payload: issued_at must be YYYY-MM-DD")
	}
	expiresAt, err := time.Parse(time.DateOnly, p.ExpiresAt)
	if err != nil {
		return payload{}, fmt.Errorf("invalid payload: expires_at must be YYYY-MM-DD")
	}
	if !time.Now().Before(expiresAt) {
		return payload{}, fmt.Errorf("invalid payload: expires_at is not in the future")
	}
	return p, nil
}

func decodePrivateKey(encoded string) (ed25519.PrivateKey, error) {
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	if len(raw) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key: want %d bytes, got %d", ed25519.PrivateKeySize, len(raw))
	}
	return ed25519.PrivateKey(raw), nil
}

func sign(priv ed25519.PrivateKey, p payload) (string, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	sig := ed25519.Sign(priv, body)
	return encodeKey(body, sig), nil
}

func encodeKey(payloadBytes, signature []byte) string {
	return base64.RawURLEncoding.EncodeToString(payloadBytes) + "." + base64.RawURLEncoding.EncodeToString(signature)
}
