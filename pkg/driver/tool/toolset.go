package tool

import (
	_ "embed"

	jsonschema "github.com/usesnipet/snipet/pkg/json_schema"
)

//go:embed tool-set-schema.json
var toolSetSchemaJSON []byte

// toolSetSchema is the JSON Schema a document must satisfy to be parsed by
// ToolSetFromSchema.
var toolSetSchema = jsonschema.MustLoad(toolSetSchemaJSON)

// Toolset is the set of tools a Driver exposes to the LLM.
type Toolset struct {
	Tools []Tool `json:"tools"`
}

// ToolSetFromSchema parses and validates schemaJSON against the embedded
// tool-set schema (tool-set-schema.json), returning the resulting Toolset.
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

// NewToolset builds a Toolset from the given Tools.
func NewToolset(tools ...Tool) Toolset {
	return Toolset{
		Tools: tools,
	}
}

// Tool describes a single callable tool: its Name, a Description the model
// uses to decide when to call it, and a JSON Schema (Parameters) describing
// its arguments.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// NewTool builds a Tool from its name, description, and parameter schema.
func NewTool(name string, description string, parameters map[string]any) Tool {
	return Tool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}
}
