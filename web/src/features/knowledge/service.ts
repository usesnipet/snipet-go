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

const knowledgeUrl = (tenantId: string) => `/api/tenants/${tenantId}/knowledge`;

const list = async (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedKnowledge>,
): Promise<PaginatedKnowledge> => {
  return http.get({
    url: knowledgeUrl(tenantId),
    schemas: {
      response: paginatedKnowledgeSchema,
    },
    ...opts,
  });
};

const findByID = async (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<Knowledge>,
): Promise<Knowledge> => {
  return http.get({
    url: `${knowledgeUrl(tenantId)}/{id}`,
    params: { id },
    schemas: {
      response: knowledgeSchema,
    },
    ...opts,
  });
};

const listItems = async (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeItem>,
): Promise<PaginatedKnowledgeItem> => {
  return http.get({
    url: `${knowledgeUrl(tenantId)}/{id}/items`,
    params: { id },
    schemas: {
      response: paginatedKnowledgeItemSchema,
    },
    ...opts,
  });
};

const create = async (
  tenantId: string,
  body: CreateKnowledge,
  opts?: ServicePostOptions<CreateKnowledge, CreateKnowledgeResponse>,
): Promise<CreateKnowledgeResponse> => {
  return http.post({
    url: knowledgeUrl(tenantId),
    body,
    schemas: {
      body: createKnowledgeSchema,
      response: createKnowledgeResponseSchema,
    },
    ...opts,
  });
};

const update = async (
  tenantId: string,
  id: string,
  body: UpdateKnowledge,
  opts?: ServicePutOptions<UpdateKnowledge, void>,
): Promise<void> => {
  return http.put({
    url: `${knowledgeUrl(tenantId)}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateKnowledgeSchema,
    },
    ...opts,
  });
};

const listDrivers = async (
  tenantId: string,
  opts?: ServiceGetOptions<ListKnowledgeDrivers>,
): Promise<DriverInfo[]> => {
  const response = await http.get({
    url: `${knowledgeUrl(tenantId)}/drivers`,
    schemas: { response: listKnowledgeDriversSchema },
    ...opts,
  });
  return response.source_drivers;
};

const sync = async (
  tenantId: string,
  id: string,
  force = false,
  opts?: ServicePostOptions<undefined, void>,
): Promise<void> => {
  return http.post({
    url: `${knowledgeUrl(tenantId)}/{id}/sync`,
    params: { id },
    searchParams: { force },
    ...opts,
  });
};

const remove = async (
  tenantId: string,
  id: string,
  opts?: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${knowledgeUrl(tenantId)}/{id}`,
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

const knowledgeIndexUrl = (tenantId: string, knowledgeId: string) =>
  `/api/tenants/${tenantId}/knowledge/${knowledgeId}/index`;

const listIndexes = async (
  tenantId: string,
  knowledgeID: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeIndex>,
): Promise<PaginatedKnowledgeIndex> => {
  return http.get({
    url: knowledgeIndexUrl(tenantId, knowledgeID),
    schemas: {
      response: paginatedKnowledgeIndexSchema,
    },
    ...opts,
  });
};

const findIndexByID = async (
  tenantId: string,
  knowledgeID: string,
  id: string,
  opts?: ServiceGetOptions<KnowledgeIndex>,
): Promise<KnowledgeIndex> => {
  return http.get({
    url: `${knowledgeIndexUrl(tenantId, knowledgeID)}/{id}`,
    params: { id },
    schemas: {
      response: knowledgeIndexSchema,
    },
    ...opts,
  });
};

const createIndex = async (
  tenantId: string,
  knowledgeID: string,
  body: CreateKnowledgeIndex,
  opts?: ServicePostOptions<CreateKnowledgeIndex, KnowledgeIndex>,
): Promise<KnowledgeIndex> => {
  return http.post({
    url: knowledgeIndexUrl(tenantId, knowledgeID),
    body,
    schemas: {
      body: createKnowledgeIndexSchema,
      response: knowledgeIndexSchema,
    },
    ...opts,
  });
};

const updateIndex = async (
  tenantId: string,
  knowledgeID: string,
  id: string,
  body: UpdateKnowledgeIndex,
  opts?: ServicePutOptions<UpdateKnowledgeIndex, void>,
): Promise<void> => {
  return http.put({
    url: `${knowledgeIndexUrl(tenantId, knowledgeID)}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateKnowledgeIndexSchema,
    },
    ...opts,
  });
};

const deleteIndex = async (
  tenantId: string,
  knowledgeID: string,
  id: string,
  opts?: ServiceDeleteOptions<void>,
): Promise<void> => {
  return http.delete({
    url: `${knowledgeIndexUrl(tenantId, knowledgeID)}/{id}`,
    params: { id },
    ...opts,
  });
};

const listIndexDrivers = async (
  tenantId: string,
  opts?: ServiceGetOptions<ListKnowledgeIndexDrivers>,
): Promise<DriverInfo[]> => {
  const response = await http.get({
    url: `${knowledgeUrl(tenantId)}/index/drivers`,
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
