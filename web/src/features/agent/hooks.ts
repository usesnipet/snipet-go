import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { agentService } from "./service";

import type { Agent, CreateAgent, PaginatedAgent, UpdateAgent } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "agent";

export const listAgentQueryKey = (tenantId: string) => [BASE_QUERY_KEY, "list", tenantId] as const;
export const useListAgent = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedAgent>,
): UseQueryResult<PaginatedAgent, Error> => {
  return useQuery({
    queryKey: listAgentQueryKey(tenantId),
    queryFn: () => agentService.list(tenantId, opts),
    enabled: !!tenantId,
  });
};

export const findByIdAgentQueryKey = (tenantId: string, id: string) =>
  [BASE_QUERY_KEY, "findById", tenantId, id] as const;
export const useFindByIdAgent = (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<Agent>,
): UseQueryResult<Agent, Error> => {
  return useQuery({
    queryKey: findByIdAgentQueryKey(tenantId, id),
    queryFn: (): Promise<Agent> =>
      agentService.findById(tenantId, id, opts),
    enabled: !!tenantId && !!id,
  });
};

export const createAgentQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateAgent = (
  opts?: ServicePostOptions<CreateAgent, Agent>,
): UseMutationResult<Agent, Error, { tenantId: string; data: CreateAgent }> => {
  return useMutation({
    mutationKey: createAgentQueryKey(),
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CreateAgent }) =>
      agentService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Agent created successfully",
        description: "The agent has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAgentQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to create agent",
        description: "The agent has not been created successfully",
        variant: "destructive",
      });
    },
  });
};

export const updateAgentQueryKey = () => [BASE_QUERY_KEY, "update"] as const;
export const useUpdateAgent = (
  opts?: ServicePutOptions<UpdateAgent, void>,
): UseMutationResult<void, Error, { tenantId: string; id: string; data: UpdateAgent }> => {
  return useMutation({
    mutationKey: updateAgentQueryKey(),
    mutationFn: ({ tenantId, id, data }: { tenantId: string; id: string; data: UpdateAgent }) =>
      agentService.update(tenantId, id, data, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "Agent updated successfully",
        description: "The agent has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAgentQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByIdAgentQueryKey(tenantId, id) });
    },
    onError: () => {
      toast({
        title: "Failed to update agent",
        description: "The agent has not been updated successfully",
        variant: "destructive",
      });
    },
  });
};

export const deleteAgentQueryKey = () => [BASE_QUERY_KEY, "delete"] as const;
export const useDeleteAgent = (
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: deleteAgentQueryKey(),
    mutationFn: ({ tenantId, id }: { tenantId: string; id: string }) =>
      agentService.delete(tenantId, id, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "Agent deleted successfully",
        description: "The agent has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAgentQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByIdAgentQueryKey(tenantId, id) });
    },
    onError: () => {
      toast({
        title: "Failed to delete agent",
        description: "The agent has not been deleted successfully",
        variant: "destructive",
      });
    },
  });
};
