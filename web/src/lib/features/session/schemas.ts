import z from "zod";

export const sessionSchema = z.object({
	id: z.uuid(),
	client_id: z.uuid(),
	agent_id: z.uuid(),
	metadata: z.record(z.string(), z.any()).default({}),
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
	})
	.default({});
export type FilterSession = z.infer<typeof filterSessionSchema>;

export const createSessionSchema = z.object({
	agent_id: z.uuid(),
	metadata: z.record(z.string(), z.string()).optional().default({}),
});
export type CreateSession = z.infer<typeof createSessionSchema>;
