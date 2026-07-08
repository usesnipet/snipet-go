import { authenticatedClient, publicClient } from "$lib/http/client";
import { createMutation, createQuery } from "@tanstack/svelte-query";

import {
  apiKeyPaginatedSchema, apiKeySchema, apiKeyWithSecretSchema, createApiKeySchema,
  updateApiKeyExpirationSchema
} from "./schemas";

import type { APIKey, APIKeyPaginated, CreateApiKey, UpdateApiKeyExpiration } from "./schemas";

const BASE_URL = "/api/api-key";

export const apiKeyService = {
  check: () => createMutation(() => ({
    mutationFn: async (apiKey: string) => {
      const res = await publicClient().get<APIKey>({
        url: `${BASE_URL}/me`, headers: { "X-API-Key": apiKey },
        schemas: { response: apiKeySchema }
      });
      return !!res;
    },
  })),
  me: () => createQuery(() => ({
    queryKey: ["api-keys", "me"],
    queryFn: () => {
      return authenticatedClient().get<APIKey>({
        url: `${BASE_URL}/me`,
        schemas: { response: apiKeySchema }
      });
    }
  })),
  list: () => createQuery(() => ({
    queryKey: ["api-keys", "list"],
    queryFn: async () => {
      const res = await authenticatedClient().get<APIKeyPaginated>({
        url: `${BASE_URL}`,
        schemas: { response: apiKeyPaginatedSchema }
      });
      return res.data;
    }
  })),
  create: (data: CreateApiKey) => createMutation(() => ({
    mutationFn: () => authenticatedClient().post({
      url: `${BASE_URL}`,
      body: data,
      schemas: { body: createApiKeySchema, response: apiKeyWithSecretSchema }
    })
  })),
  updateExpiration: (id: string, data: UpdateApiKeyExpiration) => createMutation(() => ({
    mutationFn: () => authenticatedClient().put({
      url: `${BASE_URL}/${id}/expiration`,
      body: data,
      schemas: { body: updateApiKeyExpirationSchema, response: apiKeySchema }
    })
  })),
};