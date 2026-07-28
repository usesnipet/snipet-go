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

export const listClientQueryKey = () => [BASE_QUERY_KEY] as const;
export const useListClient = (
  opts?: ServiceGetOptions<PaginatedClient>
): UseQueryResult<PaginatedClient, Error> => {
  return useQuery({
    queryKey: listClientQueryKey(),
    queryFn: () => clientService.list({ ...opts, auth: "api-key" }),
  })
}

export const listClientAgentsQueryKey = (code: string) => [BASE_QUERY_KEY, "agents", code] as const;
export const useListClientAgents = (
  code: string,
  opts?: ServiceGetOptions<PaginatedAgent>
): UseQueryResult<PaginatedAgent, Error> => {
  return useQuery({
    queryKey: listClientAgentsQueryKey(code),
    queryFn: () => clientService.listAgents(code, { auth: "jwt", ...opts }),
    enabled: !!code,
  })
}

export const findByCodeClientQueryKey = (code: string) => [BASE_QUERY_KEY, "findByCode", code];
export const useFindByCodeClient = (
  code: string,
  opts?: ServiceGetOptions<Client>
): UseQueryResult<Client, Error> => {
  return useQuery({
    queryKey: findByCodeClientQueryKey(code),
    queryFn: (): Promise<Client> =>
      clientService.findByCode(code, { ...opts, auth: "api-key" }),
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
      clientService.findPublicByCode(code, { ...opts, auth: false }),
  })
}

export const createClientQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateClient = (
  opts?: ServicePostOptions<CreateClient, Client>
): UseMutationResult<Client, Error, CreateClient> => {
  return useMutation({
    mutationKey: createClientQueryKey(),
    mutationFn: (data: CreateClient) =>
      clientService.create(data, { ...opts, auth: "api-key" }),
    onSuccess: () => {
      toast({
        title: "Client created successfully",
        description: "The client has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listClientQueryKey() });
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
  opts?: ServicePutOptions<UpdateClient, void>
): UseMutationResult<void, Error, { code: string; data: UpdateClient }> => {
  return useMutation({
    mutationKey: updateClientQueryKey(),
    mutationFn: ({ code, data }: { code: string; data: UpdateClient }) =>
      clientService.update(code, data, { ...opts, auth: "api-key" }),
    onSuccess: (_data, { code }) => {
      toast({
        title: "Client updated successfully",
        description: "The client has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listClientQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeClientQueryKey(code) });
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
  opts?: ServiceDeleteOptions<void>
): UseMutationResult<void, Error, string> => {
  return useMutation({
    mutationKey: deleteClientQueryKey(),
    mutationFn: (code: string) =>
      clientService.delete(code, { ...opts, auth: "api-key" }),
    onSuccess: (_data, code) => {
      toast({
        title: "Client deleted successfully",
        description: "The client has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listClientQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeClientQueryKey(code) });
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
