import { http } from "@/lib/http";

import {
  authenticateAnonymousSchema, authenticateResponseSchema, refreshSchema
} from "./schemas";

import type {
  AuthenticateAnonymous, AuthenticateResponse, AuthProvider, Refresh
} from "./schemas";
import type { ServicePostOptions } from "@/lib/services";

const authUrl = (clientCode: string) => `/api/client/${clientCode}`;

const authenticate = async (
  clientCode: string,
  provider: AuthProvider,
  opts: ServicePostOptions<unknown, AuthenticateResponse> = {},
): Promise<AuthenticateResponse> => {
  return http.post({
    url: `${authUrl(clientCode)}/authenticate/{provider}`,
    params: { provider },
    schemas: {
      response: authenticateResponseSchema,
    },
    ...opts,
  })
}

const authenticateAnonymous = async (
  clientCode: string,
  body: AuthenticateAnonymous = {},
  opts: ServicePostOptions<AuthenticateAnonymous, AuthenticateResponse> = {},
): Promise<AuthenticateResponse> => {
  return http.post({
    url: `${authUrl(clientCode)}/authenticate/anonymous`,
    body,
    schemas: {
      body: authenticateAnonymousSchema,
      response: authenticateResponseSchema,
    },
    ...opts,
  })
}

const refresh = async (
  clientCode: string,
  body: Refresh,
  opts: ServicePostOptions<Refresh, AuthenticateResponse> = {},
): Promise<AuthenticateResponse> => {
  return http.post({
    url: `${authUrl(clientCode)}/refresh`,
    body,
    schemas: {
      body: refreshSchema,
      response: authenticateResponseSchema,
    },
    ...opts,
  })
}

export const authService = {
  authenticate,
  authenticateAnonymous,
  refresh,
}
