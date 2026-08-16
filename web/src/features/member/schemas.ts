import { z } from "zod";

import { memberSchema, roleSchema } from "@/models/member";
import { paginatedSchema } from "@/schemas/paginated";

export { memberSchema, roleSchema } from "@/models/member";
export type { Member, Role } from "@/models/member";

export const paginatedMemberSchema = paginatedSchema(memberSchema);
export type PaginatedMember = z.infer<typeof paginatedMemberSchema>;

export const listMemberSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
  })
  .strict();

export type ListMemberSearchParams = z.infer<typeof listMemberSearchParamsSchema>;

export const updateMemberRoleSchema = z
  .object({
    role: roleSchema,
  })
  .strict();

export type UpdateMemberRole = z.infer<typeof updateMemberRoleSchema>;
