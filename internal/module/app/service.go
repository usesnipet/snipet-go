package app

import "github.com/usesnipet/snipet/config"

type Service struct {
	cfg *config.AppConfig
}

func NewService(cfg *config.AppConfig) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Config() *AppConfigDTO {
	return &AppConfigDTO{
		InheritClient: s.cfg.InheritClient,
	}
}
