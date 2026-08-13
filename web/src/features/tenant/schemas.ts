import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const tenantSchema = z.object({
  id: z.string(),
  name: z.string().min(1).max(255),
  slug: z.string().min(1).max(255),
  icon: z.string().max(255).nullable(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
}).strict();

export type Tenant = z.infer<typeof tenantSchema>;

export const paginatedTenantSchema = paginatedSchema(tenantSchema);
export type PaginatedTenant = z.infer<typeof paginatedTenantSchema>;

export const createTenantSchema = tenantSchema.pick({
  name: true,
  slug: true,
}).extend({
  icon: z.string().max(255).optional(),
}).strict();

export type CreateTenant = z.infer<typeof createTenantSchema>;

export const updateTenantSchema = tenantSchema.pick({
  name: true,
  slug: true,
  icon: true,
}).partial().strict();

export type UpdateTenant = z.infer<typeof updateTenantSchema>;
