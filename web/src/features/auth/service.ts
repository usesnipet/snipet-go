import { http } from "@/lib/http";

import {
  activateAccountSchema, authenticateResponseSchema, authorizationUrlResponseSchema, forgotPasswordSchema,
  loginSchema, providerCallbackQuerySchema, refreshSchema, registerResponseSchema, registerSchema,
  resendActivationSchema, resetPasswordSchema, setPasswordSchema
} from "./schemas";

import type {
  ActivateAccount,
  AuthenticateResponse,
  AuthProvider,
  AuthorizationUrlResponse,
  ForgotPassword,
  Login,
  ProviderCallbackQuery,
  Refresh,
  Register,
  RegisterResponse,
  ResendActivation,
  ResetPassword,
  SetPassword,
} from "./schemas";
import type {
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";

const AUTH_URL = "/api/auth";

const register = async (
  body: Register,
  opts: ServicePostOptions<Register, RegisterResponse> = {},
): Promise<RegisterResponse> => {
  return http.post({
    url: `${AUTH_URL}/register`,
    body,
    schemas: {
      body: registerSchema,
      response: registerResponseSchema,
    },
    ...opts,
  });
};

const login = async (
  body: Login,
  opts: ServicePostOptions<Login, AuthenticateResponse> = {},
): Promise<AuthenticateResponse> => {
  return http.post({
    url: `${AUTH_URL}/login`,
    body,
    schemas: {
      body: loginSchema,
      response: authenticateResponseSchema,
    },
    ...opts,
  });
};

const getAuthorizationUrl = async (
  provider: AuthProvider,
  opts: ServiceGetOptions<AuthorizationUrlResponse> = {},
): Promise<AuthorizationUrlResponse> => {
  return http.get({
    url: `${AUTH_URL}/{provider}`,
    params: { provider },
    schemas: {
      response: authorizationUrlResponseSchema,
    },
    ...opts,
  });
};

const callback = async (
  provider: AuthProvider,
  searchParams: ProviderCallbackQuery,
  opts: ServiceGetOptions<AuthenticateResponse, ProviderCallbackQuery> = {},
): Promise<AuthenticateResponse> => {
  return http.get({
    url: `${AUTH_URL}/{provider}/callback`,
    params: { provider },
    searchParams,
    schemas: {
      searchParams: providerCallbackQuerySchema,
      response: authenticateResponseSchema,
    },
    ...opts,
  });
};

const refresh = async (
  body: Refresh,
  opts: ServicePostOptions<Refresh, AuthenticateResponse> = {},
): Promise<AuthenticateResponse> => {
  return http.post({
    url: `${AUTH_URL}/refresh`,
    body,
    schemas: {
      body: refreshSchema,
      response: authenticateResponseSchema,
    },
    ...opts,
  });
};

const logout = async (
  body: Refresh,
  opts: ServicePostOptions<Refresh, void> = {},
): Promise<void> => {
  return http.post({
    url: `${AUTH_URL}/logout`,
    body,
    schemas: {
      body: refreshSchema,
    },
    ...opts,
  });
};

const setPassword = async (
  body: SetPassword,
  opts: ServicePutOptions<SetPassword, void> = {},
): Promise<void> => {
  return http.put({
    url: `${AUTH_URL}/password`,
    body,
    schemas: {
      body: setPasswordSchema,
    },
    ...opts,
  });
};

const forgotPassword = async (
  body: ForgotPassword,
  opts: ServicePostOptions<ForgotPassword, void> = {},
): Promise<void> => {
  return http.post({
    url: `${AUTH_URL}/password/forgot`,
    body,
    schemas: {
      body: forgotPasswordSchema,
    },
    ...opts,
  });
};

const resetPassword = async (
  body: ResetPassword,
  opts: ServicePostOptions<ResetPassword, void> = {},
): Promise<void> => {
  return http.post({
    url: `${AUTH_URL}/password/reset`,
    body,
    schemas: {
      body: resetPasswordSchema,
    },
    ...opts,
  });
};

const activate = async (
  body: ActivateAccount,
  opts: ServicePostOptions<ActivateAccount, void> = {},
): Promise<void> => {
  return http.post({
    url: `${AUTH_URL}/activate`,
    body,
    schemas: {
      body: activateAccountSchema,
    },
    ...opts,
  });
};

const resendActivation = async (
  body: ResendActivation,
  opts: ServicePostOptions<ResendActivation, void> = {},
): Promise<void> => {
  return http.post({
    url: `${AUTH_URL}/activate/resend`,
    body,
    schemas: {
      body: resendActivationSchema,
    },
    ...opts,
  });
};

export const authService = {
  register,
  login,
  getAuthorizationUrl,
  callback,
  refresh,
  logout,
  setPassword,
  forgotPassword,
  resetPassword,
  activate,
  resendActivation,
};
