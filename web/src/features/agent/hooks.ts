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

export const listAgentQueryKey = () => [BASE_QUERY_KEY] as const;
export const useListAgent = (
  opts?: ServiceGetOptions<PaginatedAgent>,
): UseQueryResult<PaginatedAgent, Error> => {
  return useQuery({
    queryKey: listAgentQueryKey(),
    queryFn: () => agentService.list({ ...opts, auth: "api-key" }),
  });
};

export const findByIdAgentQueryKey = (id: string) =>
  [BASE_QUERY_KEY, "findById", id] as const;
export const useFindByIdAgent = (
  id: string,
  opts?: ServiceGetOptions<Agent>,
): UseQueryResult<Agent, Error> => {
  return useQuery({
    queryKey: findByIdAgentQueryKey(id),
    queryFn: (): Promise<Agent> =>
      agentService.findById(id, { ...opts, auth: "api-key" }),
    enabled: !!id,
  });
};

export const createAgentQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateAgent = (
  opts?: ServicePostOptions<CreateAgent, Agent>,
): UseMutationResult<Agent, Error, CreateAgent> => {
  return useMutation({
    mutationKey: createAgentQueryKey(),
    mutationFn: (data: CreateAgent) =>
      agentService.create(data, { ...opts, auth: "api-key" }),
    onSuccess: () => {
      toast({
        title: "Agent created successfully",
        description: "The agent has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAgentQueryKey() });
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
): UseMutationResult<void, Error, { id: string; data: UpdateAgent }> => {
  return useMutation({
    mutationKey: updateAgentQueryKey(),
    mutationFn: ({ id, data }: { id: string; data: UpdateAgent }) =>
      agentService.update(id, data, { ...opts, auth: "api-key" }),
    onSuccess: (_data, { id }) => {
      toast({
        title: "Agent updated successfully",
        description: "The agent has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAgentQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdAgentQueryKey(id) });
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
): UseMutationResult<void, Error, string> => {
  return useMutation({
    mutationKey: deleteAgentQueryKey(),
    mutationFn: (id: string) =>
      agentService.delete(id, { ...opts, auth: "api-key" }),
    onSuccess: (_data, id) => {
      toast({
        title: "Agent deleted successfully",
        description: "The agent has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAgentQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIdAgentQueryKey(id) });
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
