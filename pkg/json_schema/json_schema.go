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

func MustLoad(schemaJSON []byte) util.JSONMap {
	schema, err := Load(schemaJSON)
	if err != nil {
		panic(err)
	}
	return schema
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

func ParseAndValidate[T any](schema util.JSONMap, data util.JSONMap) (*T, error) {
	if err := Validate(schema, data); err != nil {
		return nil, err
	}

	parsed, err := util.ParseJSONMap[T](data)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
