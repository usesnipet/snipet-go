import { http } from "@/lib/http";

import {
  agentSchema, createAgentSchema, paginatedAgentSchema, updateAgentSchema
} from "./schemas";

import type { Agent, CreateAgent, PaginatedAgent, UpdateAgent } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const AGENT_URL = "/api/agent";

const list = async (opts?: ServiceGetOptions<PaginatedAgent>): Promise<PaginatedAgent> => {
  return http.get({
    url: AGENT_URL,
    schemas: {
      response: paginatedAgentSchema,
    },
    ...opts,
  });
};

const findById = async (id: string, opts?: ServiceGetOptions<Agent>): Promise<Agent> => {
  return http.get({
    url: `${AGENT_URL}/{id}`,
    params: { id },
    schemas: { response: agentSchema },
    ...opts,
  });
};

const create = async (
  body: CreateAgent,
  opts?: ServicePostOptions<CreateAgent, Agent>,
): Promise<Agent> => {
  return http.post({
    url: AGENT_URL,
    body,
    schemas: {
      body: createAgentSchema,
      response: agentSchema,
    },
    ...opts,
  });
};

const update = async (
  id: string,
  body: UpdateAgent,
  opts?: ServicePutOptions<UpdateAgent, void>,
): Promise<void> => {
  return http.put({
    url: `${AGENT_URL}/{id}`,
    params: { id },
    body,
    schemas: {
      body: updateAgentSchema,
    },
    ...opts,
  });
};

const remove = async (id: string, opts?: ServiceDeleteOptions<void>): Promise<void> => {
  return http.delete({
    url: `${AGENT_URL}/{id}`,
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
