package config

import "time"

type AuthConfig struct {
	JWTSecret     string        `env:"JWT_SECRET, default=change-me-in-production"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION, default=15m"`
	JWTIssuer     string        `env:"JWT_ISSUER, default=https://snipet.cloud"`
	JWTAudience   string        `env:"JWT_AUDIENCE, default=https://snipet.cloud"`
}
