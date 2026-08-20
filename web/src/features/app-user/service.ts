import { http } from "@/lib/http";

import { listUserSearchParamsSchema, paginatedUserSchema } from "./schemas";

import type { ListUserSearchParams, PaginatedUser } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";

const userUrl = (appCode: string) => `/api/apps/${appCode}/user`;

const list = async (
  appCode: string,
  opts: ServiceGetOptions<PaginatedUser, ListUserSearchParams> = {},
): Promise<PaginatedUser> => {
  return http.get({
    url: userUrl(appCode),
    schemas: {
      response: paginatedUserSchema,
      searchParams: listUserSearchParamsSchema,
    },
    ...opts
  })
}
export const appUserService = {
  list,
}
