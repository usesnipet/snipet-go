import { paginatedSchema, paginationParamsSchema } from "@/schemas/paginated";
import { z } from "zod";

export const apiKeyKeySchema = z.string().length(46).startsWith("sn_");
export type ApiKeyKey = z.infer<typeof apiKeyKeySchema>;

const optionalDateSchema = z
  .union([z.coerce.date(), z.null()])
  .optional()
  .transform((value) => value ?? undefined);

export const apiKeySchema = z.object({
  id: z.string(),
  name: z.string().min(1).max(255),
  key_id: z.string(),
  active: z.boolean(),
  expires_at: optionalDateSchema,
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
}).strict();

export type ApiKey = z.infer<typeof apiKeySchema>;

export const paginatedApiKeySchema = paginatedSchema(apiKeySchema);
export type PaginatedApiKey = z.infer<typeof paginatedApiKeySchema>;

export const listApiKeySearchParamsSchema = paginationParamsSchema;
export type ListApiKeySearchParams = z.infer<typeof listApiKeySearchParamsSchema>;

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
