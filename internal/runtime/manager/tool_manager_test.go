package manager_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/internal/runtime/manager"
	"github.com/usesnipet/snipet/internal/runtime/registry"
	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

func newFakeToolDriver(name string, toolName string, call func(context.Context, tool.Call) (tool.Result, error)) tool.Driver {
	return tool.CreateDriver(
		tool.WithName(name),
		tool.WithDescription(name),
		tool.WithToolSet(tool.NewToolset(tool.NewTool(toolName, toolName, nil))),
		tool.WithAPI(tool.API{
			TestConnection: func(context.Context, jsonx.JSONMap) error { return nil },
			Call:           call,
		}),
	)
}

func TestToolManagerToolsetNamespacesByDriver(t *testing.T) {
	t.Parallel()

	reg := registry.New[tool.Driver]()
	reg.MustRegister("alpha", newFakeToolDriver("Alpha", "search", nil))
	reg.MustRegister("beta", newFakeToolDriver("Beta", "search", nil))

	manager := manager.NewTool(manager.NewDriver(reg))

	toolset, err := manager.Toolset()
	require.NoError(t, err)
	require.Len(t, toolset.Tools, 2)

	names := []string{toolset.Tools[0].Name, toolset.Tools[1].Name}
	require.Contains(t, names, "alpha__search")
	require.Contains(t, names, "beta__search")
}

func TestToolManagerCallDispatchesToOwningDriver(t *testing.T) {
	t.Parallel()

	var gotCall tool.Call
	reg := registry.New[tool.Driver]()
	reg.MustRegister("alpha", newFakeToolDriver("Alpha", "search", func(_ context.Context, call tool.Call) (tool.Result, error) {
		gotCall = call
		return tool.Result{Tool: call.Tool, Arguments: call.Arguments, Result: "ok"}, nil
	}))

	manager := manager.NewTool(manager.NewDriver(reg))

	result, err := manager.Call(context.Background(), tool.Call{
		Tool:      "alpha__search",
		Arguments: map[string]any{"q": "term"},
	})
	require.NoError(t, err)
	require.Equal(t, "search", gotCall.Tool, "driver should see the un-namespaced tool name")
	require.Equal(t, "alpha__search", result.Tool, "result should keep the namespaced name")
	require.Equal(t, "ok", result.Result)
}

func TestToolManagerCallUnknownToolReturnsErrToolNotFound(t *testing.T) {
	t.Parallel()

	toolManager := manager.NewTool(manager.NewDriver(registry.New[tool.Driver]()))

	_, err := toolManager.Call(context.Background(), tool.Call{Tool: "missing__tool"})
	require.ErrorIs(t, err, manager.ErrToolNotFound)

	_, err = toolManager.Call(context.Background(), tool.Call{Tool: "no-separator"})
	require.ErrorIs(t, err, manager.ErrToolNotFound)
}
