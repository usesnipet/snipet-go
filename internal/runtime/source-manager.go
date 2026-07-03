package runtime

import (
	"github.com/usesnipet/snipet/internal/util"
	"github.com/xeipuuv/gojsonschema"
)

type SourceManager struct {
	registry *Registry[SourceDriver]
}

func NewSourceManager(
	registry *Registry[SourceDriver],
) *SourceManager {
	return &SourceManager{
		registry: registry,
	}
}

func (m *SourceManager) GetSourceDriver(name string) (SourceDriver, error) {
	driver, ok := m.registry.Get(name)
	if !ok {
		return nil, ErrSourceDriverNotFound
	}
	return driver, nil
}

func (m *SourceManager) ValidateConfiguration(schema util.JSONMap, config util.JSONMap) error {
	referenceLoader := gojsonschema.NewGoLoader(schema)
	configLoader := gojsonschema.NewGoLoader(config)

	result, err := gojsonschema.Validate(referenceLoader, configLoader)

	if err != nil || !result.Valid() {
		return ErrInvalidConfiguration
	}

	return nil
}
