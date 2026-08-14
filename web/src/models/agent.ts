import { llmSchema } from "@/models/llm";
import { z } from "zod";

export const agentToLLMSchema = z
  .object({
    llm_id: z.uuid(),
    priority: z.number().int(),
    llm: llmSchema,
  })
  .strict();

export type AgentToLLM = z.infer<typeof agentToLLMSchema>;

export const agentSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    instructions: z.string().max(1000),
    llms: z.array(agentToLLMSchema).nullable(),
  })
  .strict();

export type Agent = z.infer<typeof agentSchema>;
