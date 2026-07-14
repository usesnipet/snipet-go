import z from "zod";

export const agentConfigurationSchema = z.object({
	llms: z.array(z.any()).default([]).nullable(),
});

export type AgentConfiguration = z.infer<typeof agentConfigurationSchema>;

export const agentSchema = z.object({
	id: z.uuid(),
	name: z.string().min(1).max(255),
	description: z.string().max(1000),
	configuration: agentConfigurationSchema,
});
export type Agent = z.infer<typeof agentSchema>;

export const agentPaginatedSchema = z.object({
	data: agentSchema.array(),
	total: z.number(),
	take: z.number(),
	skip: z.number(),
});
export type AgentPaginated = z.infer<typeof agentPaginatedSchema>;

export const filterAgentSchema = z
	.object({
		take: z.number().min(1).optional(),
		skip: z.number().min(0).optional(),
	})
	.default({});
export type FilterAgent = z.infer<typeof filterAgentSchema>;

export const createAgentSchema = z.object({
	name: z.string().min(1).max(255),
	description: z.string().max(1000).optional().default(""),
	configuration: agentConfigurationSchema,
});
export type CreateAgent = z.infer<typeof createAgentSchema>;

export const updateAgentSchema = z.object({
	name: z.string().min(1).max(255).optional(),
	description: z.string().max(1000).optional(),
});
export type UpdateAgent = z.infer<typeof updateAgentSchema>;
