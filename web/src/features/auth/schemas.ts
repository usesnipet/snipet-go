import { userMetadataSchema, userSchema } from "@/features/user/schemas";
import { z } from "zod";

export const authProviderSchema = z.enum(["oidc", "webhook"]);
export type AuthProvider = z.infer<typeof authProviderSchema>;

export const authenticateResponseSchema = z
  .object({
    access_token: z.string().min(1),
    access_token_expires_at: z.coerce.date(),
    refresh_token: z.string().min(1),
    refresh_token_expires_at: z.coerce.date(),
    user: userSchema,
  })
  .strict();

export type AuthenticateResponse = z.infer<typeof authenticateResponseSchema>;

export const authTokensSchema = authenticateResponseSchema
  .pick({
    access_token: true,
    access_token_expires_at: true,
    refresh_token: true,
    refresh_token_expires_at: true,
  })
  .strict();

export type AuthTokens = z.infer<typeof authTokensSchema>;

export const authenticateAnonymousSchema = z
  .object({
    name: z.string().max(255).optional(),
    picture: z.url().optional(),
    email: z.email().optional(),
    metadata: userMetadataSchema.optional(),
  })
  .strict();

export type AuthenticateAnonymous = z.infer<typeof authenticateAnonymousSchema>;

export const refreshSchema = z
  .object({
    refresh_token: z.string().min(1),
  })
  .strict();

export type Refresh = z.infer<typeof refreshSchema>;
