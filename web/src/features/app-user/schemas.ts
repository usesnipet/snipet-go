import { appUserSchema } from "@/models/app-user";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export {
  appUserMetadataSchema as userMetadataSchema,
  appUserSchema as userSchema,
} from "@/models/app-user";
export type {
  AppUser as User,
  AppUserMetadata as UserMetadata,
} from "@/models/app-user";

export const paginatedUserSchema = paginatedSchema(appUserSchema);
export type PaginatedUser = z.infer<typeof paginatedUserSchema>;

export const listUserSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
    name_order: z.enum(["asc", "desc"]).optional(),
  })
  .strict();

export type ListUserSearchParams = z.infer<typeof listUserSearchParamsSchema>;
