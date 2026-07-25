import { http } from "@/lib/http";

import {
  createSessionSchema, findSessionSearchParamsSchema, listMessagesSearchParamsSchema,
  listSessionSearchParamsSchema, paginatedExecutionMessageSchema, paginatedSessionSchema, sessionSchema,
  updateSessionSchema
} from "./schemas";

import type {
  CreateSession, FindSessionSearchParams, ListMessagesSearchParams, ListSessionSearchParams,
  PaginatedExecutionMessage, PaginatedSession, Session, UpdateSession
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const sessionUrl = (clientCode: string) => `/api/client/${clientCode}/session`;

const list = async (
  clientCode: string,
  opts: ServiceGetOptions<PaginatedSession> & { searchParams?: ListSessionSearchParams },
): Promise<PaginatedSession> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: sessionUrl(clientCode),
    searchParams,
    schemas: {
      response: paginatedSessionSchema,
      searchParams: listSessionSearchParamsSchema,
    },
    ...rest,
  })
}

const create = async (
  clientCode: string,
  body: CreateSession,
  opts: ServicePostOptions<CreateSession, Session>,
): Promise<Session> => {
  return http.post({
    url: sessionUrl(clientCode),
    body,
    schemas: {
      body: createSessionSchema,
      response: sessionSchema,
    },
    ...opts,
  })
}

const findById = async (
  clientCode: string,
  id: string,
  opts: ServiceGetOptions<Session> & { searchParams?: FindSessionSearchParams },
): Promise<Session> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: `${sessionUrl(clientCode)}/{id}`,
    params: { id },
    searchParams,
    schemas: {
      response: sessionSchema,
      searchParams: findSessionSearchParamsSchema,
    },
    ...rest,
  })
}

const update = async (
  clientCode: string,
  id: string,
  body: UpdateSession,
  opts: ServicePutOptions<UpdateSession, void>,
): Promise<void> => {
  return http.put({
    url: `${sessionUrl(clientCode)}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateSessionSchema,
    },
    ...opts,
  })
}

const remove = async (
  clientCode: string,
  id: string,
  opts: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${sessionUrl(clientCode)}/{id}`,
    params: { id },
    ...opts,
  })
}

const findMessages = async (
  clientCode: string,
  id: string,
  opts: ServiceGetOptions<PaginatedExecutionMessage> & {
    searchParams?: ListMessagesSearchParams
  },
): Promise<PaginatedExecutionMessage> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: `${sessionUrl(clientCode)}/{id}/messages`,
    params: { id },
    searchParams,
    schemas: {
      response: paginatedExecutionMessageSchema,
      searchParams: listMessagesSearchParamsSchema,
    },
    ...rest,
  })
}

export const sessionService = {
  list,
  create,
  findById,
  update,
  delete: remove,
  findMessages,
}
