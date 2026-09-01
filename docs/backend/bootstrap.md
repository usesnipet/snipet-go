# `cmd/api` + `internal/bootstrap` + `config/`

Where every other layer gets constructed and wired together. Nothing here
holds business logic — it's strictly composition.

## `config/`

One struct per concern (`ServerConfig`, `DatabaseConfig`, `LogConfig`,
`AuthConfig`, `SyncConfig`), composed into one `Config`, loaded from
environment variables (via `go-envconfig`, with `.env` support for local
dev — `config.Load()` walks up from the working directory looking for a
`.env` file):

```go
type Config struct {
    Server   ServerConfig   `env:", prefix=SERVER_"`
    Database DatabaseConfig `env:", prefix=DB_"`
    Log      LogConfig      `env:", prefix=LOG_"`
    Auth     AuthConfig     `env:", prefix=AUTH_"`
    Sync     SyncConfig     `env:", prefix=SYNC_"`
    Env      string         `env:"ENV, default=development"`
    DevProxy string         `env:"DEV_PROXY, default=http://localhost:5173"`
}
```

`AuthConfig` is where admin basic auth lives (`AUTH_BASIC_AUTH_USERNAME` /
`AUTH_BASIC_AUTH_PASSWORD`), alongside the JWT/refresh-token settings for
App end-user tokens.

A new cross-cutting setting gets a field (with `env:"..., default=..."`)
on the relevant sub-struct, not a bespoke `os.Getenv` call somewhere deep
in the code. `config.APIPrefix` (`"/api"`) lives here too — the one
constant every module's routes are mounted under.

## `cmd/api/main.go`

The actual binary entrypoint — deliberately tiny:

```go
func main() {
    cfg, err := config.Load()
    level, _ := logger.ParseLevel(cfg.Log.Level)
    appLogger := logger.NewLogger(level)
    bootstrap.Bootstrap(cfg, appLogger)
}
```

Its only job is to load config, build the root logger, and hand off to
`bootstrap.Bootstrap`. Anything more than that belongs in `bootstrap` or a
config struct instead.

## `internal/bootstrap.Bootstrap` — the wiring order

Everything is built by hand (no DI framework), in dependency order, inside
one function:

1. **Database** — `database.NewDatabase(cfg, logger)` (ensures the DB
   exists, opens the connection, runs migrations — see
   [migrations.md](./migrations.md) and [infra.md](./infra.md)).
2. **Repositories** — one `repository.NewXRepository(db, ...)` call per
   entity, plus `repository.NewTxManager(db)` (see
   [repository.md](./repository.md)). Some repositories depend on another
   (`sessionRepo` takes `clientRepo`) — construct in the order that
   satisfies those dependencies.
3. **Driver registries + managers** — `<kind>.Registry()` +
   `manager.NewDriver(registry)` for llm/tool/source/index, plus
   `manager.NewTool(...)` for the tool aggregator (see
   [drivers.md](./drivers.md)).
4. **Runtime engine** — `runtime.NewEngine(llmManager, toolManager, logger)`
   (see [runtime.md](./runtime.md)).
5. **Background workers** — the sync worker pool (`queue.NewPool(...).Start(...)`)
   and the knowledge sync workers that use it (see
   [infra.md](./infra.md)).
6. **Auth primitives** — `auth.NewAPIKeyGenerator`, `auth.NewKeyHasher`,
   `auth.NewJWTService(cfg.Auth, ...)` for App end-user claims
   (see [auth-middleware.md](./auth-middleware.md)).
7. **Services** — one `<module>.NewService(repo, ..., logger)` per module,
   each depending only on repository interfaces, managers, and other
   services it genuinely needs (see [modules.md](./modules.md)).
8. **Cache** — `cache.NewMemoryCache(...)` for anything a guard needs
   (see [infra.md](./infra.md)).
9. **Guards** — `guard.RequireBasicAuth` / `guard.RequireApiKey` /
   `guard.RequireAppKey` / `guard.RequireAppUser` (composed at the route
   level with `api.Or(...)`), built from the config/services/cache above
   and passed as raw `api.Gate`s into modules (see
   [auth-middleware.md](./auth-middleware.md)).
10. **Handlers** — one `<module>.NewHandler(service, middleware...)` per
    module.
11. **Routes** — `api.New()`, mount the built SPA (`web.Handler()`) as the
    catch-all, then every `handler.RegisterRoutes(r, api.Serve)` under
    `config.APIPrefix` (see [api.md](./api.md)).
12. **Serve** — `http.ListenAndServe`.

Adding a new module means adding one line to *most* of these steps (repo →
manager if it needs one → service → handler → register routes) — see the
`create-backend-module` skill's bootstrap step for the exact checklist.

## Import aliasing

Hyphenated module directories (`api-key`, `knowledge-index`) can't be
imported under their literal package name, so `bootstrap` (and anything
else importing them) aliases: `apikey "github.com/.../module/api-key"`,
`knowledgeindex "github.com/.../module/knowledge-index"`.
