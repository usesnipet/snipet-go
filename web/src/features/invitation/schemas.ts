import { z } from "zod";

import { invitationBaseSchema, invitationSchema, invitationStatusSchema } from "@/models/invitation";
import { tenantSchema } from "@/models/tenant";
import { paginatedSchema } from "@/schemas/paginated";

export { invitationSchema, invitationStatusSchema } from "@/models/invitation";
export type { Invitation, InvitationStatus } from "@/models/invitation";

export const paginatedInvitationSchema = paginatedSchema(invitationSchema);
export type PaginatedInvitation = z.infer<typeof paginatedInvitationSchema>;

// "expired" is a filter-only pseudo-status: invitations that are still
// pending but whose expiration date has passed. It has no corresponding
// InvitationStatus, since the API derives it from expires_at rather than
// storing it.
export const invitationStatusFilterSchema = z.enum([
  ...invitationStatusSchema.options,
  "expired",
]);
export type InvitationStatusFilter = z.infer<typeof invitationStatusFilterSchema>;

export const listInvitationSearchParamsSchema = z
  .object({
    status: invitationStatusFilterSchema.optional(),
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
  })
  .strict();

export type ListInvitationSearchParams = z.infer<typeof listInvitationSearchParamsSchema>;

export const createInvitationSchema = invitationBaseSchema.pick({
  email: true,
  role: true,
}).strict();

export type CreateInvitation = z.infer<typeof createInvitationSchema>;

export const acceptInvitationSchema = z
  .object({
    token: z.string().min(1),
  })
  .strict();

export type AcceptInvitation = z.infer<typeof acceptInvitationSchema>;

export const declineInvitationSchema = z
  .object({
    token: z.string().min(1),
  })
  .strict();

export type DeclineInvitation = z.infer<typeof declineInvitationSchema>;

export const invitationInfoSchema = z
  .object({
    invite: invitationSchema,
    tenant: tenantSchema,
  })
  .strict();

export type InvitationInfo = z.infer<typeof invitationInfoSchema>;
