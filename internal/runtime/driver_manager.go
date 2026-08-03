package runtime

import (
	"context"
	"fmt"

	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/internal/util"
	"github.com/usesnipet/snipet/pkg/driver"
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
)

type DriverManager[T driver.IDriver] struct {
	registry *registry.R[T]
}

func NewDriverManager[T driver.IDriver](registry *registry.R[T]) *DriverManager[T] {
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

func (m *DriverManager[T]) ValidateConfiguration(schema util.JSONMap, config util.JSONMap) error {
	return jsonschema.Validate(schema, config)
}

func (m *DriverManager[T]) ValidateConfigurations(schema util.JSONMap, configs ...util.JSONMap) error {
	for _, config := range configs {
		if err := jsonschema.Validate(schema, config); err != nil {
			return err
		}
	}
	return nil
}

func (m *DriverManager[T]) ValidateConfigurationByKey(key string, config util.JSONMap) error {
	driver, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfiguration(driver.Info().ConfigurationSchema, config)
}

func (m *DriverManager[T]) ValidateConfigurationsByKey(key string, configs ...util.JSONMap) error {
	driver, err := m.GetDriver(key)
	if err != nil {
		return err
	}
	return m.ValidateConfigurations(driver.Info().ConfigurationSchema, configs...)
}

type Configuration struct {
	Key    string       `json:"key"`
	Config util.JSONMap `json:"config"`
}

func (m *DriverManager[T]) ValidateMultipleConfigurationsByKey(configs ...Configuration) error {
	for _, c := range configs {
		if err := m.ValidateConfigurationByKey(c.Key, c.Config); err != nil {
			return err
		}
	}
	return nil
}

func (m *DriverManager[T]) Prepare(ctx context.Context, driverKey string, config util.JSONMap) (T, error) {
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

func (m *DriverManager[T]) ListDrivers(ctx context.Context) ([]driver.Info, error) {
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
