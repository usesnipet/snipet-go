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

export const listSessionQueryKey = (clientCode: string) =>
  [BASE_QUERY_KEY, clientCode] as const;
export const useListSession = (
  clientCode: string,
  opts?: Partial<ServiceGetOptions<PaginatedSession, ListSessionSearchParams>>,
): UseQueryResult<PaginatedSession, Error> => {
  return useQuery({
    queryKey: [...listSessionQueryKey(clientCode), opts?.searchParams],
    queryFn: () =>
      sessionService.list(clientCode, { auth: "api-key", ...opts }),
    enabled: !!clientCode,
  })
}

export const findByIdSessionQueryKey = (clientCode: string, id: string) =>
  [BASE_QUERY_KEY, clientCode, "findById", id] as const;
export const useFindByIdSession = (
  clientCode: string,
  id: string,
  opts?: Partial<ServiceGetOptions<Session, FindSessionSearchParams>>,
): UseQueryResult<Session, Error> => {
  return useQuery({
    queryKey: [...findByIdSessionQueryKey(clientCode, id), opts?.searchParams],
    queryFn: (): Promise<Session> =>
      sessionService.findById(clientCode, id, { auth: "api-key", ...opts }),
    enabled: !!clientCode && !!id,
  })
}

export const findMessagesSessionQueryKey = (clientCode: string, id: string) =>
  [BASE_QUERY_KEY, clientCode, "messages", id] as const;
export const useFindMessagesSession = (
  clientCode: string,
  id: string,
  opts?: ServiceGetOptions<PaginatedExecutionMessage, ListMessagesSearchParams>,
): UseQueryResult<PaginatedExecutionMessage, Error> => {
  return useQuery({
    queryKey: [...findMessagesSessionQueryKey(clientCode, id), opts?.searchParams],
    queryFn: (): Promise<PaginatedExecutionMessage> =>
      sessionService.findMessages(clientCode, id, { auth: "api-key", ...opts }),
    enabled: !!clientCode && !!id,
  })
}

export const createSessionQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateSession = (
  opts?: ServicePostOptions<CreateSession, Session>,
): UseMutationResult<
  Session,
  Error,
  { clientCode: string; data: CreateSession }
> => {
  return useMutation({
    mutationKey: createSessionQueryKey(),
    mutationFn: ({ clientCode, data }) =>
      sessionService.create(clientCode, data, { auth: "api-key", ...opts }),
    onSuccess: (_data, { clientCode }) => {
      toast({
        title: "Session created successfully",
        description: "The session has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listSessionQueryKey(clientCode) });
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
  opts?: ServicePutOptions<UpdateSession, void>,
): UseMutationResult<
  void,
  Error,
  { clientCode: string; id: string; data: UpdateSession }
> => {
  return useMutation({
    mutationKey: updateSessionQueryKey(),
    mutationFn: ({ clientCode, id, data }) =>
      sessionService.update(clientCode, id, data, { auth: "api-key", ...opts }),
    onSuccess: (_data, { clientCode, id }) => {
      toast({
        title: "Session updated successfully",
        description: "The session has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listSessionQueryKey(clientCode) });
      queryClient.invalidateQueries({ queryKey: findByIdSessionQueryKey(clientCode, id) });
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
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, { clientCode: string; id: string }> => {
  return useMutation({
    mutationKey: deleteSessionQueryKey(),
    mutationFn: ({ clientCode, id }) =>
      sessionService.delete(clientCode, id, { auth: "api-key", ...opts }),
    onSuccess: (_data, { clientCode, id }) => {
      toast({
        title: "Session deleted successfully",
        description: "The session has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listSessionQueryKey(clientCode) });
      queryClient.invalidateQueries({
        queryKey: findByIdSessionQueryKey(clientCode, id),
      });
      queryClient.invalidateQueries({
        queryKey: findMessagesSessionQueryKey(clientCode, id),
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
