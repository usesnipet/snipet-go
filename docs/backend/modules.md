# `internal/module/<name>`

A module is one business domain (`llm`, `agent`, `client`, `knowledge`, ...):
everything needed to expose CRUD (and any domain-specific operations) for
that entity over HTTP. Scaffolding one is covered step-by-step by the
[`create-backend-module` skill](../../.claude/skills/create-backend-module/SKILL.md);
this doc explains why each file exists and how they relate.

## Anatomy

```
internal/module/<name>/
  dto.go        # request/response shapes + validation tags
  service.go     # business logic, talks to repositories
  handler.go      # HTTP layer: routes, parsing, status codes
  service_test.go  # tests service.go against mocked repositories
```

Domain CRUD belongs here — never in `drivers/` or `pkg/driver/`, which are
for pluggable backend integrations (see [drivers.md](./drivers.md)), not
application domain logic.

## `dto.go`

Three DTO shapes per entity, each with a distinct job:

```go
type CreateLLMDTO struct {
    Name          string        `json:"name" validate:"required,max=255"`
    Provider      string        `json:"provider" validate:"required,max=255"`
    Configuration jsonx.JSONMap `json:"configuration" validate:"required"`
}

type UpdateLLMDTO struct {
    Name          *string       `json:"name" validate:"omitempty,max=255"`
    Provider      *string       `json:"provider" validate:"omitempty,max=255"`
    Configuration jsonx.JSONMap `json:"configuration" validate:"omitempty"`
}

type FindLLMsFilterDTO struct {
    Take *int `form:"take" validate:"omitempty,min=1"`
    Skip *int `form:"skip" validate:"omitempty,min=0"`
}

func (dto *FindLLMsFilterDTO) ToFilter() *filter.Options[model.LLM] {
    return filter.New[model.LLM](filter.PtrTake(dto.Take), filter.PtrSkip(dto.Skip))
}
```

- **Create DTO** — value fields, `validate:"required"` on what the entity
  can't exist without.
- **Update DTO** — every field is a **pointer**, `validate:"omitempty"`.
  `nil` means "leave unchanged" — this is what makes `PUT` a partial patch
  instead of a full replace (see how `service.go` uses this below).
- **List/filter DTO** — `form:"..."` tags (parsed from query string by
  `api.ParseQuery`, see [api.md](./api.md)) plus a `ToFilter()` method that
  turns it into `*filter.Options[T]` (see
  [filter-page.md](./filter-page.md)). This is the only place query-string
  filtering knowledge lives — `service.go` just takes `*filter.Options[T]`.

## `service.go`

Owns business logic; depends on repository **interfaces** (never a concrete
`*XRepository`), so it's mockable in tests and swappable in wiring.

```go
type Service struct {
    repo       repository.ILLMRepository
    llmManager *manager.Driver[llm.Driver]
}

func NewService(repo repository.ILLMRepository, llmManager *manager.Driver[llm.Driver]) *Service {
    return &Service{repo: repo, llmManager: llmManager}
}
```

Standard method set: `Filter`, `FindByID`, `Create`, `Update`, `DeleteByID`
— thin pass-throughs to the repository when there's no extra logic, doing
real work only where the domain needs it (e.g. `LLM.Create` validates
`Configuration` against the provider's JSON Schema via `llmManager` before
persisting, see [drivers.md](./drivers.md)).

For a module scoped under `/tenants/{tenant_id}/...`, every method also
takes `tenantID string` as its first parameter and starts with an
`authz.RequireMember`/`RequireAdmin` check (see
[auth-middleware.md](./auth-middleware.md)):

**The pointer-fields-mean-partial-update pattern**, from `LLM.Update`:

```go
func (s *Service) Update(ctx context.Context, tenantID, id string, dto UpdateLLMDTO) error {
    if _, err := authz.RequireMember(ctx, tenantID); err != nil {
        return err
    }
    if _, err := s.findInTenant(ctx, tenantID, id); err != nil { // fetch-then-compare, see below
        return err
    }

    updates := &model.LLM{}
    if dto.Name != nil {
        updates.Name = *dto.Name
    }
    if dto.Provider != nil {
        updates.Provider = *dto.Provider
    }
    return s.repo.UpdateByID(ctx, id, updates)
}
```

Only fields the caller actually set get copied onto `updates`; everything
else stays at its Go zero value, which GORM's `Updates` call treats as "not
included in the `SET` clause" (see [repository.md](./repository.md)) — so an
omitted field in the request body genuinely means "don't touch this column."

**Domain errors** are returned as `*apperr.Error` (`apperr.BadRequest(...)`,
etc. — see [errors.md](./errors.md)), never a bare `errors.New(...)`, so the
HTTP layer knows the right status code without the service knowing about
HTTP at all.

**Multi-step writes** that touch more than one table go through
`repository.ITxManager.WithTransaction` (see `agent.Service.Create`/
`Update`, which create the agent row and replace its LLM associations
atomically).

## `handler.go`

The only layer that knows about HTTP. Constructs an `api.Handler` (see
[api.md](./api.md)):

```go
func NewHandler(service *Service, userMiddleware api.MiddlewareFunc) api.Handler {
    return &Handler{service: service, userMiddleware: userMiddleware}
}

func (h *Handler) RegisterRoutes(r chi.Router, serve api.ServeFunc) {
    r.Route("/tenants/{tenant_id}/llm", func(r chi.Router) {
        r.Use(h.userMiddleware)
        r.Get("/", serve(h.filter))
        r.Post("/", serve(h.create))
        r.Get("/{id}", serve(h.findByID))
        r.Put("/{id}", serve(h.update))
        r.Delete("/{id}", serve(h.deleteByID))
    })
}

func (h *Handler) filter(w http.ResponseWriter, r *http.Request) error {
    data, err := h.service.Filter(r.Context(), chi.URLParam(r, "tenant_id"), query.ToFilter())
    ...
}
```

- Route group + `r.Use(...)` is where a module wires its auth requirement
  (see [auth-middleware.md](./auth-middleware.md)). Most modules' admin/CRUD
  surface lives under `/tenants/{tenant_id}/...`, gated by
  `guard.RequireUser` — every handler method reads `tenant_id` via
  `chi.URLParam` and passes it as the first argument into the matching
  service call, which checks it with `authz.RequireMember`/`RequireAdmin`
  (see [auth-middleware.md](./auth-middleware.md#internal-authz--tenant-membership-on-top-of-guardrequireuser)).
  A module can still mix gates across route groups when some routes need
  different auth than others — `agent`'s CRUD lives under
  `/tenants/{tenant_id}/agents` (`requireUser`) while `POST /agent/{id}/run`
  stays on `requireApiKey` at its own unprefixed path; `client`/`session`
  use `Or(...)` compositions for their client-widget-facing routes.
- Every handler method has the shape `func(w, r) error` — parse with
  `api.ParseBody`/`api.ParseQuery`, call the service, write with
  `api.WriteJSON`/`api.WriteNoContent`, and just `return err` on failure.
  `serve(...)` (passed in from `bootstrap`, see [bootstrap.md](./bootstrap.md))
  is what turns that returned error into the actual HTTP response.
- Status codes: list/get → `200`, create → `201`, update/delete → `204`
  (`api.WriteNoContent`). Path params via `chi.URLParam(r, "id")`.

## `service_test.go`

External test package (`package foomodule_test`), `t.Parallel()`, testify +
generated mocks:

```go
repo := mocks.NewMockIFooRepository(t)
repo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
svc := newTestService(repo)
```

A `newTestService(repo, ...)` helper wires real collaborators when they're
cheap (a hasher, a logger) and mocks when they're not (repositories, driver
managers). Assert domain failures by checking the returned `*apperr.Error`'s
status/message, not just that an error occurred.

## Wiring

A module's `Service`/`Handler` don't construct their own dependencies —
`internal/bootstrap` builds the repository, then the service, then the
handler, in that order, and calls `RegisterRoutes`. See
[bootstrap.md](./bootstrap.md).
