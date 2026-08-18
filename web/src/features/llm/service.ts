import { http } from "@/lib/http";

import {
  createLlmSchema,
  listDriversSchema,
  listLlmSearchParamsSchema,
  llmSchema,
  paginatedLlmSchema,
  updateLlmSchema,
} from "./schemas";

import type {
  CreateLlm, ListDrivers, ListLlmSearchParams, Llm, PaginatedLlm, UpdateLlm
} from "./schemas";
import type {
  ServiceDeleteOptions,
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";

const llmUrl = (tenantId: string) => `/api/tenants/${tenantId}/llm`;

const list = async (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedLlm, ListLlmSearchParams>,
): Promise<PaginatedLlm> => {
  return http.get({
    url: llmUrl(tenantId),
    schemas: {
      response: paginatedLlmSchema,
      searchParams: listLlmSearchParamsSchema,
    },
    ...opts,
  });
};

const create = async (
  tenantId: string,
  body: CreateLlm,
  opts?: ServicePostOptions<CreateLlm, Llm>,
): Promise<Llm> => {
  return http.post({
    url: llmUrl(tenantId),
    body,
    schemas: {
      body: createLlmSchema,
      response: llmSchema,
    },
    ...opts,
  });
};

const update = async (
  tenantId: string,
  id: string,
  body: UpdateLlm,
  opts?: ServicePutOptions<UpdateLlm, void>,
): Promise<void> => {
  return http.put({
    url: `${llmUrl(tenantId)}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateLlmSchema,
    },
    ...opts,
  });
};

const remove = async (
  tenantId: string,
  id: string,
  opts?: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${llmUrl(tenantId)}/{id}`,
    params: { id },
    ...opts,
  });
};

const listDrivers = async (
  tenantId: string,
  opts?: ServiceGetOptions<ListDrivers>,
): Promise<ListDrivers> => {
  return http.get({
    url: `${llmUrl(tenantId)}/drivers`,
    schemas: { response: listDriversSchema },
    ...opts,
  });
};

export const llmService = {
  list,
  create,
  update,
  delete: remove,
  listDrivers,
};
