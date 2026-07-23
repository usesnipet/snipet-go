import z from "zod";

export const appConfigSchema = z.object({
	inherit_client: z.boolean(),
	inherit_client_code: z.string(),
	inherit_client_name: z.string(),
});
export type AppConfig = z.infer<typeof appConfigSchema>;
