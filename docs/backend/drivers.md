# Drivers: `pkg/driver` + `drivers/` + registry/manager

The plugin system behind every pluggable backend integration: LLM
providers, tools, and knowledge sources/indexes. Three layers, each with a
distinct job:

```
pkg/driver/...              framework-agnostic CONTRACT (interfaces, shared types)
drivers/<kind>/<provider>/    CONCRETE IMPLEMENTATION of that contract
drivers/<kind>/driver.go        REGISTRY of every concrete implementation for that kind
internal/runtime/registry         generic name → instance store (used by the line above)
internal/runtime/manager           generic operations over a registry (validate, prepare, list, invoke)
```

## `pkg/driver` — the contracts

Framework-agnostic, no dependency on `internal/*` — see the package's own
Go doc comments (`go doc ./pkg/driver/...`) for the full contract. The
shape every driver kind shares:

```go
type Info struct {
    Key, Name, Description, Icon string
    Tags                []string
    ConfigurationSchema jsonx.JSONMap  // JSON Schema validating this driver's config
}

type IDriver interface {
    Info() Info
    TestConnection(ctx context.Context, config jsonx.JSONMap) error
}
```

Each kind extends `IDriver` with its own operations:
`pkg/driver/llm.Driver` adds `Stream`/`Generate`/`Models`/`Model`;
`pkg/driver/tool.Driver` adds `ToolSet`/`Call`;
`pkg/driver/knowledge.ISourceDriver`/`IIndexDriver` add source
iteration/reading and index read/write respectively.

## `drivers/<kind>/<provider>/` — one concrete implementation

A provider is built with the kind's `CreateDriver(opts ...Option)` functional
builder, composing `Info` fields with an `API` (the actual
`TestConnection`/`Generate`/`Stream`/`Call` functions):

```go
// drivers/llm/openai/main.go
func New() llm.Driver {
    return llm.CreateDriver(
        llm.WithName("OpenAI"),
        llm.WithDescription("OpenAI language models."),
        llm.WithConfigurationSchema(openaicompatible.DefaultConfigSchema),
        llm.WithAPI(openaicompatible.New(baseURL)),
    )
}

// drivers/tool/swapi/main.go
func New() tool.Driver {
    return tool.CreateDriver(
        tool.WithName("SWAPI"),
        tool.WithToolSetSchema(toolsJSON),   // //go:embed'd JSON Schema
        tool.WithAPI(NewAPI()),
    )
}
```

Several LLM providers (OpenAI, Groq, Mistral, Ollama, OpenRouter) share one
`API` implementation, `pkg/driver/llm/api/openai_compatible`, since they all
speak the OpenAI Chat Completions wire format — only `baseURL` and
`ConfigurationSchema` differ per provider. A genuinely different wire
protocol gets its own `API` implementation instead of being forced into
that shape.

## `drivers/<kind>/driver.go` — the registry

```go
func Registry() *registry.R[llmDriver.Driver] {
    registry := registry.New[llmDriver.Driver]()
    registry.MustRegister("openai", openai.New())
    registry.MustRegister("groq", groq.New())
    // ...
    return registry
}
```

One `Registry()` per kind (`drivers/llm`, `drivers/tool`, `drivers/source`,
`drivers/index`), each a flat list of every concrete provider available for
that kind, keyed by the string a `model.LLM.Provider` (or equivalent)
column stores. Adding a new provider means writing its `drivers/<kind>/<provider>/`
package and adding one `MustRegister` line here — nothing else in the app
needs to change. `drivers/index` and `drivers/source` are currently empty
registries (commented-out example) — knowledge indexing has no built-in
provider yet.

## `internal/runtime/registry.R[T]` — generic store

```go
type R[T any] struct { /* thread-safe map[string]T */ }
func (r *R[T]) Register(name string, value T) error
func (r *R[T]) Get(name string) (T, bool)
func (r *R[T]) Names() []string
```

Just a concurrency-safe named lookup, generic over any driver interface —
`drivers/<kind>/driver.go`'s `Registry()` return type is always
`*registry.R[SomeDriverInterface]`.

## `internal/runtime/manager` — generic operations

```go
type DriverManager[T driver.IDriver] struct { registry *registry.R[T] }

func (m *DriverManager[T]) GetDriver(key string) (T, error)
func (m *DriverManager[T]) ValidateConfigurationByKey(key string, config jsonx.JSONMap) error
func (m *DriverManager[T]) Connect(ctx context.Context, driverKey string, config jsonx.JSONMap) (T, error)  // validate + TestConnection
func (m *DriverManager[T]) ListDrivers(ctx context.Context) ([]driver.Info, error)
```

`manager.DriverManager[T]` wraps a `registry.R[T]` with the operations every
consumer actually needs — validating a config against the driver's JSON
Schema (used by `LLM.Create`/`Update`, see [modules.md](./modules.md)),
listing available drivers for a "choose a provider" UI
(`LLM.ListDrivers` → `GET /llm/drivers`), or resolving + connection-testing
one before use (`Connect`).

**`manager.Toolbox`** is a special case built *on top of* `manager.DriverManager[tool.Driver]`,
because unlike LLM/source/index, every registered tool driver's tools are
meant to be available to the agent simultaneously, not selected one at a
time:

```go
func (m *Toolbox) Toolset() (tool.Toolset, error)  // every tool from every registered driver, namespaced "<driverKey>__<toolName>"
func (m *Toolbox) Call(ctx context.Context, call tool.Call) (tool.Result, error)  // routes by the namespace prefix back to the owning driver
```

This is what `runtime.Engine` asks for the current toolset and dispatches
calls through (see [runtime.md](./runtime.md)).

## Where managers get built

Every `manager.DriverManager[T]`/`manager.Toolbox` is built once in
`bootstrap.Bootstrap` from that kind's `Registry()` and handed to whichever
service/engine needs it (`llmModule.NewService(llmRepo, llmManager)`,
`runtime.NewEngine(llmManager, toolManager, logger)`) — see
[bootstrap.md](./bootstrap.md).
