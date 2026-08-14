import { clientUserSchema } from "@/models/client-user";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export {
  clientUserMetadataSchema as userMetadataSchema,
  clientUserSchema as userSchema,
} from "@/models/client-user";
export type {
  ClientUser as User,
  ClientUserMetadata as UserMetadata,
} from "@/models/client-user";

export const paginatedUserSchema = paginatedSchema(clientUserSchema);
export type PaginatedUser = z.infer<typeof paginatedUserSchema>;

export const listUserSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
    name_order: z.enum(["asc", "desc"]).optional(),
  })
  .strict();

export type ListUserSearchParams = z.infer<typeof listUserSearchParamsSchema>;
