import { userSchema } from "@/models/user";
import { z } from "zod";

export const authProviderSchema = z.enum(["google", "github"]);
export type AuthProvider = z.infer<typeof authProviderSchema>;

export type { Challenge, User as AuthUser } from "@/models/user";
export { challengeSchema } from "@/models/user";

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

export const registerResponseSchema = z
  .object({
    user: userSchema,
  })
  .strict();

export type RegisterResponse = z.infer<typeof registerResponseSchema>;

export const registerSchema = z
  .object({
    name: z.string().min(1).max(255),
    email: z.email().max(255),
    password: z.string().min(8),
  })
  .strict();

export type Register = z.infer<typeof registerSchema>;

export const loginSchema = z
  .object({
    email: z.email().max(255),
    password: z.string().min(1),
  })
  .strict();

export type Login = z.infer<typeof loginSchema>;

export const refreshSchema = z
  .object({
    refresh_token: z.string().min(1),
  })
  .strict();

export type Refresh = z.infer<typeof refreshSchema>;

export const setPasswordSchema = z
  .object({
    new_password: z.string().min(8),
  })
  .strict();

export type SetPassword = z.infer<typeof setPasswordSchema>;

export const forgotPasswordSchema = z
  .object({
    email: z.email().max(255),
  })
  .strict();

export type ForgotPassword = z.infer<typeof forgotPasswordSchema>;

export const resetPasswordSchema = z
  .object({
    token: z.string().min(1),
    new_password: z.string().min(8),
  })
  .strict();

export type ResetPassword = z.infer<typeof resetPasswordSchema>;

export const activateAccountSchema = z
  .object({
    token: z.string().min(1),
  })
  .strict();

export type ActivateAccount = z.infer<typeof activateAccountSchema>;

export const resendActivationSchema = z
  .object({
    email: z.email().max(255),
  })
  .strict();

export type ResendActivation = z.infer<typeof resendActivationSchema>;

export const authorizationUrlResponseSchema = z
  .object({
    url: z.url(),
  })
  .strict();

export type AuthorizationUrlResponse = z.infer<typeof authorizationUrlResponseSchema>;

export const providerCallbackQuerySchema = z
  .object({
    code: z.string().min(1),
    state: z.string().min(1),
  })
  .strict();

export type ProviderCallbackQuery = z.infer<typeof providerCallbackQuerySchema>;
