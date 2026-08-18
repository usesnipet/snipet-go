import { z } from "zod";

import { userSchema, type User } from "@/models/user";

export const tokenTypeSchema = z.enum(["refresh", "activate_account", "reset_password"]);
export type TokenType = z.infer<typeof tokenTypeSchema>;

export interface Token {
  id: string;
  type: TokenType;
  user_id: string;
  expires_at: Date;
  revoked_at: Date | null;
  metadata: Record<string, unknown>;
  created_at: Date;
  user?: User | null;
}

export const tokenSchema: z.ZodType<Token> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      type: tokenTypeSchema,
      user_id: z.uuid(),
      expires_at: z.coerce.date(),
      revoked_at: z.coerce.date().nullable(),
      metadata: z.record(z.string(), z.unknown()),
      created_at: z.coerce.date(),
      user: userSchema.nullable().optional(),
    })
    .strict(),
);
