package runtime

import (
	"context"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/xeipuuv/gojsonschema"
)

type IDriver interface {
	GetConfigurationSchema(ctx context.Context) (util.JSONMap, error)
	TestConnection(ctx context.Context, config util.JSONMap) error
}

type DriverInfo struct {
	Name                string       `json:"name"`
	ConfigurationSchema util.JSONMap `json:"configuration_schema"`
}

type Manager[T IDriver] struct {
	registry *Registry[T]
}

func NewManager[T IDriver](registry *Registry[T]) *Manager[T] {
	return &Manager[T]{registry: registry}
}

func (m *Manager[T]) GetDriver(name string) (T, error) {
	driver, ok := m.registry.Get(name)
	if !ok {
		return driver, ErrDriverNotFound
	}
	return driver, nil
}

func (m *Manager[T]) ValidateConfiguration(schema util.JSONMap, config util.JSONMap) error {
	referenceLoader := gojsonschema.NewGoLoader(schema)
	configLoader := gojsonschema.NewGoLoader(config)

	result, err := gojsonschema.Validate(referenceLoader, configLoader)

	if err != nil || !result.Valid() {
		return ErrInvalidConfiguration
	}

	return nil
}

func (m *Manager[T]) Prepare(ctx context.Context, driver string, config util.JSONMap) (T, error) {
	driverInstance, err := m.GetDriver(driver)
	if err != nil {
		return driverInstance, err
	}
	configurationSchema, err := driverInstance.GetConfigurationSchema(ctx)
	if err != nil {
		return driverInstance, err
	}
	if err := m.ValidateConfiguration(configurationSchema, config); err != nil {
		return driverInstance, err
	}
	if err := driverInstance.TestConnection(ctx, config); err != nil {
		return driverInstance, ErrConnectionFailed
	}
	return driverInstance, nil
}

func (m *Manager[T]) ListDrivers(ctx context.Context) ([]DriverInfo, error) {
	names := m.registry.Names()
	drivers := make([]DriverInfo, 0, len(names))

	for _, name := range names {
		driver, err := m.GetDriver(name)
		if err != nil {
			return nil, err
		}

		schema, err := driver.GetConfigurationSchema(ctx)
		if err != nil {
			return nil, err
		}

		drivers = append(drivers, DriverInfo{
			Name:                name,
			ConfigurationSchema: schema,
		})
	}

	return drivers, nil
}
