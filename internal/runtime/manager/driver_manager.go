package manager

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/pkg/driver"
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

type Driver[T driver.IDriver] struct {
	registry *registry.R[T]
}

func NewDriver[T driver.IDriver](registry *registry.R[T]) *Driver[T] {
	return &Driver[T]{registry: registry}
}

func (m *Driver[T]) GetDriver(key string) (T, error) {
	driverInstance, ok := m.registry.Get(key)
	if !ok {
		return driverInstance, driver.ErrDriverNotFound
	}
	return driverInstance, nil
}

// Names returns the sorted keys of every registered driver.
func (m *Driver[T]) Names() []string {
	return m.registry.Names()
}

func (m *Driver[T]) ValidateConfiguration(schema jsonx.JSONMap, config jsonx.JSONMap) error {
	return jsonschema.Validate(schema, config)
}

func (m *Driver[T]) ValidateConfigurations(schema jsonx.JSONMap, configs ...jsonx.JSONMap) error {
	for _, config := range configs {
		if err := jsonschema.Validate(schema, config); err != nil {
			return err
		}
	}
	return nil
}

func (m *Driver[T]) ValidateConfigurationByKey(key string, config jsonx.JSONMap) error {
	driver, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfiguration(driver.Info().ConfigurationSchema, config)
}

func (m *Driver[T]) ValidateConfigurationsByKey(key string, configs ...jsonx.JSONMap) error {
	driver, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfigurations(driver.Info().ConfigurationSchema, configs...)
}

func (m *Driver[T]) ValidateMultipleConfigurationsByKey(configs ...Configuration) error {
	for _, c := range configs {
		if err := m.ValidateConfigurationByKey(c.Key, c.Config); err != nil {
			return err
		}
	}
	return nil
}

func (m *Driver[T]) Prepare(ctx context.Context, driverKey string, config jsonx.JSONMap) (T, error) {
	driverInstance, err := m.GetDriver(driverKey)
	if err != nil {
		return driverInstance, err
	}
	schema := driverInstance.Info().ConfigurationSchema
	if err := m.ValidateConfiguration(schema, config); err != nil {
		return driverInstance, err
	}
	if err := driverInstance.TestConnection(ctx, config); err != nil {
		return driverInstance, fmt.Errorf("%w: %v", driver.ErrDriverConnectionFailed, err)
	}
	return driverInstance, nil
}

func (m *Driver[T]) ListDrivers(ctx context.Context) ([]driver.Info, error) {
	names := m.registry.Names()
	drivers := make([]driver.Info, 0, len(names))

	for _, name := range names {
		driver, err := m.GetDriver(name)
		if err != nil {
			return nil, err
		}
		info := driver.Info()
		info.Key = name
		drivers = append(drivers, info)
	}

	return drivers, nil
}
