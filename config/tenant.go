package config

type TenantConfig struct {
	TenantName string `env:"SINGLE_TENANT_NAME, default=Snipet"` // bootstrap tenant's name
	TenantSlug string `env:"SINGLE_TENANT_SLUG, default=default"` // bootstrap tenant's slug
}
