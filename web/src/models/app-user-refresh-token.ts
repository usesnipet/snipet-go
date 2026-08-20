import { z } from "zod";

import { appUserSchema, type AppUser } from "@/models/app-user";

export interface AppUserRefreshToken {
  id: string;
  app_user_id: string;
  expires_at: Date;
  created_at: Date;
  revoked_at: Date | null;
  metadata: Record<string, unknown>;
  app_user?: AppUser | null;
}

export const appUserRefreshTokenSchema: z.ZodType<AppUserRefreshToken> = z
  .object({
    id: z.uuid(),
    app_user_id: z.uuid(),
    expires_at: z.coerce.date(),
    created_at: z.coerce.date(),
    revoked_at: z.coerce.date().nullable(),
    metadata: z.record(z.string(), z.unknown()),
    app_user: appUserSchema.nullable().optional(),
  })
  .strict();
