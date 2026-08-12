import { http } from "@/lib/http";

import { listUserSearchParamsSchema, paginatedUserSchema } from "./schemas";

import type { ListUserSearchParams, PaginatedUser } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";

const userUrl = (clientCode: string) => `/api/client/${clientCode}/user`;

const list = async (
  clientCode: string,
  opts: ServiceGetOptions<PaginatedUser, ListUserSearchParams> = {},
): Promise<PaginatedUser> => {
  return http.get({
    url: userUrl(clientCode),
    schemas: {
      response: paginatedUserSchema,
      searchParams: listUserSearchParamsSchema,
    },
    ...opts
  })
}
export const clientUserService = {
  list,
}
