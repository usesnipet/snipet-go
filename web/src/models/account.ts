import { z } from "zod";

import { userSchema, type User } from "@/models/user";

export interface Account {
  id: string;
  user_id: string;
  provider: string;
  external_id: string;
  created_at: Date;
  updated_at: Date;
  user?: User | null;
}

export const accountSchema: z.ZodType<Account> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      user_id: z.uuid(),
      provider: z.string(),
      external_id: z.string(),
      created_at: z.coerce.date(),
      updated_at: z.coerce.date(),
      user: userSchema.nullable().optional(),
    })
    .strict(),
);
