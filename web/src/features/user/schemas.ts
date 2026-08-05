import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const userMetadataSchema = z.record(z.string(), z.unknown());
export type UserMetadata = z.infer<typeof userMetadataSchema>;

export const userSchema = z
  .object({
    id: z.uuid(),
    name: z.string(),
    picture: z.string().nullable().optional(),
    email: z.string().nullable().optional(),
    metadata: userMetadataSchema,
  })
  .strict();

export type User = z.infer<typeof userSchema>;

export const paginatedUserSchema = paginatedSchema(userSchema);
export type PaginatedUser = z.infer<typeof paginatedUserSchema>;

export const listUserSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
    name_order: z.enum(["asc", "desc"]).optional(),
  })
  .strict();

export type ListUserSearchParams = z.infer<typeof listUserSearchParamsSchema>;
