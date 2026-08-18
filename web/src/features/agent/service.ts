import { http } from "@/lib/http";

import {
  agentSchema, createAgentSchema, paginatedAgentSchema, updateAgentSchema
} from "./schemas";

import type { Agent, CreateAgent, PaginatedAgent, UpdateAgent } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const agentUrl = (tenantId: string) => `/api/tenants/${tenantId}/agents`;

const list = async (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedAgent>,
): Promise<PaginatedAgent> => {
  return http.get({
    url: agentUrl(tenantId),
    schemas: {
      response: paginatedAgentSchema,
    },
    ...opts,
  });
};

const findById = async (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<Agent>,
): Promise<Agent> => {
  return http.get({
    url: `${agentUrl(tenantId)}/{id}`,
    params: { id },
    schemas: { response: agentSchema },
    ...opts,
  });
};

const create = async (
  tenantId: string,
  body: CreateAgent,
  opts?: ServicePostOptions<CreateAgent, Agent>,
): Promise<Agent> => {
  return http.post({
    url: agentUrl(tenantId),
    body,
    schemas: {
      body: createAgentSchema,
      response: agentSchema,
    },
    ...opts,
  });
};

const update = async (
  tenantId: string,
  id: string,
  body: UpdateAgent,
  opts?: ServicePutOptions<UpdateAgent, void>,
): Promise<void> => {
  return http.put({
    url: `${agentUrl(tenantId)}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateAgentSchema,
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
    url: `${agentUrl(tenantId)}/{id}`,
    params: { id },
    ...opts,
  });
};

export const agentService = {
  list,
  findById,
  create,
  update,
  delete: remove,
};
