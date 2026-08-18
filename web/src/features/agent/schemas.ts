import { agentBaseSchema, agentSchema } from "@/models/agent";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export { agentSchema, agentToLLMSchema as agentToLlmSchema } from "@/models/agent";
export type { Agent, AgentToLLM as AgentToLlm } from "@/models/agent";

export const paginatedAgentSchema = paginatedSchema(agentSchema);
export type PaginatedAgent = z.infer<typeof paginatedAgentSchema>;

/** Form / create payload shape matching CreateAgentDTO. */
export const createAgentSchema = agentBaseSchema.pick({
  name: true,
  description: true,
  instructions: true,
}).extend({
  llm_ids: z.array(z.uuid()).min(1),
}).strict();

export type CreateAgent = z.infer<typeof createAgentSchema>;

export const updateAgentSchema = createAgentSchema.partial().strict();

export type UpdateAgent = z.infer<typeof updateAgentSchema>;
