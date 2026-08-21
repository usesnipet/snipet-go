import { memberBaseSchema, memberSchema } from "@/models/member";
import { memberRoleSchema } from "@/models/member-role";
import { userBaseSchema } from "@/models/user";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export type { MemberRole } from "@/models/member-role";

export const paginatedMemberSchema = paginatedSchema(memberSchema);
export type PaginatedMember = z.infer<typeof paginatedMemberSchema>;

export const listMemberSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
  })
  .strict();

export type ListMemberSearchParams = z.infer<typeof listMemberSearchParamsSchema>;

export const updateMemberRoleSchema = memberBaseSchema.pick({
  role: true,
}).strict();

export type UpdateMemberRole = z.infer<typeof updateMemberRoleSchema>;

// Only accepted by the API on unlicensed single-tenant instances — the
// stand-in for self-registration when there's no invitation flow to reach a
// tenant through (see useGetSystemInfo().multi_tenant_enabled).
export const createMemberSchema = userBaseSchema.pick({
  name: true,
  email: true,
}).extend({
  password: z.string().min(8),
  confirm_password: z.string().min(8),
  role: memberRoleSchema,
}).strict()
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });

export type CreateMember = z.infer<typeof createMemberSchema>;
