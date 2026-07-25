import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const apiKeyKeySchema = z.string().length(46).startsWith("sn_");
export type ApiKeyKey = z.infer<typeof apiKeyKeySchema>;

export const apiKeySchema = z.object({
  id: z.string(),
  name: z.string(),
  key_id: z.string(),
  active: z.boolean(),
  expires_at: z.coerce.date().optional(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
}).strict();

export type ApiKey = z.infer<typeof apiKeySchema>;

export const paginatedApiKeySchema = paginatedSchema(apiKeySchema);
export type PaginatedApiKey = z.infer<typeof paginatedApiKeySchema>;

export const createApiKeySchema = apiKeySchema.pick({
  name: true,
  expires_at: true,
}).strict();

export type CreateApiKey = z.infer<typeof createApiKeySchema>;

export const updateApiKeyExpirationSchema = apiKeySchema.pick({
  expires_at: true,
}).strict();

export type UpdateApiKeyExpiration = z.infer<typeof updateApiKeyExpirationSchema>;

export const apiKeyWithSecretSchema = apiKeySchema.extend({
  key: apiKeyKeySchema,
}).strict();

export type ApiKeyWithSecret = z.infer<typeof apiKeyWithSecretSchema>;