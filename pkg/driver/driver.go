// Package driver defines the common driver contract shared by the concrete
// driver kinds (llm, tool, knowledge source/index): metadata description via
// Info and a connectivity check via IDriver.TestConnection.
package driver

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
)

// Info describes a driver instance for display and configuration purposes:
// its identity (Key, Name, Description, Icon, Tags) and the JSON Schema
// (ConfigurationSchema) that validates the config map passed to the driver.
type Info struct {
	Key                 string       `json:"key,omitempty"`
	Name                string       `json:"name" validate:"required"`
	Description         string       `json:"description" validate:"required"`
	Icon                string       `json:"icon" validate:"omitempty"`
	Tags                []string     `json:"tags" validate:"omitempty"`
	ConfigurationSchema util.JSONMap `json:"configuration_schema"`
}

// IDriver is the base contract implemented by every driver kind. It exposes
// the driver's Info and a way to verify that a given config can connect
// successfully, without performing any domain-specific operation.
type IDriver interface {
	Info() Info
	TestConnection(ctx context.Context, config util.JSONMap) error
}
