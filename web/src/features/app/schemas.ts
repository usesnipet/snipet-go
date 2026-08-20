import { appAuthConfigSchema, appBaseSchema, appSchema } from "@/models/app";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export { appAuthConfigSchema, appSchema, appStatusSchema } from "@/models/app";
export type { App, AppAuthConfig, AppStatus } from "@/models/app";

export const paginatedAppSchema = paginatedSchema(appSchema);
export type PaginatedApp = z.infer<typeof paginatedAppSchema>;

export const appWithSecretSchema = appSchema.and(z.object({ key: z.string() }));
export type AppWithSecret = z.infer<typeof appWithSecretSchema>;

export const createAppSchema = appBaseSchema.pick({
  name: true,
  description: true,
}).strict();

export type CreateApp = z.infer<typeof createAppSchema>;

export const updateAppSchema = appBaseSchema.pick({
  name: true,
  description: true,
}).partial().strict();

export type UpdateApp = z.infer<typeof updateAppSchema>;

export const updateAppAuthConfigSchema = z.object({
  auth_config: appAuthConfigSchema,
}).strict();

export type UpdateAppAuthConfig = z.infer<typeof updateAppAuthConfigSchema>;
