import { z } from "zod";

export const challengeSchema = z.enum(["active_account", "change_password"]);
export type Challenge = z.infer<typeof challengeSchema>;

export const userSchema = z
  .object({
    id: z.uuid(),
    name: z.string(),
    email: z.email(),
    picture: z.string().nullable().optional(),
    is_admin: z.boolean(),
    challenges: z.array(challengeSchema),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export type User = z.infer<typeof userSchema>;

export const updateProfilePictureSchema = z
  .object({
    picture: z.string().min(1).max(255),
  })
  .strict();

export type UpdateProfilePicture = z.infer<typeof updateProfilePictureSchema>;
