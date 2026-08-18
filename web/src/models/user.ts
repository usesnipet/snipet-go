import { z } from "zod";



export const challengeSchema = z.enum(["active_account", "change_password"]);
export type Challenge = z.infer<typeof challengeSchema>;

export const userSchema = z
  .object({
    id: z.uuid(),
    name: z.string(),
    email: z.email(),
    picture: z.string().nullable().optional(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
    members: z.unknown().nullable().optional(),
  })
  .strict();

export type User = z.infer<typeof userSchema>;
