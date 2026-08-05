package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/usesnipet/snipet/config"
)

const refreshTokenBytes = 32

type RefreshTokenService struct {
	expiration time.Duration
}

func NewRefreshTokenService(authConfig config.AuthConfig) *RefreshTokenService {
	return &RefreshTokenService{expiration: authConfig.RefreshTokenExpiration}
}

// GenerateRefreshToken creates a cryptographically secure opaque refresh token.
func (s *RefreshTokenService) GenerateRefreshToken() (string, time.Time, error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, err
	}
	return base64.RawURLEncoding.EncodeToString(b), time.Now().Add(s.expiration), nil
}

// HashRefreshToken returns the SHA-256 hex digest of a refresh token for storage/lookup.
func (s *RefreshTokenService) HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
