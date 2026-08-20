import { http } from "@/lib/http";

import {
  appSchema, appWithSecretSchema, createAppSchema, paginatedAppSchema, updateAppAuthConfigSchema,
  updateAppSchema
} from "./schemas";

import type {
  App, AppWithSecret, CreateApp, PaginatedApp, UpdateApp, UpdateAppAuthConfig
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const appsUrl = (tenantId: string) => `/api/tenants/${tenantId}/apps`;

const list = async (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedApp>,
): Promise<PaginatedApp> => {
  return http.get({
    url: appsUrl(tenantId),
    schemas: {
      response: paginatedAppSchema,
    },
    ...opts,
  })
}

const create = async (
  tenantId: string,
  body: CreateApp,
  opts?: ServicePostOptions<CreateApp, AppWithSecret>,
): Promise<AppWithSecret> => {
  return http.post({
    url: appsUrl(tenantId),
    body,
    schemas: {
      body: createAppSchema,
      response: appWithSecretSchema,
    },
    ...opts,
  })
}

const findByCode = async (
  tenantId: string,
  code: string,
  opts?: ServiceGetOptions<App>,
): Promise<App> => {
  return http.get({
    url: `${appsUrl(tenantId)}/{code}`,
    params: { code },
    schemas: { response: appSchema },
    ...opts,
  })
}

const update = async (
  tenantId: string,
  code: string,
  body: UpdateApp,
  opts: ServicePutOptions<UpdateApp, void>,
): Promise<void> => {
  return http.put({
    url: `${appsUrl(tenantId)}/{code}`,
    params: { code },
    body,
    schemas: {
      body: updateAppSchema,
    },
    ...opts,
  })
}

const updateAuthConfig = async (
  tenantId: string,
  code: string,
  body: UpdateAppAuthConfig,
  opts: ServicePutOptions<UpdateAppAuthConfig, void>,
): Promise<void> => {
  return http.put({
    url: `${appsUrl(tenantId)}/{code}/auth-config`,
    params: { code },
    body,
    schemas: {
      body: updateAppAuthConfigSchema,
    },
    ...opts,
  })
}

const remove = async (
  tenantId: string,
  code: string,
  opts: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${appsUrl(tenantId)}/{code}`,
    params: { code },
    ...opts,
  })
}

export const appService = {
  list,
  create,
  findByCode,
  update,
  updateAuthConfig,
  delete: remove,
}
