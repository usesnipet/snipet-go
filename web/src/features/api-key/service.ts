import { http } from "@/lib/http";

import {
  apiKeySchema, apiKeyWithSecretSchema, createApiKeySchema,
  listApiKeySearchParamsSchema, paginatedApiKeySchema, updateApiKeyExpirationSchema
} from "./schemas";

import type {
  ApiKey, ApiKeyWithSecret, CreateApiKey, ListApiKeySearchParams, PaginatedApiKey,
  UpdateApiKeyExpiration
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const apiKeysUrl = (tenantId: string) => `/api/tenants/${tenantId}/api-keys`;
const API_KEY_ME_URL = "/api/api-key/me";

const list = async (
  tenantId: string,
  opts: ServiceGetOptions<PaginatedApiKey, ListApiKeySearchParams> = {},
): Promise<PaginatedApiKey> => {
  return http.get({
    url: apiKeysUrl(tenantId),
    schemas: {
      response: paginatedApiKeySchema,
      searchParams: listApiKeySearchParamsSchema,
    },
    ...opts,
  })
}

const create = async (
  tenantId: string,
  body: CreateApiKey,
  opts: ServicePostOptions<CreateApiKey, ApiKeyWithSecret> = {},
): Promise<ApiKeyWithSecret> => {
  return http.post({
    url: apiKeysUrl(tenantId),
    body,
    schemas: {
      body: createApiKeySchema,
      response: apiKeyWithSecretSchema,
    },
    ...opts,
  })
}

const updateExpiration = async (
  tenantId: string,
  id: string,
  body: UpdateApiKeyExpiration,
  opts: ServicePutOptions<UpdateApiKeyExpiration, void> = {},
): Promise<void> => {
  return http.put({
    url: `${apiKeysUrl(tenantId)}/{id}/expiration`,
    params: { id },
    body,
    schemas: {
      body: updateApiKeyExpirationSchema,
    },
    ...opts,
  })
}

const me = async (opts: ServiceGetOptions<ApiKey> = {}): Promise<ApiKey> => {
  return http.get({
    url: API_KEY_ME_URL,
    schemas: { response: apiKeySchema },
    ...opts,
  })
}

const findById = async (
  tenantId: string,
  id: string,
  opts: ServiceGetOptions<ApiKey> = {},
): Promise<ApiKey> => {
  return http.get({
    url: `${apiKeysUrl(tenantId)}/{id}`,
    params: { id },
    schemas: { response: apiKeySchema },
    ...opts,
  })
}

const roll = async (
  tenantId: string,
  id: string,
  opts: ServicePostOptions<undefined, ApiKeyWithSecret> = {},
): Promise<ApiKeyWithSecret> => {
  return http.post({
    url: `${apiKeysUrl(tenantId)}/{id}/roll`,
    params: { id },
    schemas: {
      response: apiKeyWithSecretSchema,
    },
    ...opts,
  })
}

const remove = async (
  tenantId: string,
  id: string,
  opts: ServiceDeleteOptions<void> = {},
): Promise<void> => {
  return http.delete({
    url: `${apiKeysUrl(tenantId)}/{id}`,
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
