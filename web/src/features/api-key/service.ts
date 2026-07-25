import { http } from "@/lib/http";

import {
  apiKeySchema, apiKeyWithSecretSchema, createApiKeySchema, paginatedApiKeySchema,
  updateApiKeyExpirationSchema
} from "./schemas";

import type {
  ApiKey, ApiKeyWithSecret, CreateApiKey, PaginatedApiKey, UpdateApiKeyExpiration
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const API_KEYS_URL = "/api/api-key";

const list = async (opts: ServiceGetOptions<PaginatedApiKey>): Promise<PaginatedApiKey> => {
  return http.get({
    url: API_KEYS_URL,
    schemas: {
      response: paginatedApiKeySchema,
    },
    ...opts,
  })
}

const create = async (
  body: CreateApiKey,
  opts: ServicePostOptions<CreateApiKey, ApiKeyWithSecret>,
): Promise<ApiKeyWithSecret> => {
  return http.post({
    url: API_KEYS_URL,
    body,
    schemas: {
      body: createApiKeySchema,
      response: apiKeyWithSecretSchema,
    },
    ...opts,
  })
}

const updateExpiration = async (
  id: string,
  body: UpdateApiKeyExpiration,
  opts: ServicePutOptions<UpdateApiKeyExpiration, void>,
): Promise<void> => {
  return http.put({
    url: `${API_KEYS_URL}/{id}/expiration`,
    params: { id },
    body,
    schemas: {
      body: updateApiKeyExpirationSchema,
    },
    ...opts,
  })
}

const me = async (opts: ServiceGetOptions<ApiKey>): Promise<ApiKey> => {
  return http.get({
    url: `${API_KEYS_URL}/me`,
    schemas: { response: apiKeySchema },
    ...opts,
  })
}

const findById = async (id: string, opts: ServiceGetOptions<ApiKey>): Promise<ApiKey> => {
  return http.get({
    url: `${API_KEYS_URL}/{id}`,
    params: { id },
    schemas: { response: apiKeySchema },
    ...opts,
  })
}

const roll = async (
  id: string,
  opts: ServicePostOptions<undefined, ApiKeyWithSecret>,
): Promise<ApiKeyWithSecret> => {
  return http.post({
    url: `${API_KEYS_URL}/{id}/roll`,
    params: { id },
    schemas: {
      response: apiKeyWithSecretSchema,
    },
    ...opts,
  })
}

const remove = async (id: string, opts: ServiceDeleteOptions<void>): Promise<void> => {
  return http.delete({
    url: `${API_KEYS_URL}/{id}`,
    params: { id },
    ...opts,
  })
}

export const apiKeyService = {
  list,
  create,
  me,
  findById,
  roll,
  updateExpiration,
  delete: remove,
}
