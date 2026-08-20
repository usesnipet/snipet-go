import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { sessionService } from "./service";

import type {
  CreateSession, FindSessionSearchParams, ListMessagesSearchParams, ListSessionSearchParams,
  PaginatedExecutionMessage, PaginatedSession, Session, UpdateSession
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "session";

export const listSessionQueryKey = (appCode: string) =>
  [BASE_QUERY_KEY, appCode] as const;
export const useListSession = (
  appCode: string,
  opts: Partial<ServiceGetOptions<PaginatedSession, ListSessionSearchParams>> = {},
): UseQueryResult<PaginatedSession, Error> => {
  return useQuery({
    queryKey: [...listSessionQueryKey(appCode), opts?.searchParams],
    queryFn: () =>
      sessionService.list(appCode, opts),
    enabled: !!appCode,
  })
}

export const findByIdSessionQueryKey = (appCode: string, id: string) =>
  [BASE_QUERY_KEY, appCode, "findById", id] as const;
export const useFindByIdSession = (
  appCode: string,
  id: string,
  opts: Partial<ServiceGetOptions<Session, FindSessionSearchParams>> = {},
): UseQueryResult<Session, Error> => {
  return useQuery({
    queryKey: [...findByIdSessionQueryKey(appCode, id), opts?.searchParams],
    queryFn: (): Promise<Session> =>
      sessionService.findById(appCode, id, opts),
    enabled: !!appCode && !!id,
  })
}

export const findMessagesSessionQueryKey = (appCode: string, id: string) =>
  [BASE_QUERY_KEY, appCode, "messages", id] as const;
export const useFindMessagesSession = (
  appCode: string,
  id: string,
  opts: ServiceGetOptions<PaginatedExecutionMessage, ListMessagesSearchParams> = {},
): UseQueryResult<PaginatedExecutionMessage, Error> => {
  return useQuery({
    queryKey: [...findMessagesSessionQueryKey(appCode, id), opts?.searchParams],
    queryFn: (): Promise<PaginatedExecutionMessage> =>
      sessionService.findMessages(appCode, id, opts),
    enabled: !!appCode && !!id,
  })
}

export const createSessionQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateSession = (
  opts: ServicePostOptions<CreateSession, Session> = {},
): UseMutationResult<
  Session,
  Error,
  { appCode: string; data: CreateSession }
> => {
  return useMutation({
    mutationKey: createSessionQueryKey(),
    mutationFn: ({ appCode, data }) =>
      sessionService.create(appCode, data, opts),
    onSuccess: (_data, { appCode }) => {
      toast({
        title: "Session created successfully",
        description: "The session has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listSessionQueryKey(appCode) });
    },
    onError: () => {
      toast({
        title: "Failed to create session",
        description: "The session has not been created successfully",
        variant: "destructive",
      });
    },
  })
}

export const updateSessionQueryKey = () => [BASE_QUERY_KEY, "update"] as const;
export const useUpdateSession = (
  opts: ServicePutOptions<UpdateSession, void> = {},
): UseMutationResult<
  void,
  Error,
  { appCode: string; id: string; data: UpdateSession }
> => {
  return useMutation({
    mutationKey: updateSessionQueryKey(),
    mutationFn: ({ appCode, id, data }) =>
      sessionService.update(appCode, id, data, opts),
    onSuccess: (_data, { appCode, id }) => {
      toast({
        title: "Session updated successfully",
        description: "The session has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listSessionQueryKey(appCode) });
      queryClient.invalidateQueries({ queryKey: findByIdSessionQueryKey(appCode, id) });
    },
    onError: () => {
      toast({
        title: "Failed to update session",
        description: "The session has not been updated successfully",
        variant: "destructive",
      });
    },
  })
}

export const deleteSessionQueryKey = () => [BASE_QUERY_KEY, "delete"] as const;
export const useDeleteSession = (
  opts: ServiceDeleteOptions<void> = {},
): UseMutationResult<void, Error, { appCode: string; id: string }> => {
  return useMutation({
    mutationKey: deleteSessionQueryKey(),
    mutationFn: ({ appCode, id }) =>
      sessionService.delete(appCode, id, opts),
    onSuccess: (_data, { appCode, id }) => {
      toast({
        title: "Session deleted successfully",
        description: "The session has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listSessionQueryKey(appCode) });
      queryClient.invalidateQueries({
        queryKey: findByIdSessionQueryKey(appCode, id),
      });
      queryClient.invalidateQueries({
        queryKey: findMessagesSessionQueryKey(appCode, id),
      });
    },
    onError: () => {
      toast({
        title: "Failed to delete session",
        description: "The session has not been deleted successfully",
        variant: "destructive",
      });
    },
  })
}
