import z from "zod";

import { agentSchema } from "../agent/schemas";

export const sessionSchema = z.object({
	id: z.uuid(),
	client_id: z.uuid(),
	agent_id: z.uuid(),
	metadata: z.record(z.string(), z.any()).default({}),
	agent: agentSchema.optional(),
});
export type Session = z.infer<typeof sessionSchema>;

export const sessionPaginatedSchema = z.object({
	data: sessionSchema.array(),
	total: z.number(),
	take: z.number(),
	skip: z.number(),
});
export type SessionPaginated = z.infer<typeof sessionPaginatedSchema>;

export const filterSessionSchema = z
	.object({
		take: z.number().min(1).optional(),
		skip: z.number().min(0).optional(),
		include: z.array(z.enum(["agent"])).optional(),
	})
	.default({});
export type FilterSession = z.infer<typeof filterSessionSchema>;

export const filterSessionByIDSchema = z.object({
	include: z.array(z.enum(["agent"])).optional(),
}).default({});
export type FilterSessionByID = z.infer<typeof filterSessionByIDSchema>;

export const createSessionSchema = z.object({
	agent_id: z.uuid(),
	metadata: z.record(z.string(), z.string()).optional().default({}),
});
export type CreateSession = z.infer<typeof createSessionSchema>;
