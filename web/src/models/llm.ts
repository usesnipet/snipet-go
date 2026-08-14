import { z } from "zod";

export const llmSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    provider: z.string().min(1).max(255),
    configuration: z.record(z.string(), z.unknown()),
  })
  .strict();

export type Llm = z.infer<typeof llmSchema>;
