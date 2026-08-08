// Package driver defines the common driver contract shared by the concrete
// driver kinds (llm, tool, knowledge source/index): metadata description via
// Info and a connectivity check via IDriver.TestConnection.
package driver

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

var validate = validator.New()

// Info describes a driver instance for display and configuration purposes:
// its identity (Key, Name, Description, Icon, Tags) and the JSON Schema
// (ConfigurationSchema) that validates the config map passed to the driver.
// Key is the registry identity for the driver (see R.Register) and must be
// set by whoever builds the driver (e.g. via a CreateDriver's WithKey
// option); it is never derived or overwritten by the registry.
type Info struct {
	Key                 string        `json:"key" validate:"required"`
	Name                string        `json:"name" validate:"required"`
	Description         string        `json:"description" validate:"required"`
	Icon                string        `json:"icon" validate:"omitempty"`
	Tags                []string      `json:"tags" validate:"omitempty"`
	ConfigurationSchema jsonx.JSONMap `json:"configuration_schema"`
}

// Validate checks that Info's required fields (Key, Name, Description) are
// set. It only validates the metadata shape; it says nothing about whether
// the driver's behavior (e.g. its API funcs) is complete — see IDriver.Validate.
func (i Info) Validate() error {
	return validate.Struct(i)
}

// IDriver is the base contract implemented by every driver kind. It exposes
// the driver's Info, a way to verify that a given config can connect
// successfully, and Validate to check the driver itself is complete (valid
// Info plus every required behavior wired up) before it's trusted anywhere,
// most importantly before R.Register accepts it. A driver kind built via
// its CreateDriver already returns pre-validated; Validate exists so a
// hand-implemented driver (one that satisfies this interface without going
// through CreateDriver) can't silently skip the same checks.
type IDriver interface {
	Info() Info
	TestConnection(ctx context.Context, config jsonx.JSONMap) error
	Validate() error
}
