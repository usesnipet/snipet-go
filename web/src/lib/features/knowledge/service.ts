import { authenticatedClient } from "$lib/http/client";
import { queryClient } from "$lib/query-client";
import { createMutation, createQuery } from "@tanstack/svelte-query";
import { toast } from "svelte-sonner";

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

const listQueryKey = (filter?: FilterKnowledge) => [BASE_URL, "list", filter];
const listItemsQueryKey = (knowledgeId: string, filter?: FilterKnowledgeItem) => [BASE_URL, knowledgeId, "items", filter];
const findByIdQueryKey = (id: string) => [BASE_URL, id];

export const knowledgeService = {
  queryKeys: {
    listQueryKey,
    listItemsQueryKey,
    findByIdQueryKey,
  },
  list: (filter?: FilterKnowledge) => createQuery(() => ({
    queryKey: listQueryKey(filter),
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
    queryKey: listItemsQueryKey(knowledgeId, filter),
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
    queryKey: findByIdQueryKey(id),
    queryFn: () => authenticatedClient().get<Knowledge>({
      url: `${BASE_URL}/${id}`,
      schemas: { response: knowledgeSchema },
    }),
    refetchInterval: (query) => {
      const knowledge = query.state.data;
      return knowledge?.sync_status === "in_progress" ? 2000 : false;
    },
  })),
  create: (data: CreateKnowledge) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post<Knowledge>({
      url: BASE_URL,
      body: data,
      schemas: { body: createKnowledgeSchema, response: knowledgeSchema },
    }),
    onSuccess: () => {
      toast.success("Knowledge created.");
      queryClient.invalidateQueries({ queryKey: listQueryKey() });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  update: (id: string, data: UpdateKnowledge) => createMutation(() => ({
    mutationFn: () => authenticatedClient().put({
      url: `${BASE_URL}/${id}`,
      body: data,
      schemas: { body: updateKnowledgeSchema },
    }),
    onSuccess: () => {
      toast.success("Knowledge updated.");
      queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  delete: (id: string) => createMutation(() => ({
    mutationFn: () => authenticatedClient().delete({
      url: `${BASE_URL}/${id}`,
    }),
    onSuccess: () => {
      toast.success("Knowledge deleted.");
      queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listItemsQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listQueryKey() });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  testConnection: (data: TestConnection) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post({
      url: `${BASE_URL}/test-connection`,
      body: data,
      schemas: { body: testConnectionSchema },
    }),
    onSuccess: () => {
      toast.success("Connection test successful.");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  sync: (id: string, query?: SyncKnowledgeQuery) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post({
      url: `${BASE_URL}/${id}/sync`,
      searchParams: query,
      schemas: { searchParams: syncKnowledgeQuerySchema },
    }),
    onSuccess: () => {
      toast.success("Knowledge sync started.");
      queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listItemsQueryKey(id) });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
};
