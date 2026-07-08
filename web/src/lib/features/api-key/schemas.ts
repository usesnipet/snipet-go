import z from "zod";

export const apiKeySchema = z.object({
  id: z.string(),
  name: z.string(),
  key_id: z.string(),
  active: z.boolean(),
  expires_at: z.coerce.date().nullable(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
});

export type APIKey = z.infer<typeof apiKeySchema>;

export const apiKeyPaginatedSchema = z.object({
  data: apiKeySchema.array(),
  total: z.number(),
  take: z.number(),
  skip: z.number(),
});
export type APIKeyPaginated = z.infer<typeof apiKeyPaginatedSchema>;

export const apiKeyWithSecretSchema = apiKeySchema.extend({
  key: z.string(),
});
export type APIKeyWithSecret = z.infer<typeof apiKeyWithSecretSchema>;

export const checkApiKeySchema = z.object({ apiKey: z.string().trim().min(30).startsWith("sn_") });
export type CheckApiKey = z.infer<typeof checkApiKeySchema>;

export const createApiKeySchema = z.object({
  name: z.string().min(1).max(255),
  expires_at: z.date().optional(),
});
export type CreateApiKey = z.infer<typeof createApiKeySchema>;

export const updateApiKeyExpirationSchema = z.object({
  expires_at: z.date().optional(),
});
export type UpdateApiKeyExpiration = z.infer<typeof updateApiKeyExpirationSchema>;