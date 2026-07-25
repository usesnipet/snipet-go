import type { ServiceGetOptions } from "@/lib/services"
import { useQuery } from "@tanstack/react-query";

import { getAppConfig, getSystemInfo } from "./service";

import type { AppConfig, SystemInfo } from "./schemas";

export const useGetSystemInfo = (opts?: ServiceGetOptions<SystemInfo>) => {
  return useQuery({
    queryKey: ["system-info"],
    queryFn: () => getSystemInfo(opts),
  })
}

export const useGetAppConfig = (opts?: ServiceGetOptions<AppConfig>) => {
  return useQuery({
    queryKey: ["app-config"],
    queryFn: () => getAppConfig(opts),
  })
}