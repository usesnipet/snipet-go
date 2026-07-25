import http from "@/lib/http";

import { appConfigSchema, systemInfoSchema } from "./schemas";

import type { AppConfig, SystemInfo } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";

const APP_URL = "/api/app";

export const getSystemInfo = async (opts?: ServiceGetOptions<SystemInfo>) => {
  return http.get({
    url: `${APP_URL}/system-info`,
    schemas: {
      response: systemInfoSchema,
    },
    ...opts,
  })
}

export const getAppConfig = async (opts?: ServiceGetOptions<AppConfig>) => {
  return http.get({
    url: `${APP_URL}/config`,
    schemas: {
      response: appConfigSchema,
    },
    ...opts,
  })
}