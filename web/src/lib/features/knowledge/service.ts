import { authenticatedClient } from "$lib/http/client";
import { createMutation, createQuery } from "@tanstack/svelte-query";

import {
  createKnowledgeSchema, filterKnowledgeItemSchema, filterKnowledgeSchema, knowledgeItemPaginatedSchema,
  knowledgePaginatedSchema, knowledgeSchema, syncKnowledgeQuerySchema, testConnectionSchema,
  updateKnowledgeSchema
} from "./schemas";

import type {
  CreateKnowledge,
  FilterKnowledge,
  FilterKnowledgeItem,
  Knowledge,
  KnowledgePaginated,
  KnowledgeItemPaginated,
  SyncKnowledgeQuery,
  TestConnection,
  UpdateKnowledge,
} from "./schemas";

const BASE_URL = "/api/knowledge";

export const knowledgeService = {
  list: (filter?: FilterKnowledge) => createQuery(() => ({
    queryKey: [BASE_URL, "list", filter],
    queryFn: async () => {
      const res = await authenticatedClient().get<KnowledgePaginated>({
        url: BASE_URL,
        searchParams: filter,
        schemas: { response: knowledgePaginatedSchema, searchParams: filterKnowledgeSchema },
      });
      return res.data;
    },
  })),
  listItems: (knowledgeId: string, filter?: FilterKnowledgeItem) => createQuery(() => ({
    queryKey: [BASE_URL, knowledgeId, "items", filter],
    queryFn: async () => {
      const res = await authenticatedClient().get<KnowledgeItemPaginated>({
        url: `${BASE_URL}/${knowledgeId}/items`,
        searchParams: filter,
        schemas: { response: knowledgeItemPaginatedSchema, searchParams: filterKnowledgeItemSchema },
      });
      return res.data;
    },
  })),
  findById: (id: string) => createQuery(() => ({
    queryKey: [BASE_URL, id],
    queryFn: () => authenticatedClient().get<Knowledge>({
      url: `${BASE_URL}/${id}`,
      schemas: { response: knowledgeSchema },
    }),
  })),
  create: (data: CreateKnowledge) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post({
      url: BASE_URL,
      body: data,
      schemas: { body: createKnowledgeSchema, response: knowledgeSchema },
    }),
  })),
  update: (id: string, data: UpdateKnowledge) => createMutation(() => ({
    mutationFn: () => authenticatedClient().put({
      url: `${BASE_URL}/${id}`,
      body: data,
      schemas: { body: updateKnowledgeSchema },
    }),
  })),
  delete: (id: string) => createMutation(() => ({
    mutationFn: () => authenticatedClient().delete({
      url: `${BASE_URL}/${id}`,
    }),
  })),
  testConnection: (data: TestConnection) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post({
      url: `${BASE_URL}/test-connection`,
      body: data,
      schemas: { body: testConnectionSchema },
    }),
  })),
  sync: (id: string, query?: SyncKnowledgeQuery) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post({
      url: `${BASE_URL}/${id}/sync`,
      searchParams: query,
      schemas: { searchParams: syncKnowledgeQuerySchema },
    }),
  })),
};
