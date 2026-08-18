import { z } from "zod";

import { clientSchema, type Client } from "@/models/client";
import { sessionSchema, type Session } from "@/models/session";

export const clientUserMetadataSchema = z.record(z.string(), z.unknown());
export type ClientUserMetadata = z.infer<typeof clientUserMetadataSchema>;

export interface ClientUser {
  id: string;
  name: string;
  picture?: string | null;
  email?: string | null;
  metadata: ClientUserMetadata;
  client_user_to_sessions: ClientUserToSession[] | null;
  client_to_client_users: ClientToClientUser[] | null;
}

export const clientUserSchema: z.ZodType<ClientUser> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      name: z.string(),
      picture: z.string().nullable().optional(),
      email: z.string().nullable().optional(),
      metadata: clientUserMetadataSchema,
      client_user_to_sessions: z.array(clientUserToSessionSchema).nullable(),
      client_to_client_users: z.array(clientToClientUserSchema).nullable(),
    })
    .strict(),
);

export interface ClientToClientUser {
  client_id: string;
  client_user_id: string;
  external_id: string | null;
  client?: Client | null;
  client_user?: ClientUser | null;
}

export const clientToClientUserSchema: z.ZodType<ClientToClientUser> = z.lazy(() =>
  z
    .object({
      client_id: z.uuid(),
      client_user_id: z.uuid(),
      external_id: z.string().nullable(),
      client: clientSchema.nullable().optional(),
      client_user: clientUserSchema.nullable().optional(),
    })
    .strict(),
);

export interface ClientUserToSession {
  user_id: string;
  session_id: string;
  client_user?: ClientUser | null;
  session?: Session | null;
}

export const clientUserToSessionSchema: z.ZodType<ClientUserToSession> = z.lazy(() =>
  z
    .object({
      user_id: z.uuid(),
      session_id: z.uuid(),
      client_user: clientUserSchema.nullable().optional(),
      session: sessionSchema.nullable().optional(),
    })
    .strict(),
);
