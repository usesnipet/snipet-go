import { http } from "@/lib/http";

import {
  createLlmSchema, listDriversSchema, listLlmSearchParamsSchema, llmSchema, paginatedLlmSchema,
  updateLlmSchema
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

const llmUrl = () => `/api/llm`;

const list = async (
  opts?: ServiceGetOptions<PaginatedLlm, ListLlmSearchParams>,
): Promise<PaginatedLlm> => {
  return http.get({
    url: llmUrl(),
    schemas: {
      response: paginatedLlmSchema,
      searchParams: listLlmSearchParamsSchema,
    },
    ...opts,
  });
};

const create = async (
  body: CreateLlm,
  opts?: ServicePostOptions<CreateLlm, Llm>,
): Promise<Llm> => {
  return http.post({
    url: llmUrl(),
    body,
    schemas: {
      body: createLlmSchema,
      response: llmSchema,
    },
    ...opts,
  });
};

const update = async (
  id: string,
  body: UpdateLlm,
  opts?: ServicePutOptions<UpdateLlm, void>,
): Promise<void> => {
  return http.put({
    url: `${llmUrl()}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateLlmSchema,
    },
    ...opts,
  });
};

const remove = async (
  id: string,
  opts?: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${llmUrl()}/{id}`,
    params: { id },
    ...opts,
  });
};

const listDrivers = async (
  opts?: ServiceGetOptions<ListDrivers>,
): Promise<ListDrivers> => {
  return http.get({
    url: `${llmUrl()}/drivers`,
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
