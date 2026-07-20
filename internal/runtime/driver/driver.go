package driver

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

type Info struct {
	Name                string       `json:"name" validate:"required"`
	Description         string       `json:"description" validate:"required"`
	Icon                string       `json:"icon" validate:"omitempty"`
	Tags                []string     `json:"tags" validate:"omitempty"`
	ConfigurationSchema util.JSONMap `json:"configuration_schema"`
}

type IDriver interface {
	Info() Info
	TestConnection(ctx context.Context, config util.JSONMap) error
}
