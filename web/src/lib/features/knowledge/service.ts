import { authenticatedClient } from "$lib/http/client";
import { queryClient } from "$lib/query-client";
import { createMutation, createQuery } from "@tanstack/svelte-query";
import { toast } from "svelte-sonner";

import {
  createKnowledgeSchema, driversSchema, filterKnowledgeItemSchema, filterKnowledgeSchema,
  knowledgeItemPaginatedSchema, knowledgePaginatedSchema, knowledgeSchema, syncKnowledgeQuerySchema,
  testConnectionSchema, updateKnowledgeSchema
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
  Drivers,
} from "./schemas";
const BASE_URL = "/api/knowledge";

const listQueryKey = (filter?: FilterKnowledge) => [BASE_URL, "list", filter];
const listItemsQueryKey = (knowledgeId: string, filter?: FilterKnowledgeItem) => [BASE_URL, knowledgeId, "items", filter];
const findByIdQueryKey = (id: string) => [BASE_URL, id];
const listDriversQueryKey = () => [BASE_URL, "drivers"];
export const knowledgeService = {
  queryKeys: {
    listQueryKey,
    listItemsQueryKey,
    findByIdQueryKey,
    listDriversQueryKey,
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
  create: () => createMutation(() => ({
    mutationFn: (data: CreateKnowledge) => authenticatedClient().post<Knowledge>({
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
  update: () => createMutation(() => ({
    mutationFn: ({ data, id }: { id: string, data: UpdateKnowledge }) => authenticatedClient().put({
      url: `${BASE_URL}/${id}`,
      body: data,
      schemas: { body: updateKnowledgeSchema },
    }),
    onSuccess: (_, { id }) => {
      toast.success("Knowledge updated.");
      queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listQueryKey() });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  delete: () => createMutation(() => ({
    mutationFn: (id: string) => authenticatedClient().delete({
      url: `${BASE_URL}/${id}`,
    }),
    onSuccess: (_, id) => {
      toast.success("Knowledge deleted.");
      queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listItemsQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listQueryKey() });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  testConnection: () => createMutation(() => ({
    mutationFn: (data: TestConnection) => authenticatedClient().post({
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
  sync: () => createMutation(() => ({
    mutationFn: ({ id, query }: { id: string, query?: SyncKnowledgeQuery }) => authenticatedClient().post({
      url: `${BASE_URL}/${id}/sync`,
      searchParams: query,
      schemas: { searchParams: syncKnowledgeQuerySchema },
    }),
    onSuccess: (_, { id }) => {
      toast.success("Knowledge sync started.");
      queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listItemsQueryKey(id) });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  })),
  listDrivers: () => createQuery(() => ({
    queryKey: listDriversQueryKey(),
    queryFn: () => authenticatedClient().get<Drivers>({
      url: `${BASE_URL}/drivers`,
      schemas: { response: driversSchema },
    }),
  })),
};
