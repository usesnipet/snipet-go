package auth_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/auth"
	"github.com/usesnipet/snipet/internal/model"
)

func testAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		JWTSecret:     "test-secret-key-with-enough-length",
		JWTExpiration: time.Hour,
		JWTIssuer:     "https://test.snipet.cloud",
		JWTAudience:   "https://test.snipet.cloud",
	}
}

func testUser() *model.User {
	return &model.User{
		ID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name: "Test User",
	}
}

func TestGenerateTokenReturnsBearerToken(t *testing.T) {
	t.Parallel()

	service := auth.NewJWTService(testAuthConfig())
	token, _, err := service.GenerateToken("client-abc", testUser())
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(token, "Bearer "))
}

func TestGenerateTokenEmbedsExpectedClaims(t *testing.T) {
	t.Parallel()

	cfg := testAuthConfig()
	service := auth.NewJWTService(cfg)
	user := testUser()

	_, claims, err := service.GenerateToken("client-abc", user)
	require.NoError(t, err)

	assert.Equal(t, "client-abc", claims.ClientCode)
	assert.Equal(t, user.ID.String(), claims.Subject)
	assert.Equal(t, cfg.JWTIssuer, claims.Issuer)
	assert.Equal(t, cfg.JWTAudience, claims.Audience[0])
	assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now().Add(cfg.JWTExpiration), claims.ExpiresAt.Time, time.Second)
}

func TestVerifyTokenAcceptsValidToken(t *testing.T) {
	t.Parallel()

	service := auth.NewJWTService(testAuthConfig())
	user := testUser()

	token, expectedClaims, err := service.GenerateToken("client-abc", user)
	require.NoError(t, err)

	claims, err := service.VerifyToken(token)
	require.NoError(t, err)

	assert.Equal(t, expectedClaims.ClientCode, claims.ClientCode)
	assert.Equal(t, expectedClaims.Subject, claims.Subject)
	assert.Equal(t, expectedClaims.Issuer, claims.Issuer)
}

func TestVerifyTokenAcceptsTokenWithoutBearerPrefix(t *testing.T) {
	t.Parallel()

	service := auth.NewJWTService(testAuthConfig())
	token, _, err := service.GenerateToken("client-abc", testUser())
	require.NoError(t, err)

	claims, err := service.VerifyToken(strings.TrimPrefix(token, "Bearer "))
	require.NoError(t, err)

	assert.Equal(t, "client-abc", claims.ClientCode)
}

func TestVerifyTokenRejectsInvalidToken(t *testing.T) {
	t.Parallel()

	service := auth.NewJWTService(testAuthConfig())

	claims, err := service.VerifyToken("Bearer not.a.valid.token")
	require.Error(t, err)
	assert.Empty(t, claims.ClientCode)
}

func TestVerifyTokenRejectsWrongSecret(t *testing.T) {
	t.Parallel()

	issuer := auth.NewJWTService(testAuthConfig())
	token, _, err := issuer.GenerateToken("client-abc", testUser())
	require.NoError(t, err)

	otherSecret := testAuthConfig()
	otherSecret.JWTSecret = "another-secret-key-with-enough-length"
	verifier := auth.NewJWTService(otherSecret)

	claims, err := verifier.VerifyToken(token)
	require.Error(t, err)
	assert.Empty(t, claims.ClientCode)
}

func TestVerifyTokenRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	cfg := testAuthConfig()
	cfg.JWTExpiration = -time.Minute
	service := auth.NewJWTService(cfg)

	token, _, err := service.GenerateToken("client-abc", testUser())
	require.NoError(t, err)

	claims, err := service.VerifyToken(token)
	require.Error(t, err)
	assert.Empty(t, claims.ClientCode)
}
