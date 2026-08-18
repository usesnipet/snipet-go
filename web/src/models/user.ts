import { z } from "zod";

import { accountSchema, type Account } from "@/models/account";
import { memberSchema, type Member } from "@/models/member";
import { tokenSchema, type Token } from "@/models/token";

export const challengeSchema = z.enum(["active_account", "change_password"]);
export type Challenge = z.infer<typeof challengeSchema>;

export interface User {
  id: string;
  name: string;
  email: string;
  picture?: string | null;
  created_at: Date;
  updated_at: Date;
  members: Member[] | null;
  accounts: Account[] | null;
  tokens: Token[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const userBaseSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    email: z.email().max(255),
    picture: z.string().nullable().optional(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export const userSchema: z.ZodType<User> = z.lazy(() =>
  userBaseSchema
    .extend({
      members: z.array(memberSchema).nullable(),
      accounts: z.array(accountSchema).nullable(),
      tokens: z.array(tokenSchema).nullable(),
    })
    .strict(),
);
