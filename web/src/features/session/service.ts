import { http, httpSse } from "@/lib/http";

import {
  createSessionSchema, findSessionSearchParamsSchema, listMessagesSearchParamsSchema,
  listSessionSearchParamsSchema, paginatedExecutionMessageSchema, paginatedSessionSchema, runSessionSchema,
  sessionSchema, updateSessionSchema
} from "./schemas";

import type {
  CreateSession, FindSessionSearchParams, ListMessagesSearchParams, ListSessionSearchParams,
  PaginatedExecutionMessage, PaginatedSession, RunSession, Session, UpdateSession
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { SseEventHandler } from "@/lib/http";

const sessionUrl = (appCode: string) => `/api/apps/${appCode}/session`;

const list = async (
  appCode: string,
  opts: ServiceGetOptions<PaginatedSession, ListSessionSearchParams>,
): Promise<PaginatedSession> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: sessionUrl(appCode),
    searchParams,
    schemas: {
      response: paginatedSessionSchema,
      searchParams: listSessionSearchParamsSchema,
    },
    ...rest,
  })
}

const create = async (
  appCode: string,
  body: CreateSession,
  opts: ServicePostOptions<CreateSession, Session>,
): Promise<Session> => {
  return http.post({
    url: sessionUrl(appCode),
    body,
    schemas: {
      body: createSessionSchema,
      response: sessionSchema,
    },
    ...opts,
  })
}

const findById = async (
  appCode: string,
  id: string,
  opts: ServiceGetOptions<Session, FindSessionSearchParams>,
): Promise<Session> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: `${sessionUrl(appCode)}/{id}`,
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
  appCode: string,
  id: string,
  body: UpdateSession,
  opts: ServicePutOptions<UpdateSession, void>,
): Promise<void> => {
  return http.put({
    url: `${sessionUrl(appCode)}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateSessionSchema,
    },
    ...opts,
  })
}

const remove = async (
  appCode: string,
  id: string,
  opts: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${sessionUrl(appCode)}/{id}`,
    params: { id },
    ...opts,
  })
}

const findMessages = async (
  appCode: string,
  id: string,
  opts: ServiceGetOptions<PaginatedExecutionMessage, ListMessagesSearchParams>,
): Promise<PaginatedExecutionMessage> => {
  const { searchParams, ...rest } = opts;
  return http.get({
    url: `${sessionUrl(appCode)}/{id}/messages`,
    params: { id },
    searchParams,
    schemas: {
      response: paginatedExecutionMessageSchema,
      searchParams: listMessagesSearchParamsSchema,
    },
    ...rest,
  })
}

export type RunSessionOptions = {
  signal?: AbortSignal;
  onEvent: SseEventHandler;
};

const run = async (
  appCode: string,
  id: string,
  body: RunSession,
  opts: RunSessionOptions,
): Promise<void> => {
  const { onEvent, signal } = opts;
  return httpSse({
    url: `${sessionUrl(appCode)}/{id}/run`,
    params: { id },
    method: "POST",
    body,
    signal,
    schemas: {
      body: runSessionSchema,
    },
    onEvent,
  });
};

export const sessionService = {
  list,
  create,
  findById,
  update,
  delete: remove,
  findMessages,
  run,
}
