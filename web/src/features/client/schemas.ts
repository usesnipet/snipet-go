import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const clientConfigSchema = z.object({
  oidc: z.object({
    issuer: z.url(),
    audience: z.url(),
    enabled: z.boolean(),
  }).strict(),
  webhook: z.object({
    url: z.url(),
    enabled: z.boolean(),
  }).strict(),
  anonymous: z.object({
    enabled: z.boolean(),
  }).strict(),
}).strict();

export type ClientConfig = z.infer<typeof clientConfigSchema>;

export const clientSchema = z.object({
  id: z.string(),
  code: z.string(),
  name: z.string().min(1).max(255),
  config: clientConfigSchema,
}).strict();

export type Client = z.infer<typeof clientSchema>;

export const paginatedClientSchema = paginatedSchema(clientSchema);
export type PaginatedClient = z.infer<typeof paginatedClientSchema>;

export const createClientSchema = clientSchema.pick({
  name: true,
  config: true,
}).strict();

export type CreateClient = z.infer<typeof createClientSchema>;

export const updateClientSchema = clientSchema.pick({
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
