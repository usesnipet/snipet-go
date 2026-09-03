import { z } from "zod";

import { agentSchema, type Agent } from "@/models/agent";
import { appToAppUserSchema, type AppToAppUser } from "@/models/app-user";
import { sessionSchema, type Session } from "@/models/session";

/** Accepts "" from form inputs and coerces to undefined; input/output stay `string | undefined`. */
const urlSchema = z.union([
  z.url(),
  z.literal("").transform(() => undefined),
  z.undefined(),
]);

export const appAuthConfigSchema = z
  .object({
    oidc: z
      .object({
        issuer: urlSchema,
        audience: urlSchema,
        enabled: z.boolean(),
      })
      .strict(),
    webhook: z
      .object({
        url: urlSchema,
        enabled: z.boolean(),
      })
      .strict(),
    anonymous: z
      .object({
        enabled: z.boolean(),
      })
      .strict(),
  })
  .strict();

export type AppAuthConfig = z.infer<typeof appAuthConfigSchema>;

export const appStatusSchema = z.enum(["pending", "active", "deactivated"]);
export type AppStatus = z.infer<typeof appStatusSchema>;

export interface App {
  id: string;
  code: string;
  name: string;
  description: string;
  status: AppStatus;
  public: boolean;
  last_verified_at: Date | null;
  key_id: string;
  auth_config: AppAuthConfig;
  sessions: Session[] | null;
  app_to_users: AppToAppUser[] | null;
  app_to_agents: AppToAgent[] | null;
}

export interface AppToAgent {
  app_id: string;
  agent_id: string;
  agent?: Agent | null;
  app?: App | null;
}

export const appToAgentSchema: z.ZodType<AppToAgent> = z.lazy(() =>
  z
    .object({
      app_id: z.uuid(),
      agent_id: z.uuid(),
      agent: agentSchema.nullable().optional(),
      app: appSchema.nullable().optional(),
    })
    .strict(),
);

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const appBaseSchema = z
  .object({
    id: z.string(),
    code: z.string(),
    name: z.string().min(1).max(255),
    description: z.string().max(1000),
    status: appStatusSchema,
    public: z.boolean(),
    last_verified_at: z.coerce.date().nullable(),
    key_id: z.string(),
    auth_config: appAuthConfigSchema,
  })
  .strict();

export const appSchema: z.ZodType<App> = z.lazy(() =>
  appBaseSchema
    .extend({
      sessions: z.array(sessionSchema).nullable(),
      app_to_users: z.array(appToAppUserSchema).nullable(),
      app_to_agents: z.array(appToAgentSchema).nullable(),
    })
    .strict(),
);
