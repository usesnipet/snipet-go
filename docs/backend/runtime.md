# `internal/runtime` — the agent execution engine

Runs an agent's LLM ↔ tool conversation loop and publishes what happens as
a stream of typed events, so callers (HTTP streaming, persistence, future
subscribers) can observe an execution without the engine knowing anything
about HTTP or a database.

```
internal/runtime/
  engine.go             Engine — the turn loop
  execution/              Execution — state + event bus for one run
  generator/                Generate — one LLM call → one assistant message
  tool_executor/               ToolExecutor — runs the tool calls a message requested
  manager/                       Driver[T]/Tool — see drivers.md
  registry/                        R[T] — see drivers.md
```

## `execution.Execution` — state + event bus for one run

```go
type Execution struct {
    Agent     *Agent
    publisher IPublisher
    Config    Config          // MaxTurns, Metadata
    Status       Status       // pending → running → completed | failed | max_turns | cancelled
    Messages     []msg.Message
    Turns        int
}

exe, err := execution.NewExecution(
    execution.WithAgent(agent),
    execution.WithInitialMessages(msg.NewMessage(msg.RoleUser, "...")),
    execution.WithMaxTurns(10),
)
```

`Execution` owns the conversation (`Messages`) and turn count for one run,
and is the only thing both `Engine` and any caller-supplied subscriber talk
to. Its mutating methods (`AddMessage`, `CompleteTurn`, `Finish`,
`SetStatus`, `SetError`, `SetMaxTurnsReachedError`, `Cancel`) all end by
publishing a matching event — state changes and notifications happen
atomically from the caller's point of view, so a subscriber can never
observe a state change without the event that describes it (or vice versa).

Every status change goes through one private `transition` helper that sets
`Status` + `ErrorMessage` and publishes `StatusChangedEvent`. The terminal
outcomes then publish a **second, specific** event on top of it:
`Finish` → `FinishedEvent`, `SetError` → `FailedEvent{Error}`,
`SetMaxTurnsReachedError` → `MaxTurnsReachedEvent{MaxTurns}`,
`Cancel` → `CancelledEvent`. So a subscriber that only tracks status
switches on `StatusChangedEvent` alone; one that wants to react to *how* a
run ended switches on the specific events.

## Events and subscribers

```go
type Subscriber interface {
    Handle(context.Context, IEvent) error
}
```

Events are a closed set (sealed via an embedded `event{}` marker, same
pattern as `pkg/driver/llm.StreamEvent` — a subscriber type-switches on the
concrete type): `StartedEvent`, `StatusChangedEvent`, `FinishedEvent`,
`FailedEvent`, `MaxTurnsReachedEvent`, `CancelledEvent`, `MessageAddedEvent`,
`TurnStartedEvent`, `TurnCompletedEvent`, `MessageDeltaEvent` (incremental
text while the LLM streams), `MessageAttemptFailedEvent`,
`ToolCallStartedEvent`, `ToolCallResultEvent`.

Two subscribers ship today, both in `internal/module/agent/subscriber/`,
both attached by `agent.Service.Run`:

- **`Persistence`** — writes every event to the database (`ExecutionMessage`
  rows, `Execution.Status`/`Turns` updates) via the repositories, so an
  execution survives past the request that started it.
- **`SSE`** — translates each event into a named
  Server-Sent-Event frame (`api.SSEWriter`, see [api.md](./api.md)) so a
  client watching `POST /agent/{id}/run` sees the conversation as it
  happens.

A new cross-cutting concern (e.g. metrics, audit logging) is a new
`Subscriber` implementation attached the same way — it doesn't require
touching `Engine` or `Execution`.

**Current behavior to know before touching this:** `LocalPublisher.Publish`
stops at the first subscriber whose `Handle` returns an error, and doesn't
call the rest — this is intentional today because every attached subscriber
is considered required (if `Persistence` fails, the run is supposed to
stop). It's worth revisiting once an *optional* subscriber exists (e.g. a
best-effort logger) whose failure shouldn't abort the run — at that point
isolate per-subscriber failures instead of aborting on the first one, and
add locking (`Subscribe`/`Publish` aren't currently safe to call
concurrently on the same `Execution`).

## `Engine` — the turn loop

```go
engine := runtime.NewEngine(llmManager, toolManager, logger)
err := engine.Start(ctx, exe)
```

`NewEngine` builds two collaborators once and holds them for the engine's
lifetime: a `generator.Generator` (LLM driver manager + a prefixed logger)
and a `tool_executor.ToolExecutor`.

`Start` validates the agent (an LLM key + config it can resolve and
validate via `manager.DriverManager[llm.Driver]`, see [drivers.md](./drivers.md)),
sets `StatusRunning`, then loops calling `step` until it returns a terminal
`StepResult` (a string enum; the zero value `StepResultInvalid` is only ever
paired with a non-nil error, meaning "read the error, ignore the result"):

```go
StepContinue        // assistant responded with tool calls, or a non-final message — loop again
StepFinish           // assistant sent a final message with no tool calls — exe.Finish()
StepCancel            // ctx was cancelled — exe.Cancel()
StepMaxTurnsReached    // exe.Turns >= exe.Config.MaxTurns — exe.SetMaxTurnsReachedError()
```

One `step`:

1. `exe.StartTurn()` — increments `Turns` and publishes `TurnStartedEvent`
   (only after the context and max-turns checks pass, so the max-turns turn
   never emits a spurious start).
2. Resolve the current `Toolset` from `manager.Toolbox` (empty toolset if it
   can't resolve, rather than aborting the whole execution).
3. `generator.Generate` — one streamed LLM call producing one assistant
   `msg.Message` (see below); its `MessageDeltaEvent`/`MessageAttemptFailedEvent`
   are published *during* generation, before the message is even complete.
4. `exe.AddMessage` — appends it and publishes `MessageAddedEvent`.
5. If the message has tool calls, `toolExecutor.Run` executes them and
   appends a `RoleTool` message per result (see below) — the loop then
   continues, feeding the tool results back to the LLM on the next turn.
6. Otherwise, `message.IsFinal()` decides `StepFinish` vs `StepContinue`.

## `generator.Generator` — one LLM call

`NewGenerator(llms, logger)` is built once by `NewEngine`; `Generate(ctx,
exe, toolset)` is called per turn. It resolves the agent's configured
`llm.Driver`, checks the resolved `Model` supports
`ModelCapabilitiesToolCall` (returns `ErrModelNotSupportToolCall` otherwise
— the engine currently requires tool-call-capable models, it doesn't
degrade to a non-tool-call mode), then calls `Driver.Stream` and drains it:
`TextDeltaEvent`s become `MessageDeltaEvent` publishes and accumulate into
the final message's `Content`; `ToolCallEvent`s accumulate into `ToolCalls`.
A stream error after partial content was already published emits
`MessageAttemptFailedEvent` — the caller (the engine, via `step`'s returned
error) is expected to treat that attempt's content as discarded, never as
part of the conversation history.

## `tool_executor.ToolExecutor` — running tool calls

For each `tool.Call` on the assistant's message: publishes
`ToolCallStartedEvent`, calls `manager.Toolbox.Call` (routes to the owning
driver by namespace prefix, see [drivers.md](./drivers.md)), publishes
`ToolCallResultEvent` (with `Error` set instead of `Result` on failure — a tool
failure does **not** abort the execution, it's surfaced to the LLM as an
error string so it can react), and appends a `RoleTool` message so the next
turn's prompt includes the result.

## How a module uses this

`internal/module/agent`'s `Run` (see [modules.md](./modules.md)) is the
glue: it loads the `Agent`/creates the `model.Execution` row, converts it to
a runtime `Execution` (`model.Execution.ToRuntimeExecution`, see
[model.md](./model.md)), subscribes `Persistence` + whatever the caller
passed (the HTTP handler passes an `SSE` subscriber), and calls
`engine.Start`. A different caller (a scheduled job, a test) can drive the
same engine with different subscribers — nothing here is HTTP-specific.
