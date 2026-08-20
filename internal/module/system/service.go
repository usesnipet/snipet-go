package system

import (
	"github.com/usesnipet/snipet/internal/license"
	"github.com/usesnipet/snipet/version"
)

type Service struct {
	license *license.Service
}

func NewService(license *license.Service) *Service {
	return &Service{license: license}
}

// Info exposes MultiTenantEnabled so the frontend can branch its
// onboarding UI (single vs multi tenant) without duplicating the license
// check — enforcement itself still happens server-side in auth.Register and
// tenant.Create.
func (s *Service) Info() *InfoDTO {
	return &InfoDTO{
		Version:            version.Version,
		MultiTenantEnabled: s.license.Info().Valid,
	}
}
