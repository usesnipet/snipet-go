import { useQuery } from "@tanstack/react-query";

import { appUserService } from "./service";

import type { ListUserSearchParams, PaginatedUser } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";
import type { UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "app-user";

export const listAppUserQueryKey = (appCode: string) =>
  [BASE_QUERY_KEY, appCode] as const;

export const useListAppUser = (
  appCode: string,
  opts?: ServiceGetOptions<PaginatedUser, ListUserSearchParams>,
): UseQueryResult<PaginatedUser, Error> => {
  return useQuery({
    queryKey: [...listAppUserQueryKey(appCode), opts?.searchParams],
    queryFn: () =>
      appUserService.list(appCode, opts),
    enabled: !!appCode,
  })
}
