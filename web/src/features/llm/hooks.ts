import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { llmService } from "./service";

import type {
  CreateLlm, ListDrivers, ListLlmSearchParams, Llm, PaginatedLlm, UpdateLlm
} from "./schemas";
import type {
  ServiceDeleteOptions,
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "llm";

export const listLlmQueryKey = (tenantId: string) => [BASE_QUERY_KEY, "list", tenantId] as const;
export const useListLlm = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedLlm, ListLlmSearchParams>,
): UseQueryResult<PaginatedLlm, Error> => {
  return useQuery({
    queryKey: [...listLlmQueryKey(tenantId), opts?.searchParams],
    queryFn: () => llmService.list(tenantId, opts),
    enabled: !!tenantId,
  });
};

export const createLlmQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateLlm = (
  opts?: ServicePostOptions<CreateLlm, Llm>,
): UseMutationResult<Llm, Error, { tenantId: string; data: CreateLlm }> => {
  return useMutation({
    mutationKey: createLlmQueryKey(),
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CreateLlm }) =>
      llmService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "LLM created successfully",
        description: "The LLM has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to create LLM",
        description: "The LLM has not been created successfully",
        variant: "destructive",
      });
    },
  });
};

export const updateLlmQueryKey = () => [BASE_QUERY_KEY, "update"] as const;
export const useUpdateLlm = (
  opts?: ServicePutOptions<UpdateLlm, void>,
): UseMutationResult<void, Error, { tenantId: string; id: string; data: UpdateLlm }> => {
  return useMutation({
    mutationKey: updateLlmQueryKey(),
    mutationFn: ({ tenantId, id, data }: { tenantId: string; id: string; data: UpdateLlm }) =>
      llmService.update(tenantId, id, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "LLM updated successfully",
        description: "The LLM has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to update LLM",
        description: "The LLM has not been updated successfully",
        variant: "destructive",
      });
    },
  });
};

export const deleteLlmQueryKey = () => [BASE_QUERY_KEY, "delete"] as const;
export const useDeleteLlm = (
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: deleteLlmQueryKey(),
    mutationFn: ({ tenantId, id }: { tenantId: string; id: string }) =>
      llmService.delete(tenantId, id, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "LLM deleted successfully",
        description: "The LLM has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to delete LLM",
        description: "The LLM has not been deleted successfully",
        variant: "destructive",
      });
    },
  });
};

export const listLlmDriversQueryKey = (tenantId: string) =>
  [BASE_QUERY_KEY, "drivers", tenantId] as const;
export const useListLlmDrivers = (
  tenantId: string,
  opts?: ServiceGetOptions<ListDrivers>,
): UseQueryResult<ListDrivers, Error> => {
  return useQuery({
    queryKey: listLlmDriversQueryKey(tenantId),
    queryFn: () => llmService.listDrivers(tenantId, opts),
    enabled: !!tenantId,
  });
};
