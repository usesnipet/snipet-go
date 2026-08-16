import { z } from "zod";

import { invitationSchema, invitationStatusSchema } from "@/models/invitation";
import { roleSchema } from "@/models/member";
import { paginatedSchema } from "@/schemas/paginated";

export { invitationSchema, invitationStatusSchema } from "@/models/invitation";
export type { Invitation, InvitationStatus } from "@/models/invitation";

export const paginatedInvitationSchema = paginatedSchema(invitationSchema);
export type PaginatedInvitation = z.infer<typeof paginatedInvitationSchema>;

export const listInvitationSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
  })
  .strict();

export type ListInvitationSearchParams = z.infer<typeof listInvitationSearchParamsSchema>;

export const createInvitationSchema = z
  .object({
    email: z.email().max(255),
    role: roleSchema,
  })
  .strict();

export type CreateInvitation = z.infer<typeof createInvitationSchema>;

export const acceptInvitationSchema = z
  .object({
    token: z.string().min(1),
  })
  .strict();

export type AcceptInvitation = z.infer<typeof acceptInvitationSchema>;
