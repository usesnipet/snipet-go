package config

import "time"

type AuthConfig struct {
	JWTSecret     string        `env:"JWT_SECRET, default=change-me-in-production"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION, default=15m"`
	JWTIssuer     string        `env:"JWT_ISSUER, default=https://snipet.cloud"`
	JWTAudience   string        `env:"JWT_AUDIENCE, default=https://snipet.cloud"`

	RefreshTokenExpiration         time.Duration `env:"REFRESH_TOKEN_EXPIRATION, default=720h"`
	ActivateAccountTokenExpiration time.Duration `env:"ACTIVATE_ACCOUNT_TOKEN_EXPIRATION, default=24h"`
	ResetPasswordTokenExpiration   time.Duration `env:"RESET_PASSWORD_TOKEN_EXPIRATION, default=1h"`
	TenantInvitationExpiration     time.Duration `env:"TENANT_INVITATION_EXPIRATION, default=168h"`

	AppURL string `env:"APP_URL, default=http://localhost:5173"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL"`

	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURL  string `env:"GITHUB_REDIRECT_URL"`
}
