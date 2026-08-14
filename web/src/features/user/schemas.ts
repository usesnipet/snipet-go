import { z } from "zod";

export { challengeSchema, userSchema } from "@/models/user";
export type { Challenge, User } from "@/models/user";

export const updateProfilePictureSchema = z
  .object({
    picture: z.string().min(1).max(255),
  })
  .strict();

export type UpdateProfilePicture = z.infer<typeof updateProfilePictureSchema>;
