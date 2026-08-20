import http from "@/lib/http";

import { systemInfoSchema } from "./schemas";

import type { SystemInfo } from "./schemas";
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