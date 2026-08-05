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

export const listLlmQueryKey = () => [BASE_QUERY_KEY] as const;
export const useListLlm = (
  opts?: ServiceGetOptions<PaginatedLlm, ListLlmSearchParams>,
): UseQueryResult<PaginatedLlm, Error> => {
  return useQuery({
    queryKey: [...listLlmQueryKey(), opts?.searchParams],
    queryFn: () => llmService.list({ ...opts, auth: "api-key" }),
  });
};

export const createLlmQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateLlm = (
  opts?: ServicePostOptions<CreateLlm, Llm>,
): UseMutationResult<Llm, Error, CreateLlm> => {
  return useMutation({
    mutationKey: createLlmQueryKey(),
    mutationFn: (data: CreateLlm) =>
      llmService.create(data, { ...opts, auth: "api-key" }),
    onSuccess: () => {
      toast({
        title: "LLM created successfully",
        description: "The LLM has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey() });
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
): UseMutationResult<void, Error, { id: string; data: UpdateLlm }> => {
  return useMutation({
    mutationKey: updateLlmQueryKey(),
    mutationFn: ({ id, data }: { id: string; data: UpdateLlm }) =>
      llmService.update(id, data, { ...opts, auth: "api-key" }),
    onSuccess: () => {
      toast({
        title: "LLM updated successfully",
        description: "The LLM has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey() });
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
): UseMutationResult<void, Error, string> => {
  return useMutation({
    mutationKey: deleteLlmQueryKey(),
    mutationFn: (id: string) =>
      llmService.delete(id, { ...opts, auth: "api-key" }),
    onSuccess: () => {
      toast({
        title: "LLM deleted successfully",
        description: "The LLM has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey() });
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

export const listLlmDriversQueryKey = () => [BASE_QUERY_KEY, "drivers"] as const;
export const useListLlmDrivers = (
  opts?: ServiceGetOptions<ListDrivers>,
): UseQueryResult<ListDrivers, Error> => {
  return useQuery({
    queryKey: listLlmDriversQueryKey(),
    queryFn: () => llmService.listDrivers({ ...opts, auth: "api-key" }),
  });
};
