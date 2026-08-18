import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { apiKeyService } from "./service";

import type {
  ApiKey, ApiKeyWithSecret, CreateApiKey, ListApiKeySearchParams, PaginatedApiKey,
  UpdateApiKeyExpiration
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "api-key";

export const listApiKeyQueryKey = (tenantId: string) => [BASE_QUERY_KEY, "list", tenantId] as const;
export const useListApiKey = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedApiKey, ListApiKeySearchParams>
): UseQueryResult<PaginatedApiKey, Error> => {
  return useQuery({
    queryKey: [...listApiKeyQueryKey(tenantId), opts?.searchParams],
    queryFn: () => apiKeyService.list(tenantId, opts),
    enabled: !!tenantId,
  })
}

export const createApiKeyQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateApiKey = (
  opts?: ServicePostOptions<CreateApiKey, ApiKeyWithSecret>
): UseMutationResult<ApiKeyWithSecret, Error, { tenantId: string; data: CreateApiKey }> => {
  return useMutation({
    mutationKey: createApiKeyQueryKey(),
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CreateApiKey }) =>
      apiKeyService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "API Key created successfully",
        description: "The API key has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listApiKeyQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to create API Key",
        description: "The API Key has not been created successfully",
        variant: "destructive",
      });
    }
  })
}

export const meApiKeyQueryKey = () => [BASE_QUERY_KEY, "me"];
export const useMeApiKey = (
  opts?: ServiceGetOptions<ApiKey>
): UseQueryResult<ApiKey, Error> => {
  return useQuery({
    queryKey: meApiKeyQueryKey(),
    queryFn: (): Promise<ApiKey> =>
      apiKeyService.me(opts),
  })
}

export const findByIdApiKeyQueryKey = (tenantId: string, id: string) =>
  [BASE_QUERY_KEY, "findById", tenantId, id];
export const useFindByIdApiKey = (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<ApiKey>
): UseQueryResult<ApiKey, Error> => {
  return useQuery({
    queryKey: findByIdApiKeyQueryKey(tenantId, id),
    queryFn: (): Promise<ApiKey> =>
      apiKeyService.findById(tenantId, id, opts),
    enabled: !!tenantId && !!id,
  })
}

export const rollApiKeyQueryKey = () => [BASE_QUERY_KEY, "roll"];
export const useRollApiKey = (
  opts?: ServicePostOptions<undefined, ApiKeyWithSecret>
): UseMutationResult<ApiKeyWithSecret, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: rollApiKeyQueryKey(),
    mutationFn: ({ tenantId, id }: { tenantId: string; id: string }) =>
      apiKeyService.roll(tenantId, id, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "API Key rolled successfully",
        description: "The API key has been rolled successfully",
      });
      queryClient.invalidateQueries({ queryKey: listApiKeyQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: meApiKeyQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdApiKeyQueryKey(tenantId, id) });
    },
    onError: () => {
      toast({
        title: "Failed to roll API key",
        description: "The API key has not been rolled successfully",
        variant: "destructive",
      });
    }
  })
}

export const updateExpirationApiKeyQueryKey = () => [BASE_QUERY_KEY, "updateExpiration"];
export const useUpdateExpirationApiKey = (
  opts?: ServicePutOptions<UpdateApiKeyExpiration, void>
): UseMutationResult<void, Error, { tenantId: string; id: string; data: UpdateApiKeyExpiration }> => {
  return useMutation({
    mutationKey: updateExpirationApiKeyQueryKey(),
    mutationFn: ({ tenantId, id, data }: { tenantId: string; id: string; data: UpdateApiKeyExpiration }) =>
      apiKeyService.updateExpiration(tenantId, id, data, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "API Key expiration updated successfully",
        description: "The API key expiration has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listApiKeyQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: meApiKeyQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdApiKeyQueryKey(tenantId, id) });
    },
    onError: () => {
      toast({
        title: "Failed to update API Key expiration",
        description: "The API Key expiration has not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const deleteApiKeyQueryKey = () => [BASE_QUERY_KEY, "delete"];
export const useDeleteApiKey = (
  opts?: ServiceDeleteOptions<void>
): UseMutationResult<void, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: deleteApiKeyQueryKey(),
    mutationFn: ({ tenantId, id }: { tenantId: string; id: string }) =>
      apiKeyService.delete(tenantId, id, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "API Key deleted successfully",
        description: "The API key has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listApiKeyQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: meApiKeyQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdApiKeyQueryKey(tenantId, id) });
    },
    onError: () => {
      toast({
        title: "Failed to delete API Key",
        description: "The API key has not been deleted successfully",
        variant: "destructive",
      });
    }
  })
}
