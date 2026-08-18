import { z } from "zod";

export const apiKeyKeySchema = z.string().length(46).startsWith("sn_");
export type ApiKeyKey = z.infer<typeof apiKeyKeySchema>;

const optionalDateSchema = z
  .union([z.coerce.date(), z.null()])
  .optional()
  .transform((value) => value ?? undefined);

export const apiKeySchema = z
  .object({
    id: z.string(),
    tenant_id: z.string(),
    name: z.string().min(1).max(255),
    key_id: z.string(),
    active: z.boolean(),
    expires_at: optionalDateSchema,
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export type ApiKey = z.infer<typeof apiKeySchema>;
