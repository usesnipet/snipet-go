import z from "zod";

export const clientConfigSchema = z.object({
	oidc: z
		.object({
			enabled: z.boolean().default(false),
			issuer: z.string().default(""),
			audience: z.string().default(""),
		})
		.default({ enabled: false, issuer: "", audience: "" }),
	webhook: z
		.object({
			enabled: z.boolean().default(false),
			url: z.string().default(""),
		})
		.default({ enabled: false, url: "" }),
});

export type ClientConfig = z.infer<typeof clientConfigSchema>;

export const clientSchema = z.object({
	id: z.uuid(),
	code: z.string().min(1),
	name: z.string().min(1).max(255),
	config: clientConfigSchema,
});
export type Client = z.infer<typeof clientSchema>;

export const clientPaginatedSchema = z.object({
	data: clientSchema.array(),
	total: z.number(),
	take: z.number(),
	skip: z.number(),
});
export type ClientPaginated = z.infer<typeof clientPaginatedSchema>;

export const filterClientSchema = z
	.object({
		take: z.number().min(1).optional(),
		skip: z.number().min(0).optional(),
	})
	.default({});
export type FilterClient = z.infer<typeof filterClientSchema>;

export const createClientSchema = z.object({
	name: z.string().min(1).max(255),
	config: clientConfigSchema.default({
		oidc: { enabled: false, issuer: "", audience: "" },
		webhook: { enabled: false, url: "" },
	}),
});
export type CreateClient = z.infer<typeof createClientSchema>;

export const updateClientSchema = z.object({
	name: z.string().min(1).max(255),
	config: clientConfigSchema,
});
export type UpdateClient = z.infer<typeof updateClientSchema>;
