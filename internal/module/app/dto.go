package app

type AppConfigDTO struct {
	InheritClient     bool   `json:"inherit_client"`
	InheritClientCode string `json:"inherit_client_code"`
	InheritClientName string `json:"inherit_client_name"`
}

type SystemInfoDTO struct {
	Version string `json:"version"`
	// MultiTenantEnabled mirrors license.Info().Valid — lets the frontend
	// pick single- vs multi-tenant onboarding/registration copy. The actual
	// gate stays server-side (auth.Register, tenant.Create).
	MultiTenantEnabled bool `json:"multi_tenant_enabled"`
}
