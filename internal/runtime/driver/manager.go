package driver

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/registry"
	"github.com/usesnipet/snipet/internal/util"
	jsonschema "github.com/usesnipet/snipet/internal/util/json_schema"
)

type Manager[T IDriver] struct {
	registry *registry.R[T]
}

func NewManager[T IDriver](registry *registry.R[T]) *Manager[T] {
	return &Manager[T]{registry: registry}
}

func (m *Manager[T]) GetDriver(key string) (T, error) {
	driver, ok := m.registry.Get(key)
	if !ok {
		return driver, ErrDriverNotFound
	}
	return driver, nil
}

func (m *Manager[T]) ValidateConfiguration(schema util.JSONMap, config util.JSONMap) error {
	return jsonschema.Validate(schema, config)
}

func (m *Manager[T]) ValidateConfigurations(schema util.JSONMap, configs ...util.JSONMap) error {
	for _, config := range configs {
		if err := jsonschema.Validate(schema, config); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager[T]) ValidateConfigurationByKey(key string, config util.JSONMap) error {
	driver, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfiguration(driver.Info().ConfigurationSchema, config)
}

func (m *Manager[T]) ValidateConfigurationsByKey(key string, configs ...util.JSONMap) error {
	driver, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfigurations(driver.Info().ConfigurationSchema, configs...)
}

type Configuration struct {
	Key    string
	Config util.JSONMap
}

func (m *Manager[T]) ValidateMultipleConfigurationsByKey(configs ...Configuration) error {
	for _, c := range configs {
		if err := m.ValidateConfigurationByKey(c.Key, c.Config); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager[T]) Prepare(ctx context.Context, driver string, config util.JSONMap) (T, error) {
	driverInstance, err := m.GetDriver(driver)
	if err != nil {
		return driverInstance, err
	}
	schema := driverInstance.Info().ConfigurationSchema
	if err := m.ValidateConfiguration(schema, config); err != nil {
		return driverInstance, err
	}
	if err := driverInstance.TestConnection(ctx, config); err != nil {
		return driverInstance, fmt.Errorf("%w: %v", ErrDriverConnectionFailed, err)
	}
	return driverInstance, nil
}

func (m *Manager[T]) ListDrivers(ctx context.Context) ([]Info, error) {
	names := m.registry.Names()
	drivers := make([]Info, 0, len(names))

	for _, name := range names {
		driver, err := m.GetDriver(name)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driver.Info())
	}

	return drivers, nil
}
