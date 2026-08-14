# `src/models/`

The Zod schema for a feature's core entity — the shape with an `id` and
relations to other entities — lives in `src/models/<feature>.ts`, one file
per feature that has such an entity, not inside that feature's own
`schemas.ts`.

```typescript
// src/models/agent.ts
import { z } from "zod";

import { llmSchema } from "@/models/llm";

export const agentToLlmSchema = z
  .object({ llm_id: z.uuid(), priority: z.number().int(), llm: llmSchema })
  .strict();
export type AgentToLlm = z.infer<typeof agentToLlmSchema>;

export const agentSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    instructions: z.string().max(1000),
    llms: z.array(agentToLlmSchema).nullable(),
  })
  .strict();
export type Agent = z.infer<typeof agentSchema>;
```

## Why this exists

Entities reference each other — an `Agent` embeds `Llm`s, a `Session`
embeds an `Agent`. If those schemas lived in each feature's `schemas.ts`,
one feature would have to import another feature's `schemas.ts` to build
the relation, and the reverse relation (when it's added later) would
create an import cycle between the two features. Pulling every entity into
its own file under `src/models/`, with no feature ever importing another
feature's `schemas.ts` for this purpose, keeps the dependency direction
one-way: `models/*` never imports from `features/*`, and any file may
import from `models/*`.

## What goes here vs. `features/<feature>/schemas.ts`

- **`models/<feature>.ts`**: the entity schema and any value objects
  nested inside it (e.g. `clientConfigSchema` inside `Client`). Nothing
  else — no DTOs, no pagination wrappers, no search params.
- **`features/<feature>/schemas.ts`**: re-exports the entity from
  `@/models/<feature>` (so existing imports of `@/features/<feature>/schemas`
  keep working), then defines everything derived from it — create/update
  DTOs, `paginatedSchema(...)` wrappers, search params, response envelopes.
  See [features/schemas.md](./features/schemas.md).

```typescript
// features/llm/schemas.ts
import { llmSchema } from "@/models/llm";

export { llmSchema } from "@/models/llm";
export type { Llm } from "@/models/llm";

export const createLlmSchema = llmSchema
  .pick({ name: true, provider: true, configuration: true })
  .strict();
export type CreateLlm = z.infer<typeof createLlmSchema>;
```

A feature that needs another feature's entity (e.g. `session` embedding
`agent`, `agent` embedding `llm`) imports it from `@/models/<other>`
directly — never from `@/features/<other>/schemas`.

## When to skip this

Not every feature has an entity. A feature whose schemas are pure
config/singleton shapes with no `id` and no relations to another entity
(e.g. `app`'s `systemInfoSchema`/`appConfigSchema`) has no reason to split
— its schema stays entirely in that feature's own `schemas.ts`.

## Relation to `src/schemas/`

`src/models/` and `src/schemas/` (see [schemas.md](./schemas.md)) are both
"lives outside one feature folder" but for different reasons: `src/schemas/`
holds generic, domain-agnostic shapes reused by shape (`paginatedSchema`,
`driverInfoSchema`) — they don't represent an entity and don't participate
in entity relations. `src/models/` holds one specific feature's entity,
pulled out only because other entities reference it.
