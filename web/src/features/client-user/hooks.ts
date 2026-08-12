import { useQuery } from "@tanstack/react-query";

import { clientUserService } from "./service";

import type { ListUserSearchParams, PaginatedUser } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";
import type { UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "client-user";

export const listClientUserQueryKey = (clientCode: string) =>
  [BASE_QUERY_KEY, clientCode] as const;

export const useListClientUser = (
  clientCode: string,
  opts?: ServiceGetOptions<PaginatedUser, ListUserSearchParams>,
): UseQueryResult<PaginatedUser, Error> => {
  return useQuery({
    queryKey: [...listClientUserQueryKey(clientCode), opts?.searchParams],
    queryFn: () =>
      clientUserService.list(clientCode, opts),
    enabled: !!clientCode,
  })
}