import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { knowledgeIndexService, knowledgeService } from "./service";

import type {
  CreateKnowledge,
  CreateKnowledgeIndex,
  CreateKnowledgeResponse,
  Knowledge,
  KnowledgeIndex,
  ListKnowledgeDrivers,
  ListKnowledgeIndexDrivers,
  PaginatedKnowledge,
  PaginatedKnowledgeIndex,
  PaginatedKnowledgeItem,
  UpdateKnowledge,
  UpdateKnowledgeIndex,
} from "./schemas";
import type {
  ServiceDeleteOptions,
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";
import type { DriverInfo } from "@/schemas/driver";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "knowledge";

export const listKnowledgeQueryKey = () => [BASE_QUERY_KEY] as const;
export const useListKnowledge = (
  opts?: ServiceGetOptions<PaginatedKnowledge>,
): UseQueryResult<PaginatedKnowledge, Error> => {
  return useQuery({
    queryKey: listKnowledgeQueryKey(),
    queryFn: () => knowledgeService.list(opts),
    refetchInterval: (query) => {
      const knowledgeList = query.state.data?.data ?? [];
      const isPending = knowledgeList.some(k => k.sync_status === "pending");
      if (isPending) return 1000;
      const isSyncing = knowledgeList.some(k => k.sync_status === "in_progress");
      return isSyncing ? 5000 : false;
    },
  });
};

export const knowledgeQueryKey = (id: string) => [BASE_QUERY_KEY, id] as const;
export const useKnowledge = (
  id: string,
  opts?: ServiceGetOptions<Knowledge>,
): UseQueryResult<Knowledge, Error> => {
  return useQuery({
    queryKey: knowledgeQueryKey(id),
    queryFn: () => knowledgeService.findByID(id, opts),
    enabled: Boolean(id),
    refetchInterval: (query) => {
      const knowledge = query.state.data;
      if (!knowledge) return false;
      const isPending = knowledge.sync_status === "pending";
      if (isPending) return 1000;
      const isSyncing = knowledge.sync_status === "in_progress";
      return isSyncing ? 5000 : false;
    },
  });
};

export const listKnowledgeItemsQueryKey = (id: string) =>
  [BASE_QUERY_KEY, id, "items"] as const;
export const useListKnowledgeItems = (
  id: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeItem>,
): UseQueryResult<PaginatedKnowledgeItem, Error> => {
  return useQuery({
    queryKey: [...listKnowledgeItemsQueryKey(id), opts?.searchParams],
    queryFn: () => knowledgeService.listItems(id, opts),
    enabled: Boolean(id),
  });
};

export const createKnowledgeQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateKnowledge = (
  opts?: ServicePostOptions<CreateKnowledge, CreateKnowledgeResponse>,
): UseMutationResult<CreateKnowledgeResponse, Error, CreateKnowledge> => {
  return useMutation({
    mutationKey: createKnowledgeQueryKey(),
    mutationFn: (data: CreateKnowledge) =>
      knowledgeService.create(data, opts),
    onSuccess: ({ knowledge }) => {
      toast({
        title: "Knowledge created successfully",
        description: "The knowledge source has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(knowledge.id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey() });
    },
    onError: () => {
      toast({
        title: "Failed to create knowledge",
        description: "The knowledge source has not been created successfully",
        variant: "destructive",
      });
    },
  });
};

export const updateKnowledgeQueryKey = () => [BASE_QUERY_KEY, "update"] as const;
export const useUpdateKnowledge = (
  opts?: ServicePutOptions<UpdateKnowledge, void>,
): UseMutationResult<void, Error, { id: string; data: UpdateKnowledge }> => {
  return useMutation({
    mutationKey: updateKnowledgeQueryKey(),
    mutationFn: ({ id, data }: { id: string; data: UpdateKnowledge }) =>
      knowledgeService.update(id, data, opts),
    onSuccess: (_data, { id }) => {
      toast({
        title: "Knowledge updated successfully",
        description: "The knowledge source has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey() });
    },
    onError: () => {
      toast({
        title: "Failed to update knowledge",
        description: "The knowledge source has not been updated successfully",
        variant: "destructive",
      });
    },
  });
};

export const listKnowledgeDriversQueryKey = () =>
  [BASE_QUERY_KEY, "drivers"] as const;
export const useListKnowledgeDrivers = (
  opts?: ServiceGetOptions<ListKnowledgeDrivers>,
): UseQueryResult<DriverInfo[], Error> => {
  return useQuery({
    queryKey: listKnowledgeDriversQueryKey(),
    queryFn: () => knowledgeService.listDrivers(opts),
  });
};

export const syncKnowledgeQueryKey = () => [BASE_QUERY_KEY, "sync"] as const;
export const useSyncKnowledge = (
  opts?: ServicePostOptions<undefined, void>,
): UseMutationResult<void, Error, { id: string; force?: boolean }> => {
  return useMutation({
    mutationKey: syncKnowledgeQueryKey(),
    mutationFn: ({ id, force = false }) =>
      knowledgeService.sync(id, force, opts),
    onSuccess: (_data, { id, force }) => {
      toast({
        title: force ? "Full resync started" : "Sync started",
        description: force
          ? "A full resync of the knowledge source has been queued"
          : "A sync of the knowledge source has been queued",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeItemsQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey() });
    },
    onError: (_error, { force }) => {
      toast({
        title: force ? "Failed to start full resync" : "Failed to start sync",
        description: "The knowledge source sync could not be started",
        variant: "destructive",
      });
    },
  });
};

export const deleteKnowledgeQueryKey = () => [BASE_QUERY_KEY, "delete"] as const;
export const useDeleteKnowledge = (
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, string> => {
  return useMutation({
    mutationKey: deleteKnowledgeQueryKey(),
    mutationFn: (id: string) =>
      knowledgeService.delete(id, opts),
    onSuccess: (_data, id) => {
      toast({
        title: "Knowledge deleted successfully",
        description: "The knowledge source has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey() });
    },
    onError: () => {
      toast({
        title: "Failed to delete knowledge",
        description: "The knowledge source has not been deleted successfully",
        variant: "destructive",
      });
    },
  });
};

const INDEX_QUERY_KEY = "knowledge-index";

export const listKnowledgeIndexesQueryKey = (knowledgeID: string) =>
  [INDEX_QUERY_KEY, knowledgeID] as const;
export const useListKnowledgeIndexes = (
  knowledgeID: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeIndex>,
): UseQueryResult<PaginatedKnowledgeIndex, Error> => {
  return useQuery({
    queryKey: listKnowledgeIndexesQueryKey(knowledgeID),
    queryFn: () =>
      knowledgeIndexService.list(knowledgeID, opts),
    enabled: Boolean(knowledgeID),
  });
};

export const knowledgeIndexQueryKey = (knowledgeID: string, id: string) =>
  [INDEX_QUERY_KEY, knowledgeID, id] as const;
export const useKnowledgeIndex = (
  knowledgeID: string,
  id: string,
  opts?: ServiceGetOptions<KnowledgeIndex>,
): UseQueryResult<KnowledgeIndex, Error> => {
  return useQuery({
    queryKey: knowledgeIndexQueryKey(knowledgeID, id),
    queryFn: () =>
      knowledgeIndexService.findByID(knowledgeID, id, opts),
    enabled: Boolean(knowledgeID) && Boolean(id),
  });
};

export const createKnowledgeIndexQueryKey = () => [INDEX_QUERY_KEY, "create"] as const;
export const useCreateKnowledgeIndex = (
  opts?: ServicePostOptions<CreateKnowledgeIndex, KnowledgeIndex>,
): UseMutationResult<KnowledgeIndex, Error, { knowledgeID: string; data: CreateKnowledgeIndex }> => {
  return useMutation({
    mutationKey: createKnowledgeIndexQueryKey(),
    mutationFn: ({ knowledgeID, data }) =>
      knowledgeIndexService.create(knowledgeID, data, opts),
    onSuccess: (_data, { knowledgeID }) => {
      toast({
        title: "Index created successfully",
        description: "The knowledge index has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listKnowledgeIndexesQueryKey(knowledgeID) });
    },
    onError: () => {
      toast({
        title: "Failed to create index",
        description: "The knowledge index has not been created successfully",
        variant: "destructive",
      });
    },
  });
};

export const updateKnowledgeIndexQueryKey = () => [INDEX_QUERY_KEY, "update"] as const;
export const useUpdateKnowledgeIndex = (
  opts?: ServicePutOptions<UpdateKnowledgeIndex, void>,
): UseMutationResult<
  void,
  Error,
  { knowledgeID: string; id: string; data: UpdateKnowledgeIndex }
> => {
  return useMutation({
    mutationKey: updateKnowledgeIndexQueryKey(),
    mutationFn: ({ knowledgeID, id, data }) =>
      knowledgeIndexService.update(knowledgeID, id, data, opts),
    onSuccess: (_data, { knowledgeID, id }) => {
      toast({
        title: "Index updated successfully",
        description: "The knowledge index has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeIndexQueryKey(knowledgeID, id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeIndexesQueryKey(knowledgeID) });
    },
    onError: () => {
      toast({
        title: "Failed to update index",
        description: "The knowledge index has not been updated successfully",
        variant: "destructive",
      });
    },
  });
};

export const deleteKnowledgeIndexQueryKey = () => [INDEX_QUERY_KEY, "delete"] as const;
export const useDeleteKnowledgeIndex = (
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, { knowledgeID: string; id: string }> => {
  return useMutation({
    mutationKey: deleteKnowledgeIndexQueryKey(),
    mutationFn: ({ knowledgeID, id }) =>
      knowledgeIndexService.delete(knowledgeID, id, opts),
    onSuccess: (_data, { knowledgeID }) => {
      toast({
        title: "Index deleted successfully",
        description: "The knowledge index has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listKnowledgeIndexesQueryKey(knowledgeID) });
    },
    onError: () => {
      toast({
        title: "Failed to delete index",
        description: "The knowledge index has not been deleted successfully",
        variant: "destructive",
      });
    },
  });
};

export const listKnowledgeIndexDriversQueryKey = () => [INDEX_QUERY_KEY, "drivers"] as const;
export const useListKnowledgeIndexDrivers = (
  opts?: ServiceGetOptions<ListKnowledgeIndexDrivers>,
): UseQueryResult<DriverInfo[], Error> => {
  return useQuery({
    queryKey: listKnowledgeIndexDriversQueryKey(),
    queryFn: () => knowledgeIndexService.listDrivers(opts),
  });
};
