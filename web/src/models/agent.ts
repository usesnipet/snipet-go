import { z } from "zod";

import { llmSchema } from "@/models/llm";

export const agentToLLMSchema = z
  .object({
    llm_id: z.uuid(),
    priority: z.number().int(),
    llm: llmSchema,
  })
  .strict();

export type AgentToLLM = z.infer<typeof agentToLLMSchema>;

export const agentToKnowledgeSchema = z
  .object({
    agent_id: z.uuid(),
    knowledge_id: z.uuid(),
    active: z.boolean(),
  })
  .strict();

export type AgentToKnowledge = z.infer<typeof agentToKnowledgeSchema>;

export interface Agent {
  id: string;
  name: string;
  description: string;
  instructions: string;
  llms: AgentToLLM[] | null;
  knowledge: AgentToKnowledge[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const agentBaseSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    instructions: z.string().max(1000),
  })
  .strict();

export const agentSchema: z.ZodType<Agent> = z.lazy(() =>
  agentBaseSchema
    .extend({
      llms: z.array(agentToLLMSchema).nullable(),
      knowledge: z.array(agentToKnowledgeSchema).nullable(),
    })
    .strict(),
);
