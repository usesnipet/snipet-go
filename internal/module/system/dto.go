package system

type InfoDTO struct {
	Version string `json:"version"`
	// MultiTenantEnabled mirrors license.Info().Valid — lets the frontend
	// pick single- vs multi-tenant onboarding/registration copy. The actual
	// gate stays server-side (auth.Register, tenant.Create).
	MultiTenantEnabled bool `json:"multi_tenant_enabled"`
}
