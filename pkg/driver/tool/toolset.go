package tool

import (
	_ "embed"

	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
)

//go:embed tool-set-schema.json
var toolSetSchemaJSON []byte

var toolSetSchema = jsonschema.MustLoad(toolSetSchemaJSON)

type Toolset struct {
	Tools []Tool `json:"tools"`
}

func ToolSetFromSchema(schemaJSON []byte) (Toolset, error) {
	set, err := jsonschema.ParseAndValidate[Toolset](
		toolSetSchema,
		jsonschema.MustLoad(schemaJSON),
	)
	if err != nil {
		return Toolset{}, err
	}
	return *set, nil
}

func NewToolset(tools ...Tool) Toolset {
	return Toolset{
		Tools: tools,
	}
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

func NewTool(name string, description string, parameters map[string]any) Tool {
	return Tool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}
}
