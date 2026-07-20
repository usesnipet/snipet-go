package jsonschema

import (
	"encoding/json"
	"fmt"

	"github.com/usesnipet/snipet/internal/util"
	"github.com/xeipuuv/gojsonschema"
)

func Load(schemaJSON []byte) (util.JSONMap, error) {
	var schema util.JSONMap
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}
	return schema, nil
}

func Validate(schema util.JSONMap, json util.JSONMap) error {
	referenceLoader := gojsonschema.NewGoLoader(schema)
	jsonLoader := gojsonschema.NewGoLoader(json)

	result, err := gojsonschema.Validate(referenceLoader, jsonLoader)

	if err != nil {
		return err
	}
	if !result.Valid() {
		firstErr := result.Errors()[0]
		return fmt.Errorf("%s: %s", firstErr.Field(), firstErr.Description())
	}

	return nil
}

func ParseAndValidate[T any](schema util.JSONMap, json util.JSONMap) (*T, error) {
	parsed, err := util.ParseJSONMap[T](schema)
	if err != nil {
		return nil, err
	}

	if err := Validate(schema, json); err != nil {
		return nil, err
	}

	return &parsed, nil
}
