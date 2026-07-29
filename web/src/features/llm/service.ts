import { http } from "@/lib/http";

import {
  createLlmSchema,
  listDriversSchema,
  llmSchema,
  paginatedLlmSchema,
  updateLlmSchema,
} from "./schemas";

import type { CreateLlm, ListDrivers, Llm, PaginatedLlm, UpdateLlm } from "./schemas";
import type {
  ServiceDeleteOptions,
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";

const LLM_URL = "/api/llm";

const list = async (opts?: ServiceGetOptions<PaginatedLlm>): Promise<PaginatedLlm> => {
  return http.get({
    url: LLM_URL,
    schemas: {
      response: paginatedLlmSchema,
    },
    ...opts,
  });
};

const create = async (
  body: CreateLlm,
  opts?: ServicePostOptions<CreateLlm, Llm>,
): Promise<Llm> => {
  return http.post({
    url: LLM_URL,
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
    url: `${LLM_URL}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateLlmSchema,
    },
    ...opts,
  });
};

const remove = async (id: string, opts?: ServiceDeleteOptions<void>): Promise<void> => {
  return http.delete({
    url: `${LLM_URL}/{id}`,
    params: { id },
    ...opts,
  });
};

const listDrivers = async (opts?: ServiceGetOptions<ListDrivers>): Promise<ListDrivers> => {
  return http.get({
    url: `${LLM_URL}/drivers`,
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
