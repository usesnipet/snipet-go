package driver

import (
	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// ConfigurationSchema converts a JSON schema document into a jsonx.JSONMap.
func ConfigurationSchema(schemaJSON []byte) (jsonx.JSONMap, error) {
	return jsonschema.Load(schemaJSON)
}

// MustConfigurationSchema is like ConfigurationSchema but panics on error.
func MustConfigurationSchema(schemaJSON []byte) jsonx.JSONMap {
	schema, err := ConfigurationSchema(schemaJSON)
	if err != nil {
		panic(err)
	}
	return schema
}
