import { z } from "zod";

export interface Llm {
  id: string;
  name: string;
  provider: string;
  configuration: Record<string, unknown>;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const llmBaseSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    provider: z.string().min(1).max(255),
    configuration: z.record(z.string(), z.unknown()),
  })
  .strict();

export const llmSchema: z.ZodType<Llm> = llmBaseSchema;
