import { http } from "@/lib/http";

import {
  createKnowledgeResponseSchema, createKnowledgeSchema, knowledgeSchema, listKnowledgeDriversSchema,
  paginatedKnowledgeSchema, syncKnowledgeResponseSchema, updateKnowledgeSchema
} from "./schemas";

import type {
  CreateKnowledge,
  CreateKnowledgeResponse,
  Knowledge,
  ListKnowledgeDrivers,
  PaginatedKnowledge,
  SyncKnowledgeResponse,
  UpdateKnowledge,
} from "./schemas";
import type {
  ServiceDeleteOptions,
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";
import type { DriverInfo } from "@/schemas/driver";

const KNOWLEDGE_URL = "/api/knowledge";

const list = async (
  opts?: ServiceGetOptions<PaginatedKnowledge>,
): Promise<PaginatedKnowledge> => {
  return http.get({
    url: KNOWLEDGE_URL,
    schemas: {
      response: paginatedKnowledgeSchema,
    },
    ...opts,
  });
};

const findByID = async (
  id: string,
  opts?: ServiceGetOptions<Knowledge>,
): Promise<Knowledge> => {
  return http.get({
    url: `${KNOWLEDGE_URL}/{id}`,
    params: { id },
    schemas: {
      response: knowledgeSchema,
    },
    ...opts,
  });
};

const create = async (
  body: CreateKnowledge,
  opts?: ServicePostOptions<CreateKnowledge, CreateKnowledgeResponse>,
): Promise<CreateKnowledgeResponse> => {
  return http.post({
    url: KNOWLEDGE_URL,
    body,
    schemas: {
      body: createKnowledgeSchema,
      response: createKnowledgeResponseSchema,
    },
    ...opts,
  });
};

const update = async (
  id: string,
  body: UpdateKnowledge,
  opts?: ServicePutOptions<UpdateKnowledge, void>,
): Promise<void> => {
  return http.put({
    url: `${KNOWLEDGE_URL}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateKnowledgeSchema,
    },
    ...opts,
  });
};

const listDrivers = async (
  opts?: ServiceGetOptions<ListKnowledgeDrivers>,
): Promise<DriverInfo[]> => {
  const response = await http.get({
    url: `${KNOWLEDGE_URL}/drivers`,
    schemas: { response: listKnowledgeDriversSchema },
    ...opts,
  });
  return response.source_drivers;
};

const sync = async (
  id: string,
  force = false,
  opts?: ServicePostOptions<undefined, SyncKnowledgeResponse>,
): Promise<SyncKnowledgeResponse> => {
  return http.post({
    url: `${KNOWLEDGE_URL}/{id}/sync`,
    params: { id },
    searchParams: { force },
    schemas: {
      response: syncKnowledgeResponseSchema,
    },
    ...opts,
  });
};

const remove = async (id: string, opts?: ServiceDeleteOptions<void>): Promise<void> => {
  return http.delete({
    url: `${KNOWLEDGE_URL}/{id}`,
    params: { id },
    ...opts,
  });
};

export const knowledgeService = {
  list,
  findByID,
  create,
  update,
  listDrivers,
  sync,
  delete: remove,
};
