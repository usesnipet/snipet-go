import { http } from "@/lib/http";

import { listUserSearchParamsSchema, paginatedUserSchema, userSchema } from "./schemas";

import type { ListUserSearchParams, PaginatedUser, User } from "./schemas";
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

const me = async (
  clientCode: string,
  opts?: ServiceGetOptions<User>,
): Promise<User> => {
  return http.get({
    url: `${userUrl(clientCode)}/me`,
    schemas: { response: userSchema },
    ...opts,
  })
}

export const userService = {
  list,
  me,
}
