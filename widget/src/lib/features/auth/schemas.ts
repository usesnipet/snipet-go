import z from "zod";

export const userSchema = z.object({
	id: z.uuid(),
	name: z.string().min(1),
	picture: z.string().nullable(),
	email: z.string().nullable(),
	metadata: z.record(z.string(), z.any()).default({}),
});
export type User = z.infer<typeof userSchema>;

export const authenticateResponseSchema = z.object({
	access_token: z.string().min(1),
	refresh_token: z.string().min(1),
	user: userSchema,
});
export type AuthenticateResponse = z.infer<typeof authenticateResponseSchema>;

export const authenticateAnonymousSchema = z.object({
	name: z.string().max(255).optional(),
	picture: z.url().optional(),
	email: z.email().optional(),
	metadata: z.record(z.string(), z.any()).optional().default({}),
});
export type AuthenticateAnonymous = z.input<typeof authenticateAnonymousSchema>;

export const refreshSchema = z.object({
	refresh_token: z.string().min(1),
});
export type Refresh = z.input<typeof refreshSchema>;
