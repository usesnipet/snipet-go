import { tenantBaseSchema, tenantSchema } from "@/models/tenant";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export { tenantSchema } from "@/models/tenant";
export type { Tenant } from "@/models/tenant";

export const paginatedTenantSchema = paginatedSchema(tenantSchema);
export type PaginatedTenant = z.infer<typeof paginatedTenantSchema>;

export const createTenantSchema = tenantBaseSchema.pick({
  name: true,
  slug: true,
}).extend({
  icon: z.string().max(255).optional(),
}).strict();

export type CreateTenant = z.infer<typeof createTenantSchema>;

export const updateTenantSchema = tenantBaseSchema.pick({
  name: true,
  slug: true,
  icon: true,
}).partial().strict();

export type UpdateTenant = z.infer<typeof updateTenantSchema>;
