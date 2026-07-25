import { useNavigate } from "@/hooks/use-navigate";
import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { apiKeyService } from "./service";
import { useApiKeyStore } from "./store";

import type { Paginated } from "@/schemas/paginated";
import type { ApiKey, ApiKeyKey, ApiKeyWithSecret, CreateApiKey, UpdateApiKeyExpiration } from "./schemas";
import type { ServiceGetOptions, ServicePostOptions, ServicePutOptions } from "@/lib/services";
const BASE_QUERY_KEY = "api-key";

export const listApiKeyQueryKey = () => [BASE_QUERY_KEY];
export const useListApiKey = (opts?: ServiceGetOptions<Paginated<ApiKey>>) => {
  return useQuery({
    queryKey: listApiKeyQueryKey(),
    queryFn: () => apiKeyService.list({ auth: "api-key" , ...opts }),
  })
}

export const createApiKeyQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateApiKey = (opts?: ServicePostOptions<CreateApiKey, ApiKeyWithSecret>) => {
  return useMutation({
    mutationKey: createApiKeyQueryKey(),
    mutationFn: (data: CreateApiKey) => apiKeyService.create(data, { auth: "api-key" , ...opts }),
    onSuccess: () => {
      toast({
        title: "API Key created successfully",
        description: "The API key has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listApiKeyQueryKey() });
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
export const useMeApiKey = (opts?: ServiceGetOptions<ApiKey>) => {
  return useQuery({
    queryKey: meApiKeyQueryKey(),
    queryFn: () => apiKeyService.me({ auth: "api-key" , ...opts }),
  })
}

export const findByIdApiKeyQueryKey = (id: string) => [BASE_QUERY_KEY, "findById", id];
export const useFindByIdApiKey = (id: string, opts?: ServiceGetOptions<ApiKey>) => {
  return useQuery({
    queryKey: findByIdApiKeyQueryKey(id),
    queryFn: () => apiKeyService.findById(id, { auth: "api-key" , ...opts }),
  })
}

export const rollApiKeyQueryKey = (id: string) => [BASE_QUERY_KEY, "roll", id];
export const useRollApiKey = (id: string, opts?: ServicePostOptions<undefined, ApiKeyWithSecret>) => {
  return useMutation({
    mutationKey: rollApiKeyQueryKey(id),
    mutationFn: () => apiKeyService.roll(id, { auth: "api-key" , ...opts }),
    onSuccess: () => {
      toast({
        title: "API Key rolled successfully",
        description: "The API key has been rolled successfully",
      });
      queryClient.invalidateQueries({ queryKey: meApiKeyQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdApiKeyQueryKey(id) });
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

export const updateExpirationApiKeyQueryKey = (id: string) => [BASE_QUERY_KEY, "updateExpiration", id];
export const useUpdateExpirationApiKey = (
  id: string,
  opts?: ServicePutOptions<UpdateApiKeyExpiration, ApiKeyWithSecret>
) => {
  return useMutation({
    mutationKey: updateExpirationApiKeyQueryKey(id),
    mutationFn: (data: UpdateApiKeyExpiration) => apiKeyService.updateExpiration(id, data, { auth: "api-key" , ...opts }),
    onSuccess: () => {
      toast({
        title: "API Key expiration updated successfully",
        description: "The API key expiration has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: meApiKeyQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdApiKeyQueryKey(id) });
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

export const apiKeyLoginQueryKey = () => [BASE_QUERY_KEY, "login"];
export const useApiKeyLogin = (opts?: ServiceGetOptions<ApiKey>) => {
  const navigate = useNavigate();
  const setApiKey = useApiKeyStore((state) => state.set)
  return useMutation({
    mutationKey: apiKeyLoginQueryKey(),
    mutationFn: async (data: ApiKeyKey) => apiKeyService.me({
      ...opts,
      headers: {
        "X-API-Key": data,
        ...opts?.headers,
      },
    }),
    onSuccess: (_, variables) => {
      setApiKey(variables)
      toast({
        title: "API Key logged in successfully",
        description: "The API key has been logged in successfully",
      });
      navigate("/admin");
    },
    onError: () => {
      toast({
        title: "Failed to login API key",
        description: "The API key has not been logged in successfully",
        variant: "destructive",
      });
    }
  })
}