import { http } from "@/lib/http";

import { listUserSearchParamsSchema, paginatedUserSchema } from "./schemas";

import type { ListUserSearchParams, PaginatedUser } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";

const userUrl = (clientCode: string) => `/api/client/${clientCode}/user`;

const list = async (
  clientCode: string,
  opts: ServiceGetOptions<PaginatedUser, ListUserSearchParams>,
): Promise<PaginatedUser> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: userUrl(clientCode),
    searchParams,
    schemas: {
      response: paginatedUserSchema,
      searchParams: listUserSearchParamsSchema,
    },
    ...rest,
  })
}

export const userService = {
  list,
}
