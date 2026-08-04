package driver

import (
	"github.com/usesnipet/snipet/internal/util"
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
)

// ConfigurationSchema converts a JSON schema document into a util.JSONMap.
func ConfigurationSchema(schemaJSON []byte) (util.JSONMap, error) {
	return jsonschema.Load(schemaJSON)
}

// MustConfigurationSchema is like ConfigurationSchema but panics on error.
func MustConfigurationSchema(schemaJSON []byte) util.JSONMap {
	schema, err := ConfigurationSchema(schemaJSON)
	if err != nil {
		panic(err)
	}
	return schema
}
