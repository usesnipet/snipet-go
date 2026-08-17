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

// Only accepted by the API on unlicensed single-tenant instances — the
// stand-in for self-registration when there's no invitation flow to reach a
// tenant through (see useGetSystemInfo().multi_tenant_enabled).
export const createMemberSchema = z
  .object({
    name: z.string().min(1).max(255),
    email: z.email().max(255),
    password: z.string().min(8),
    confirm_password: z.string().min(8),
    role: roleSchema,
  })
  .strict()
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });

export type CreateMember = z.infer<typeof createMemberSchema>;
