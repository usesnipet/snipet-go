import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { clientService } from "./service";

import type {
  Client, ClientPublic, CreateClient, PaginatedClient, UpdateClient
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";
import type { PaginatedAgent } from "../agent/schemas";

const BASE_QUERY_KEY = "client";

export const listClientQueryKey = (tenantId: string) => [BASE_QUERY_KEY, "list", tenantId] as const;
export const useListClient = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedClient>
): UseQueryResult<PaginatedClient, Error> => {
  return useQuery({
    queryKey: listClientQueryKey(tenantId),
    queryFn: () => clientService.list(tenantId, opts),
    enabled: !!tenantId,
  })
}

export const listClientAgentsQueryKey = (code: string) => [BASE_QUERY_KEY, "agents", code] as const;
export const useListClientAgents = (
  code: string,
  opts?: ServiceGetOptions<PaginatedAgent>
): UseQueryResult<PaginatedAgent, Error> => {
  return useQuery({
    queryKey: listClientAgentsQueryKey(code),
    queryFn: () => clientService.listAgents(code, opts),
    enabled: !!code,
  })
}

export const findByCodeClientQueryKey = (tenantId: string, code: string) =>
  [BASE_QUERY_KEY, "findByCode", tenantId, code];
export const useFindByCodeClient = (
  tenantId: string,
  code: string,
  opts?: ServiceGetOptions<Client>
): UseQueryResult<Client, Error> => {
  return useQuery({
    queryKey: findByCodeClientQueryKey(tenantId, code),
    queryFn: (): Promise<Client> =>
      clientService.findByCode(tenantId, code, opts),
    enabled: !!tenantId && !!code,
  })
}

export const findPublicByCodeClientQueryKey = (code: string) =>
  [BASE_QUERY_KEY, "findPublicByCode", code];
export const useFindPublicByCodeClient = (
  code: string,
  opts?: ServiceGetOptions<ClientPublic>
): UseQueryResult<ClientPublic, Error> => {
  return useQuery({
    queryKey: findPublicByCodeClientQueryKey(code),
    queryFn: (): Promise<ClientPublic> =>
      clientService.findPublicByCode(code, opts),
    enabled: !!code,
  })
}

export const createClientQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateClient = (
  opts?: ServicePostOptions<CreateClient, Client>
): UseMutationResult<Client, Error, { tenantId: string; data: CreateClient }> => {
  return useMutation({
    mutationKey: createClientQueryKey(),
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CreateClient }) =>
      clientService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Client created successfully",
        description: "The client has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listClientQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to create client",
        description: "The client has not been created successfully",
        variant: "destructive",
      });
    }
  })
}

export const updateClientQueryKey = () => [BASE_QUERY_KEY, "update"];
export const useUpdateClient = (
  opts: ServicePutOptions<UpdateClient, void> = {}
): UseMutationResult<void, Error, { tenantId: string; code: string; data: UpdateClient }> => {
  return useMutation({
    mutationKey: updateClientQueryKey(),
    mutationFn: ({ tenantId, code, data }: { tenantId: string; code: string; data: UpdateClient }) =>
      clientService.update(tenantId, code, data, opts),
    onSuccess: (_data, { tenantId, code }) => {
      toast({
        title: "Client updated successfully",
        description: "The client has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listClientQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeClientQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to update client",
        description: "The client has not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const deleteClientQueryKey = () => [BASE_QUERY_KEY, "delete"];
export const useDeleteClient = (
  opts: ServiceDeleteOptions<void> = {}
): UseMutationResult<void, Error, { tenantId: string; code: string }> => {
  return useMutation({
    mutationKey: deleteClientQueryKey(),
    mutationFn: ({ tenantId, code }: { tenantId: string; code: string }) =>
      clientService.delete(tenantId, code, opts),
    onSuccess: (_data, { tenantId, code }) => {
      toast({
        title: "Client deleted successfully",
        description: "The client has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listClientQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeClientQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to delete client",
        description: "The client has not been deleted successfully",
        variant: "destructive",
      });
    }
  })
}
