import { z } from "zod";

import { llmSchema } from "@/models/llm";
import { tenantSchema, type Tenant } from "@/models/tenant";

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
  tenant_id: string;
  name: string;
  description: string;
  instructions: string;
  llms: AgentToLLM[] | null;
  knowledge: AgentToKnowledge[] | null;
  tenant?: Tenant | null;
}

export const agentSchema: z.ZodType<Agent> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      tenant_id: z.uuid(),
      name: z.string().min(1).max(255),
      description: z.string().max(1000),
      instructions: z.string().max(1000),
      llms: z.array(agentToLLMSchema).nullable(),
      knowledge: z.array(agentToKnowledgeSchema).nullable(),
      tenant: tenantSchema.nullable().optional(),
    })
    .strict(),
);
