package manager

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/pkg/driver"
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// DriverManager is a thin facade over a driver.Registry that adds the
// runtime-facing concerns the bare registry doesn't handle: validating a
// config map against a driver's schema, and dialing a driver to confirm the
// config actually connects (see Connect).
type DriverManager[T driver.IDriver] struct {
	registry *driver.Registry[T]
}

func NewDriverManager[T driver.IDriver](registry *driver.Registry[T]) *DriverManager[T] {
	return &DriverManager[T]{registry: registry}
}

func (m *DriverManager[T]) GetDriver(key string) (T, error) {
	driverInstance, ok := m.registry.Get(key)
	if !ok {
		return driverInstance, driver.ErrDriverNotFound
	}
	return driverInstance, nil
}

// Names returns the sorted keys of every registered driver.
func (m *DriverManager[T]) Names() []string {
	return m.registry.Names()
}

// ValidateConfiguration checks config against a driver's JSON Schema.
func (m *DriverManager[T]) ValidateConfiguration(schema jsonx.JSONMap, config jsonx.JSONMap) error {
	return jsonschema.Validate(schema, config)
}

// ValidateConfigurationByKey looks up the driver by key and validates config
// against its schema.
func (m *DriverManager[T]) ValidateConfigurationByKey(key string, config jsonx.JSONMap) error {
	driverInstance, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfiguration(driverInstance.Info().ConfigurationSchema, config)
}

// Connect resolves the driver by key, validates config against its schema and
// runs its connectivity check, returning the ready-to-use driver instance.
func (m *DriverManager[T]) Connect(ctx context.Context, driverKey string, config jsonx.JSONMap) (T, error) {
	driverInstance, err := m.GetDriver(driverKey)
	if err != nil {
		return driverInstance, err
	}
	if err := m.ValidateConfiguration(driverInstance.Info().ConfigurationSchema, config); err != nil {
		return driverInstance, err
	}
	if err := driverInstance.TestConnection(ctx, config); err != nil {
		return driverInstance, fmt.Errorf("%w: %v", driver.ErrDriverConnectionFailed, err)
	}
	return driverInstance, nil
}

// ListDrivers returns the Info of every registered driver, sorted by key.
func (m *DriverManager[T]) ListDrivers(ctx context.Context) ([]driver.Info, error) {
	names := m.registry.Names()
	drivers := make([]driver.Info, 0, len(names))

	for _, name := range names {
		driverInstance, err := m.GetDriver(name)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driverInstance.Info())
	}

	return drivers, nil
}
