import { z } from "zod";

import { clientUserSchema, type ClientUser } from "@/models/client-user";

export interface ClientUserRefreshToken {
  id: string;
  client_user_id: string;
  expires_at: Date;
  created_at: Date;
  revoked_at: Date | null;
  metadata: Record<string, unknown>;
  client_user?: ClientUser | null;
}

export const clientUserRefreshTokenSchema: z.ZodType<ClientUserRefreshToken> = z
  .object({
    id: z.uuid(),
    client_user_id: z.uuid(),
    expires_at: z.coerce.date(),
    created_at: z.coerce.date(),
    revoked_at: z.coerce.date().nullable(),
    metadata: z.record(z.string(), z.unknown()),
    client_user: clientUserSchema.nullable().optional(),
  })
  .strict();
