import { z } from "zod";

import { appSchema, type App } from "@/models/app";
import { sessionSchema, type Session } from "@/models/session";

export const appUserMetadataSchema = z.record(z.string(), z.unknown());
export type AppUserMetadata = z.infer<typeof appUserMetadataSchema>;

export interface AppUser {
  id: string;
  name: string;
  picture?: string | null;
  email?: string | null;
  metadata: AppUserMetadata;
  app_user_to_sessions: AppUserToSession[] | null;
  app_to_App_users: AppToAppUser[] | null;
}

export const appUserSchema: z.ZodType<AppUser> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      name: z.string(),
      picture: z.string().nullable().optional(),
      email: z.string().nullable().optional(),
      metadata: appUserMetadataSchema,
      app_user_to_sessions: z.array(appUserToSessionSchema).nullable(),
      app_to_App_users: z.array(appToAppUserSchema).nullable(),
    })
    .strict(),
);

export interface AppToAppUser {
  app_id: string;
  app_user_id: string;
  external_id: string | null;
  app?: App | null;
  app_user?: AppUser | null;
}

export const appToAppUserSchema: z.ZodType<AppToAppUser> = z.lazy(() =>
  z
    .object({
      app_id: z.uuid(),
      app_user_id: z.uuid(),
      external_id: z.string().nullable(),
      app: appSchema.nullable().optional(),
      app_user: appUserSchema.nullable().optional(),
    })
    .strict(),
);

export interface AppUserToSession {
  user_id: string;
  session_id: string;
  app_user?: AppUser | null;
  session?: Session | null;
}

export const appUserToSessionSchema: z.ZodType<AppUserToSession> = z.lazy(() =>
  z
    .object({
      user_id: z.uuid(),
      session_id: z.uuid(),
      app_user: appUserSchema.nullable().optional(),
      session: sessionSchema.nullable().optional(),
    })
    .strict(),
);
