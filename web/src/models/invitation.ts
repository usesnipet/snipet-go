import { z } from "zod";

import { roleSchema } from "@/models/member";

export const invitationStatusSchema = z.enum(["pending", "accepted", "declined"]);
export type InvitationStatus = z.infer<typeof invitationStatusSchema>;

export const invitationSchema = z
  .object({
    id: z.uuid(),
    tenant_id: z.uuid(),
    email: z.email(),
    role: roleSchema,
    status: invitationStatusSchema,
    expires_at: z.coerce.date(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export type Invitation = z.infer<typeof invitationSchema>;
