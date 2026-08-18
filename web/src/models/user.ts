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

export const userSchema: z.ZodType<User> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      name: z.string(),
      email: z.email(),
      picture: z.string().nullable().optional(),
      created_at: z.coerce.date(),
      updated_at: z.coerce.date(),
      members: z.array(memberSchema).nullable(),
      accounts: z.array(accountSchema).nullable(),
      tokens: z.array(tokenSchema).nullable(),
    })
    .strict(),
);
