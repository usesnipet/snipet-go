import { roleSchema } from "@/models/member";
import { tenantSchema } from "@/models/tenant";
import { z } from "zod";

import type { MemberRole } from "@/models/member";
import type { Tenant } from "@/models/tenant";

export const invitationStatusSchema = z.enum(["pending", "accepted", "declined"]);
export type InvitationStatus = z.infer<typeof invitationStatusSchema>;

export interface Invitation {
  id: string;
  tenant_id: string;
  email: string;
  role: MemberRole;
  status: InvitationStatus;
  expires_at: Date;
  created_at: Date;
  updated_at: Date;
  tenant?: Tenant | null;
}

export const invitationSchema: z.ZodType<Invitation> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      tenant_id: z.uuid(),
      email: z.email(),
      role: roleSchema,
      status: invitationStatusSchema,
      expires_at: z.coerce.date(),
      created_at: z.coerce.date(),
      updated_at: z.coerce.date(),
      tenant: tenantSchema.nullable().optional(),
    })
    .strict(),
);
