import { clientBaseSchema, clientSchema } from "@/models/client";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export { clientConfigSchema, clientSchema } from "@/models/client";
export type { Client, ClientConfig } from "@/models/client";

export const paginatedClientSchema = paginatedSchema(clientSchema);
export type PaginatedClient = z.infer<typeof paginatedClientSchema>;

export const createClientSchema = clientBaseSchema.pick({
  name: true,
  config: true,
}).strict();

export type CreateClient = z.infer<typeof createClientSchema>;

export const updateClientSchema = clientBaseSchema.pick({
  name: true,
  config: true,
}).partial().strict();

export type UpdateClient = z.infer<typeof updateClientSchema>;

export const clientPublicSchema = z.object({
  code: z.string(),
  name: z.string(),
  allow_anonymous: z.boolean(),
}).strict();

export type ClientPublic = z.infer<typeof clientPublicSchema>;
