import { http } from "@/lib/http";

import { paginatedAgentSchema } from "../agent/schemas";

import {
  clientPublicSchema, clientSchema, createClientSchema, paginatedClientSchema, updateClientSchema
} from "./schemas";

import type { PaginatedAgent } from "../agent/schemas";
import type {
  Client, ClientPublic, CreateClient, PaginatedClient, UpdateClient
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const CLIENTS_PUBLIC_URL = "/api/clients";
const clientsUrl = (tenantId: string) => `/api/tenants/${tenantId}/clients`;

const list = async (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedClient>,
): Promise<PaginatedClient> => {
  return http.get({
    url: clientsUrl(tenantId),
    schemas: {
      response: paginatedClientSchema,
    },
    ...opts,
  })
}
const listAgents = async (code: string, opts?: ServiceGetOptions<PaginatedAgent>): Promise<PaginatedAgent> => {
  return http.get({
    url: `${CLIENTS_PUBLIC_URL}/{code}/agents`,
    params: { code },
    schemas: {
      response: paginatedAgentSchema,
    },
    ...opts,
  })
}

const create = async (
  tenantId: string,
  body: CreateClient,
  opts?: ServicePostOptions<CreateClient, Client>,
): Promise<Client> => {
  return http.post({
    url: clientsUrl(tenantId),
    body,
    schemas: {
      body: createClientSchema,
      response: clientSchema,
    },
    ...opts,
  })
}

const findByCode = async (
  tenantId: string,
  code: string,
  opts?: ServiceGetOptions<Client>,
): Promise<Client> => {
  return http.get({
    url: `${clientsUrl(tenantId)}/{code}`,
    params: { code },
    schemas: { response: clientSchema },
    ...opts,
  })
}

const findPublicByCode = async (
  code: string,
  opts?: ServiceGetOptions<ClientPublic>,
): Promise<ClientPublic> => {
  return http.get({
    url: `${CLIENTS_PUBLIC_URL}/{code}/public`,
    params: { code },
    schemas: { response: clientPublicSchema },
    ...opts,
  })
}

const update = async (
  tenantId: string,
  code: string,
  body: UpdateClient,
  opts: ServicePutOptions<UpdateClient, void>,
): Promise<void> => {
  return http.put({
    url: `${clientsUrl(tenantId)}/{code}`,
    params: { code },
    body,
    schemas: {
      body: updateClientSchema,
    },
    ...opts,
  })
}

const remove = async (
  tenantId: string,
  code: string,
  opts: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${clientsUrl(tenantId)}/{code}`,
    params: { code },
    ...opts,
  })
}

export const clientService = {
  list,
  listAgents,
  create,
  findByCode,
  findPublicByCode,
  update,
  delete: remove,
}
