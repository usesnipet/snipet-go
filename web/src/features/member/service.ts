import { http } from "@/lib/http";

import {
  createMemberSchema, listMemberSearchParamsSchema, memberSchema, paginatedMemberSchema,
  updateMemberRoleSchema
} from "./schemas";

import type { CreateMember, ListMemberSearchParams, Member, PaginatedMember, UpdateMemberRole } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const membersUrl = (tenantId: string) => `/api/tenants/${tenantId}/members`;

const filter = async (
  tenantId: string,
  opts: ServiceGetOptions<PaginatedMember, ListMemberSearchParams> = {},
): Promise<PaginatedMember> => {
  return http.get({
    url: membersUrl(tenantId),
    schemas: {
      response: paginatedMemberSchema,
      searchParams: listMemberSearchParamsSchema,
    },
    ...opts,
  })
}

const create = async (
  tenantId: string,
  body: CreateMember,
  opts: ServicePostOptions<CreateMember, Member> = {},
): Promise<Member> => {
  return http.post({
    url: membersUrl(tenantId),
    body,
    schemas: {
      body: createMemberSchema,
      response: memberSchema,
    },
    ...opts,
  })
}

const updateRole = async (
  tenantId: string,
  id: string,
  body: UpdateMemberRole,
  opts: ServicePutOptions<UpdateMemberRole, void> = {},
): Promise<void> => {
  return http.put({
    url: `${membersUrl(tenantId)}/{id}/role`,
    params: { id },
    body,
    schemas: { body: updateMemberRoleSchema },
    ...opts,
  })
}

const remove = async (
  tenantId: string,
  id: string,
  opts: ServiceDeleteOptions<void> = {},
): Promise<void> => {
  return http.delete({
    url: `${membersUrl(tenantId)}/{id}`,
    params: { id },
    ...opts,
  })
}

export const memberService = {
  filter,
  create,
  updateRole,
  remove,
}
