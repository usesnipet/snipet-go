import { tenantSchema } from "@/models/tenant";
import { userSchema } from "@/models/user";
import { z } from "zod";

import type { Tenant } from "@/models/tenant";
import type { User } from "@/models/user";
import { memberRoleSchema, type MemberRole } from "./member-role";

export interface Member {
  id: string;
  user_id: string;
  tenant_id: string;
  role: MemberRole;
  is_active: boolean;
  created_at: Date;
  updated_at: Date;
  user?: User | null;
  tenant?: Tenant | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const memberBaseSchema = z
  .object({
    id: z.uuid(),
    user_id: z.uuid(),
    tenant_id: z.uuid(),
    role: memberRoleSchema,
    is_active: z.boolean(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export const memberSchema: z.ZodType<Member> = z.lazy(() =>
  memberBaseSchema
    .extend({
      user: userSchema.nullable(),
      tenant: tenantSchema.nullable(),
    })
    .strict(),
);
