package app

import (
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/license"
	"github.com/usesnipet/snipet/version"
)

type Service struct {
	cfg     *config.AppConfig
	license *license.Service
}

func NewService(cfg *config.AppConfig, license *license.Service) *Service {
	return &Service{cfg: cfg, license: license}
}

func (s *Service) Config() *AppConfigDTO {
	return &AppConfigDTO{
		InheritClient:     s.cfg.InheritClient,
		InheritClientCode: s.cfg.InheritClientCode,
		InheritClientName: s.cfg.InheritClientName,
	}
}

// SystemInfo exposes MultiTenantEnabled so the frontend can branch its
// onboarding UI (single vs multi tenant) without duplicating the license
// check — enforcement itself still happens server-side in auth.Register and
// tenant.Create.
func (s *Service) SystemInfo() *SystemInfoDTO {
	return &SystemInfoDTO{
		Version:            version.Version,
		MultiTenantEnabled: s.license.Info().Valid,
	}
}
