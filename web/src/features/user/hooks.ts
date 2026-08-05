import { useQuery } from "@tanstack/react-query";

import { userService } from "./service";

import type { ListUserSearchParams, PaginatedUser, User } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";
import type { UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "user";

export const listUserQueryKey = (clientCode: string) =>
  [BASE_QUERY_KEY, clientCode] as const;

export const useListUser = (
  clientCode: string,
  opts?: ServiceGetOptions<PaginatedUser, ListUserSearchParams>,
): UseQueryResult<PaginatedUser, Error> => {
  return useQuery({
    queryKey: [...listUserQueryKey(clientCode), opts?.searchParams],
    queryFn: () =>
      userService.list(clientCode, { ...opts, auth: "api-key" }),
    enabled: !!clientCode,
  })
}

export const meUserQueryKey = (clientCode: string) => [BASE_QUERY_KEY, clientCode, "me"] as const;
export const useMeUser = (
  clientCode: string,
  opts?: ServiceGetOptions<User>,
): UseQueryResult<User, Error> => {
  return useQuery({
    queryKey: meUserQueryKey(clientCode),
    queryFn: () => userService.me(clientCode, { ...opts, auth: "jwt" }),
    enabled: !!clientCode,
  })
}
