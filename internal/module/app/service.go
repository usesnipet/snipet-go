package app

import (
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/version"
)

type Service struct {
	cfg *config.AppConfig
}

func NewService(cfg *config.AppConfig) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Config() *AppConfigDTO {
	return &AppConfigDTO{
		InheritClient:     s.cfg.InheritClient,
		InheritClientCode: s.cfg.InheritClientCode,
		InheritClientName: s.cfg.InheritClientName,
	}
}

func (s *Service) SystemInfo() *SystemInfoDTO {
	return &SystemInfoDTO{
		Version: version.Version,
	}
}
