import z from "zod";

export const appConfigSchema = z.object({
	inherit_client: z.boolean(),
});
export type AppConfig = z.infer<typeof appConfigSchema>;
