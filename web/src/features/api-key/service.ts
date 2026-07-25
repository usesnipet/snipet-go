import { http } from "@/lib/http";

import {
  apiKeySchema, apiKeyWithSecretSchema, createApiKeySchema, paginatedApiKeySchema,
  updateApiKeyExpirationSchema
} from "./schemas";

import type { ApiKey, UpdateApiKeyExpiration, ApiKeyWithSecret, CreateApiKey } from "./schemas";
import type { Paginated } from "@/schemas/paginated";
import type { ServiceGetOptions, ServicePostOptions, ServicePutOptions } from "@/lib/services";

const API_KEYS_URL = "/api/api-key";

const list = async (opts?: ServiceGetOptions<Paginated<ApiKey>>) => {
  return http.get({
    url: API_KEYS_URL,
    schemas: {
      response: paginatedApiKeySchema,
    },
    ...opts,
  })
}

const create = async (body: CreateApiKey, opts?: ServicePostOptions<CreateApiKey, ApiKeyWithSecret>) => {
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
  opts?: ServicePutOptions<UpdateApiKeyExpiration, ApiKeyWithSecret>,
) => {
  return http.put({
    url: `${API_KEYS_URL}/{id}/expiration`,
    params: { id },
    body,
    schemas: {
      body: updateApiKeyExpirationSchema,
      response: apiKeyWithSecretSchema,
    },
    ...opts,
  })
}

const me = async (opts?: ServiceGetOptions<ApiKey>) => {
  return http.get({
    url: `${API_KEYS_URL}/me`,
    schemas: { response: apiKeySchema },
    ...opts,
  })
}

const findById = async (id: string, opts?: ServiceGetOptions<ApiKey>) => {
  return http.get({
    url: `${API_KEYS_URL}/{id}`,
    params: { id },
    schemas: { response: apiKeySchema },
    ...opts,
  })
}

const roll = async (id: string, opts?: ServicePostOptions<undefined, ApiKeyWithSecret>) => {
  return http.post({
    url: `${API_KEYS_URL}/{id}/roll`,
    params: { id },
    schemas: {
      response: apiKeyWithSecretSchema,
    },
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
}