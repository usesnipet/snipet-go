import { appAuthConfigSchema, appBaseSchema, appSchema } from "@/models/app";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export { appAuthConfigSchema, appSchema, appStatusSchema, appToAgentSchema } from "@/models/app";
export type { App, AppAuthConfig, AppStatus, AppToAgent } from "@/models/app";

export const paginatedAppSchema = paginatedSchema(appSchema);
export type PaginatedApp = z.infer<typeof paginatedAppSchema>;

export const appWithSecretSchema = appSchema.and(z.object({ key: z.string() }));
export type AppWithSecret = z.infer<typeof appWithSecretSchema>;

/** Writable app fields shared by the create and update payloads. */
const appWritableSchema = appBaseSchema.pick({
  name: true,
  description: true,
  public: true,
}).extend({
  agent_ids: z.array(z.uuid()),
});

export const createAppSchema = appWritableSchema.strict();

export type CreateApp = z.infer<typeof createAppSchema>;

export const updateAppSchema = appWritableSchema.partial().strict();

export type UpdateApp = z.infer<typeof updateAppSchema>;

export const linkAppAgentsSchema = z.object({
  agent_ids: z.array(z.uuid()),
}).strict();

export type LinkAppAgents = z.infer<typeof linkAppAgentsSchema>;

export const updateAppAuthConfigSchema = z.object({
  auth_config: appAuthConfigSchema,
}).strict();

export type UpdateAppAuthConfig = z.infer<typeof updateAppAuthConfigSchema>;
