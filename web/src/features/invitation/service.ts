import { http } from "@/lib/http";
import { memberSchema } from "@/models/member";

import {
  acceptInvitationSchema, createInvitationSchema, invitationSchema,
  listInvitationSearchParamsSchema, paginatedInvitationSchema
} from "./schemas";

import type {
  AcceptInvitation, CreateInvitation, Invitation, ListInvitationSearchParams, PaginatedInvitation
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions
} from "@/lib/services";
import type { Member } from "@/models/member";

const invitationsUrl = (tenantId: string) => `/api/tenants/${tenantId}/invitations`;
const ACCEPT_INVITATION_URL = "/api/invitations/accept";

const filter = async (
  tenantId: string,
  opts: ServiceGetOptions<PaginatedInvitation, ListInvitationSearchParams> = {},
): Promise<PaginatedInvitation> => {
  return http.get({
    url: invitationsUrl(tenantId),
    schemas: {
      response: paginatedInvitationSchema,
      searchParams: listInvitationSearchParamsSchema,
    },
    ...opts,
  })
}

const create = async (
  tenantId: string,
  body: CreateInvitation,
  opts: ServicePostOptions<CreateInvitation, Invitation> = {},
): Promise<Invitation> => {
  return http.post({
    url: invitationsUrl(tenantId),
    body,
    schemas: {
      body: createInvitationSchema,
      response: invitationSchema,
    },
    ...opts,
  })
}

const remove = async (
  tenantId: string,
  id: string,
  opts: ServiceDeleteOptions<void> = {},
): Promise<void> => {
  return http.delete({
    url: `${invitationsUrl(tenantId)}/{id}`,
    params: { id },
    ...opts,
  })
}

const accept = async (
  body: AcceptInvitation,
  opts: ServicePostOptions<AcceptInvitation, Member> = {},
): Promise<Member> => {
  return http.post({
    url: ACCEPT_INVITATION_URL,
    body,
    schemas: {
      body: acceptInvitationSchema,
      response: memberSchema,
    },
    ...opts,
  })
}

export const invitationService = {
  filter,
  create,
  remove,
  accept,
}
