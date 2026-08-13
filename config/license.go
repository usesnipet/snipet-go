package config

type LicenseConfig struct {
	LicenseKey string `env:"LICENSE_KEY"` // empty => unlicensed, single-tenant only
}
