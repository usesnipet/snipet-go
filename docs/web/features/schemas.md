# `schemas.ts`

Every type the feature works with — API request/response shapes, DTOs
derived from its entity — is defined as a Zod schema, with the TypeScript
type derived from it via `z.infer`. The schema is the source of truth; the
type is a byproduct.

The entity itself (if the feature has one) isn't defined here — it lives in
[`@/models/<feature>`](../models.md) and is re-exported:

```typescript
import { llmSchema } from "@/models/llm";

export { llmSchema } from "@/models/llm";
export type { Llm } from "@/models/llm";
```

## Conventions

- **`.strict()`** on object schemas, so an unexpected field in an API
  response fails loudly instead of passing through silently — it means the
  backend contract drifted and the schema needs updating.
- **Derive create/update DTOs from the entity schema** (imported from
  `@/models/<feature>`) rather than redefining them: `.pick()` for a create
  payload, `.partial()` for an update payload, `.extend()` when a DTO needs
  a field the read model doesn't.
  ```typescript
  export const createLlmSchema = llmSchema
    .pick({ name: true, provider: true, configuration: true })
    .strict();
  export type CreateLlm = z.infer<typeof createLlmSchema>;

  export const updateLlmSchema = createLlmSchema.partial().strict();
  export type UpdateLlm = z.infer<typeof updateLlmSchema>;
  ```
- **List/pagination responses** reuse the shared factory from
  `@/schemas/paginated` instead of hand-rolling an envelope:
  ```typescript
  export const paginatedLlmSchema = paginatedSchema(llmSchema);
  export type PaginatedLlm = z.infer<typeof paginatedLlmSchema>;

  export const listLlmSearchParamsSchema = paginationParamsSchema;
  export type ListLlmSearchParams = z.infer<typeof listLlmSearchParamsSchema>;
  ```
- **Dates from the API** are ISO strings — coerce them with
  `z.coerce.date()` in the schema rather than parsing manually downstream.
- **Reuse cross-feature schemas** from `@/schemas` (see
  [../schemas.md](../schemas.md)) instead of duplicating them — e.g.
  `driverInfoSchema` is shared by every feature that lists driver-backed
  resources (`llm`, `knowledge`).

## Why Zod here, and not just TypeScript types

Schemas aren't just types — `service.ts` runs them at runtime
(`schemas.body` / `schemas.response` passed to `lib/http`, see
[service.md](./service.md)) to validate what actually goes over the wire.
A malformed response is caught and surfaced as a parse error instead of
silently propagating `undefined`s through the UI.
