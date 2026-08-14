---
name: create-web-feature
description: Scaffold a new frontend feature folder (schemas, service, hooks, store, components) under web/src/features following this project's conventions. Use when creating or adding a new feature, page, or API-backed module in the web app.
---

# Web Features

Create domain features under `web/src/features/<feature-name>/` (kebab-case). Include only the layers the feature needs.

## Structure

```
web/src/features/<feature>/
  schemas.ts              # required when there are types/API payloads
  service.ts              # HTTP calls
  hooks.ts                # React Query wrappers
  store.ts                # optional client/persisted state (Zustand)
  components/             # optional feature-owned UI
```

## schemas.ts

- Define Zod schemas; export types with `z.infer<typeof schema>`.
- Prefer `.strict()` on objects; use `.pick()` / `.extend()` for create/update DTOs.
- Coerce dates with `z.coerce.date()` when the API returns ISO strings.
- Reuse shared schemas (e.g. `paginatedSchema` from `@/schemas/paginated`).

```typescript
export const fooSchema = z.object({ id: z.string(), name: z.string() }).strict();
export type Foo = z.infer<typeof fooSchema>;

export const createFooSchema = fooSchema.pick({ name: true }).strict();
export type CreateFoo = z.infer<typeof createFooSchema>;
```

## service.ts

- Call `@/lib/http` (`http.get` / `post` / `put` / `delete`).
- Accept `ServiceGetOptions` / `ServicePostOptions` / `ServicePutOptions` from `@/lib/services`.
- Pass Zod `schemas` for `body` and/or `response`; use `params` for path placeholders (`{id}`).
- Prefer a named object export (`fooService`) when there are multiple methods; named function exports are fine for small features.

```typescript
const FOO_URL = "/api/foo";

const list = async (opts: ServiceGetOptions<Paginated<Foo>> = {}) =>
  http.get({ url: FOO_URL, schemas: { response: paginatedFooSchema }, ...opts });

export const fooService = { list };
```

## hooks.ts

- Wrap services with `useQuery` / `useMutation` from `@tanstack/react-query`.
- Export query/mutation key factories: `listFooQueryKey`, `useListFoo`, etc.
- On mutations: toast success/error via `@/hooks/use-toast`, then `queryClient.invalidateQueries`.
- Forward service opts; set auth in the hook when required (`auth: "api-key"`).

## store.ts (optional)

- Use Zustand `create`.

## components/ (optional)

- Own only UI tightly coupled to the feature.
- Import hooks/schemas/store from the parent feature folder (`../hooks`, etc.).
- Shared/generic UI stays in `web/src/components`.

## New routes

If this feature adds a routed page, register it in `web/src/router.tsx` — routes there must be lazy-loaded (see project convention in `web/CLAUDE.md`).
