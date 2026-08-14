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
const CLIENTS_URL = "/api/clients";

const list = async (opts?: ServiceGetOptions<PaginatedClient>): Promise<PaginatedClient> => {
  return http.get({
    url: CLIENTS_URL,
    schemas: {
      response: paginatedClientSchema,
    },
    ...opts,
  })
}
const listAgents = async (code: string, opts?: ServiceGetOptions<PaginatedAgent>): Promise<PaginatedAgent> => {
  return http.get({
    url: `${CLIENTS_URL}/{code}/agents`,
    params: { code },
    schemas: {
      response: paginatedAgentSchema,
    },
    ...opts,
  })
}

const create = async (
  body: CreateClient,
  opts?: ServicePostOptions<CreateClient, Client>,
): Promise<Client> => {
  return http.post({
    url: CLIENTS_URL,
    body,
    schemas: {
      body: createClientSchema,
      response: clientSchema,
    },
    ...opts,
  })
}

const findByCode = async (code: string, opts?: ServiceGetOptions<Client>): Promise<Client> => {
  return http.get({
    url: `${CLIENTS_URL}/{code}`,
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
    url: `${CLIENTS_URL}/{code}/public`,
    params: { code },
    schemas: { response: clientPublicSchema },
    ...opts,
  })
}

const update = async (
  code: string,
  body: UpdateClient,
  opts: ServicePutOptions<UpdateClient, void>,
): Promise<void> => {
  return http.put({
    url: `${CLIENTS_URL}/{code}`,
    params: { code },
    body,
    schemas: {
      body: updateClientSchema,
    },
    ...opts,
  })
}

const remove = async (code: string, opts: ServiceDeleteOptions<void>): Promise<void> => {
  return http.delete({
    url: `${CLIENTS_URL}/{code}`,
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
