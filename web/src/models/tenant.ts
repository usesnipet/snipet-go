import { z } from "zod";

export const tenantSchema = z
  .object({
    id: z.string(),
    name: z.string().min(1).max(255),
    slug: z.string().min(1).max(255),
    icon: z.string().max(255).nullable(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export type Tenant = z.infer<typeof tenantSchema>;
