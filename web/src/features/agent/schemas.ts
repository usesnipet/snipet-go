import { driverInfoSchema } from "@/schemas/driver";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const llmConfigSchema = z
  .object({
    key: z.string().min(1),
    config: z.record(z.string(), z.unknown()),
  })
  .strict();

export type LlmConfig = z.infer<typeof llmConfigSchema>;

export const toolConfigSchema = z.record(z.string(), z.record(z.string(), z.unknown()));

export type ToolConfig = z.infer<typeof toolConfigSchema>;

export const agentConfigurationSchema = z
  .object({
    llms: z.array(llmConfigSchema).nullable(),
    tools: toolConfigSchema.nullable(),
  })
  .strict();

export type AgentConfiguration = z.infer<typeof agentConfigurationSchema>;

export const agentSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    instructions: z.string().max(1000),
    configuration: agentConfigurationSchema,
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
    llms: z.array(llmConfigSchema).min(1),
    tools: toolConfigSchema,
  })
  .strict();

export type CreateAgent = z.infer<typeof createAgentSchema>;

export const updateAgentSchema = createAgentSchema.partial().strict();

export type UpdateAgent = z.infer<typeof updateAgentSchema>;

export const listDriversSchema = z.array(driverInfoSchema);
export type ListDrivers = z.infer<typeof listDriversSchema>;