package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/usesnipet/snipet/pkg/driver/tool"
)

// toolNameSeparator namespaces tool names by the driver that owns them
// (e.g. "swapi__list_people") to avoid collisions between drivers. It must
// only contain characters valid in OpenAI-compatible function names
// ([a-zA-Z0-9_-]), so "." is not an option.
const toolNameSeparator = "__"

// ToolManager aggregates every registered tool driver into a single toolset
// the LLM can be given, and dispatches calls back to the owning driver.
type ToolManager struct {
	drivers *DriverManager[tool.Driver]
}

func NewToolManager(drivers *DriverManager[tool.Driver]) *ToolManager {
	return &ToolManager{drivers: drivers}
}

// Toolset returns every tool from every registered driver, by default all
// tools installed are made available to the LLM.
func (m *ToolManager) Toolset() (tool.Toolset, error) {
	var tools []tool.Tool
	for _, key := range m.drivers.Names() {
		driverInstance, err := m.drivers.GetDriver(key)
		if err != nil {
			return tool.Toolset{}, err
		}
		for _, t := range driverInstance.ToolSet().Tools {
			tools = append(tools, tool.NewTool(namespacedToolName(key, t.Name), t.Description, t.Parameters))
		}
	}
	return tool.NewToolset(tools...), nil
}

// Call dispatches a namespaced tool call to the driver that owns it.
func (m *ToolManager) Call(ctx context.Context, call tool.Call) (tool.Result, error) {
	driverKey, toolName, ok := strings.Cut(call.Tool, toolNameSeparator)
	if !ok {
		return tool.Result{}, fmt.Errorf("%w: %q", ErrToolNotFound, call.Tool)
	}

	driverInstance, err := m.drivers.GetDriver(driverKey)
	if err != nil {
		return tool.Result{}, fmt.Errorf("%w: %q", ErrToolNotFound, call.Tool)
	}

	result, err := driverInstance.Call(ctx, tool.Call{Tool: toolName, Arguments: call.Arguments})
	if err != nil {
		return tool.Result{}, err
	}
	result.Tool = call.Tool
	return result, nil
}

func namespacedToolName(driverKey, toolName string) string {
	return driverKey + toolNameSeparator + toolName
}
