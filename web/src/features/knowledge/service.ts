import { http } from "@/lib/http";

import {
  createKnowledgeIndexSchema,
  createKnowledgeResponseSchema,
  createKnowledgeSchema,
  knowledgeIndexSchema,
  knowledgeSchema,
  listKnowledgeDriversSchema,
  listKnowledgeIndexDriversSchema,
  paginatedKnowledgeIndexSchema,
  paginatedKnowledgeItemSchema,
  paginatedKnowledgeSchema,
  updateKnowledgeIndexSchema,
  updateKnowledgeSchema,
} from "./schemas";

import type {
  CreateKnowledge,
  CreateKnowledgeIndex,
  CreateKnowledgeResponse,
  Knowledge,
  KnowledgeIndex,
  ListKnowledgeDrivers,
  ListKnowledgeIndexDrivers,
  PaginatedKnowledge,
  PaginatedKnowledgeIndex,
  PaginatedKnowledgeItem,
  UpdateKnowledge,
  UpdateKnowledgeIndex,
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

const listItems = async (
  id: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeItem>,
): Promise<PaginatedKnowledgeItem> => {
  return http.get({
    url: `${KNOWLEDGE_URL}/{id}/items`,
    params: { id },
    schemas: {
      response: paginatedKnowledgeItemSchema,
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
  opts?: ServicePostOptions<undefined, void>,
): Promise<void> => {
  return http.post({
    url: `${KNOWLEDGE_URL}/{id}/sync`,
    params: { id },
    searchParams: { force },
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
  listItems,
  create,
  update,
  listDrivers,
  sync,
  delete: remove,
};

const listIndexes = async (
  knowledgeID: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeIndex>,
): Promise<PaginatedKnowledgeIndex> => {
  return http.get({
    url: `${KNOWLEDGE_URL}/{knowledge_id}/index`,
    params: { knowledge_id: knowledgeID },
    schemas: {
      response: paginatedKnowledgeIndexSchema,
    },
    ...opts,
  });
};

const findIndexByID = async (
  knowledgeID: string,
  id: string,
  opts?: ServiceGetOptions<KnowledgeIndex>,
): Promise<KnowledgeIndex> => {
  return http.get({
    url: `${KNOWLEDGE_URL}/{knowledge_id}/index/{id}`,
    params: { knowledge_id: knowledgeID, id },
    schemas: {
      response: knowledgeIndexSchema,
    },
    ...opts,
  });
};

const createIndex = async (
  knowledgeID: string,
  body: CreateKnowledgeIndex,
  opts?: ServicePostOptions<CreateKnowledgeIndex, KnowledgeIndex>,
): Promise<KnowledgeIndex> => {
  return http.post({
    url: `${KNOWLEDGE_URL}/{knowledge_id}/index`,
    params: { knowledge_id: knowledgeID },
    body,
    schemas: {
      body: createKnowledgeIndexSchema,
      response: knowledgeIndexSchema,
    },
    ...opts,
  });
};

const updateIndex = async (
  knowledgeID: string,
  id: string,
  body: UpdateKnowledgeIndex,
  opts?: ServicePutOptions<UpdateKnowledgeIndex, void>,
): Promise<void> => {
  return http.put({
    url: `${KNOWLEDGE_URL}/{knowledge_id}/index/{id}`,
    params: { knowledge_id: knowledgeID, id },
    body,
    schemas: {
      body: updateKnowledgeIndexSchema,
    },
    ...opts,
  });
};

const deleteIndex = async (
  knowledgeID: string,
  id: string,
  opts?: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${KNOWLEDGE_URL}/{knowledge_id}/index/{id}`,
    params: { knowledge_id: knowledgeID, id },
    ...opts,
  });
};

const listIndexDrivers = async (
  opts?: ServiceGetOptions<ListKnowledgeIndexDrivers>,
): Promise<DriverInfo[]> => {
  const response = await http.get({
    url: `${KNOWLEDGE_URL}/index/drivers`,
    schemas: { response: listKnowledgeIndexDriversSchema },
    ...opts,
  });
  return response.index_drivers;
};

export const knowledgeIndexService = {
  list: listIndexes,
  findByID: findIndexByID,
  create: createIndex,
  update: updateIndex,
  delete: deleteIndex,
  listDrivers: listIndexDrivers,
};
