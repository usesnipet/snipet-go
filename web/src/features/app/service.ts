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

const appsUrl = () => "/api/apps";

const list = async (
  opts?: ServiceGetOptions<PaginatedApp>,
): Promise<PaginatedApp> => {
  return http.get({
    url: appsUrl(),
    schemas: {
      response: paginatedAppSchema,
    },
    ...opts,
  })
}

const create = async (
  body: CreateApp,
  opts?: ServicePostOptions<CreateApp, AppWithSecret>,
): Promise<AppWithSecret> => {
  return http.post({
    url: appsUrl(),
    body,
    schemas: {
      body: createAppSchema,
      response: appWithSecretSchema,
    },
    ...opts,
  })
}

const findByCode = async (
  code: string,
  opts?: ServiceGetOptions<App>,
): Promise<App> => {
  return http.get({
    url: `${appsUrl()}/{code}`,
    params: { code },
    schemas: { response: appSchema },
    ...opts,
  })
}

const update = async (
  code: string,
  body: UpdateApp,
  opts: ServicePutOptions<UpdateApp, void>,
): Promise<void> => {
  return http.put({
    url: `${appsUrl()}/{code}`,
    params: { code },
    body,
    schemas: {
      body: updateAppSchema,
    },
    ...opts,
  })
}

const updateAuthConfig = async (
  code: string,
  body: UpdateAppAuthConfig,
  opts: ServicePutOptions<UpdateAppAuthConfig, void>,
): Promise<void> => {
  return http.put({
    url: `${appsUrl()}/{code}/auth-config`,
    params: { code },
    body,
    schemas: {
      body: updateAppAuthConfigSchema,
    },
    ...opts,
  })
}

const roll = async (
  code: string,
  opts?: ServicePostOptions<undefined, AppWithSecret>,
): Promise<AppWithSecret> => {
  return http.post({
    url: `${appsUrl()}/{code}/roll`,
    params: { code },
    schemas: {
      response: appWithSecretSchema,
    },
    ...opts,
  })
}

const setActive = async (
  code: string,
  active: boolean,
  opts: ServicePutOptions<void, void> = {},
): Promise<void> => {
  return http.put({
    url: `${appsUrl()}/{code}/${active ? "active" : "disabled"}`,
    params: { code },
    ...opts,
  })
}

const remove = async (
  code: string,
  opts: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${appsUrl()}/{code}`,
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
  roll,
  setActive,
  delete: remove,
}
