import { llmSchema } from "@/features/llm/schemas";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const agentToLlmSchema = z
  .object({
    llm_id: z.uuid(),
    priority: z.number().int(),
    llm: llmSchema,
  })
  .strict();

export type AgentToLlm = z.infer<typeof agentToLlmSchema>;

export const agentSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    instructions: z.string().max(1000),
    llms: z.array(agentToLlmSchema),
  })
  .strict();

export type Agent = z.infer<typeof agentSchema>;

export const paginatedAgentSchema = paginatedSchema(agentSchema);
export type PaginatedAgent = z.infer<typeof paginatedAgentSchema>;

/** Form / create payload shape matching CreateAgentDTO. */
export const createAgentSchema = z
  .object({
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    instructions: z.string().max(1000),
    llm_ids: z.array(z.uuid()).min(1),
  })
  .strict();

export type CreateAgent = z.infer<typeof createAgentSchema>;

export const updateAgentSchema = createAgentSchema.partial().strict();

export type UpdateAgent = z.infer<typeof updateAgentSchema>;
