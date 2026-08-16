import { z } from "zod";

import { tenantSchema } from "@/models/tenant";
import { userSchema } from "@/models/user";

export const roleSchema = z.enum(["admin", "user"]);
export type Role = z.infer<typeof roleSchema>;

export const memberSchema = z
  .object({
    id: z.uuid(),
    user_id: z.uuid(),
    tenant_id: z.uuid(),
    role: roleSchema,
    is_active: z.boolean(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
    user: userSchema,
    tenant: tenantSchema,
  })
  .strict();

export type Member = z.infer<typeof memberSchema>;
