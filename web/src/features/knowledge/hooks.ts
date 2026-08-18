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

export const listKnowledgeQueryKey = (tenantId: string) => [BASE_QUERY_KEY, "list", tenantId] as const;
export const useListKnowledge = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedKnowledge>,
): UseQueryResult<PaginatedKnowledge, Error> => {
  return useQuery({
    queryKey: listKnowledgeQueryKey(tenantId),
    queryFn: () => knowledgeService.list(tenantId, opts),
    enabled: !!tenantId,
    refetchInterval: (query) => {
      const knowledgeList = query.state.data?.data ?? [];
      const isPending = knowledgeList.some(k => k.sync_status === "pending");
      if (isPending) return 1000;
      const isSyncing = knowledgeList.some(k => k.sync_status === "in_progress");
      return isSyncing ? 5000 : false;
    },
  });
};

export const knowledgeQueryKey = (tenantId: string, id: string) =>
  [BASE_QUERY_KEY, tenantId, id] as const;
export const useKnowledge = (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<Knowledge>,
): UseQueryResult<Knowledge, Error> => {
  return useQuery({
    queryKey: knowledgeQueryKey(tenantId, id),
    queryFn: () => knowledgeService.findByID(tenantId, id, opts),
    enabled: Boolean(tenantId) && Boolean(id),
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

export const listKnowledgeItemsQueryKey = (tenantId: string, id: string) =>
  [BASE_QUERY_KEY, tenantId, id, "items"] as const;
export const useListKnowledgeItems = (
  tenantId: string,
  id: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeItem>,
): UseQueryResult<PaginatedKnowledgeItem, Error> => {
  return useQuery({
    queryKey: [...listKnowledgeItemsQueryKey(tenantId, id), opts?.searchParams],
    queryFn: () => knowledgeService.listItems(tenantId, id, opts),
    enabled: Boolean(tenantId) && Boolean(id),
  });
};

export const createKnowledgeQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateKnowledge = (
  opts?: ServicePostOptions<CreateKnowledge, CreateKnowledgeResponse>,
): UseMutationResult<CreateKnowledgeResponse, Error, { tenantId: string; data: CreateKnowledge }> => {
  return useMutation({
    mutationKey: createKnowledgeQueryKey(),
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CreateKnowledge }) =>
      knowledgeService.create(tenantId, data, opts),
    onSuccess: ({ knowledge }, { tenantId }) => {
      toast({
        title: "Knowledge created successfully",
        description: "The knowledge source has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(tenantId, knowledge.id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey(tenantId) });
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
): UseMutationResult<void, Error, { tenantId: string; id: string; data: UpdateKnowledge }> => {
  return useMutation({
    mutationKey: updateKnowledgeQueryKey(),
    mutationFn: ({ tenantId, id, data }: { tenantId: string; id: string; data: UpdateKnowledge }) =>
      knowledgeService.update(tenantId, id, data, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "Knowledge updated successfully",
        description: "The knowledge source has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(tenantId, id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey(tenantId) });
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

export const listKnowledgeDriversQueryKey = (tenantId: string) =>
  [BASE_QUERY_KEY, "drivers", tenantId] as const;
export const useListKnowledgeDrivers = (
  tenantId: string,
  opts?: ServiceGetOptions<ListKnowledgeDrivers>,
): UseQueryResult<DriverInfo[], Error> => {
  return useQuery({
    queryKey: listKnowledgeDriversQueryKey(tenantId),
    queryFn: () => knowledgeService.listDrivers(tenantId, opts),
    enabled: !!tenantId,
  });
};

export const syncKnowledgeQueryKey = () => [BASE_QUERY_KEY, "sync"] as const;
export const useSyncKnowledge = (
  opts?: ServicePostOptions<undefined, void>,
): UseMutationResult<void, Error, { tenantId: string; id: string; force?: boolean }> => {
  return useMutation({
    mutationKey: syncKnowledgeQueryKey(),
    mutationFn: ({ tenantId, id, force = false }) =>
      knowledgeService.sync(tenantId, id, force, opts),
    onSuccess: (_data, { tenantId, id, force }) => {
      toast({
        title: force ? "Full resync started" : "Sync started",
        description: force
          ? "A full resync of the knowledge source has been queued"
          : "A sync of the knowledge source has been queued",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(tenantId, id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeItemsQueryKey(tenantId, id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey(tenantId) });
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
): UseMutationResult<void, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: deleteKnowledgeQueryKey(),
    mutationFn: ({ tenantId, id }: { tenantId: string; id: string }) =>
      knowledgeService.delete(tenantId, id, opts),
    onSuccess: (_data, { tenantId, id }) => {
      toast({
        title: "Knowledge deleted successfully",
        description: "The knowledge source has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeQueryKey(tenantId, id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeQueryKey(tenantId) });
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

export const listKnowledgeIndexesQueryKey = (tenantId: string, knowledgeID: string) =>
  [INDEX_QUERY_KEY, tenantId, knowledgeID] as const;
export const useListKnowledgeIndexes = (
  tenantId: string,
  knowledgeID: string,
  opts?: ServiceGetOptions<PaginatedKnowledgeIndex>,
): UseQueryResult<PaginatedKnowledgeIndex, Error> => {
  return useQuery({
    queryKey: listKnowledgeIndexesQueryKey(tenantId, knowledgeID),
    queryFn: () =>
      knowledgeIndexService.list(tenantId, knowledgeID, opts),
    enabled: Boolean(tenantId) && Boolean(knowledgeID),
  });
};

export const knowledgeIndexQueryKey = (tenantId: string, knowledgeID: string, id: string) =>
  [INDEX_QUERY_KEY, tenantId, knowledgeID, id] as const;
export const useKnowledgeIndex = (
  tenantId: string,
  knowledgeID: string,
  id: string,
  opts?: ServiceGetOptions<KnowledgeIndex>,
): UseQueryResult<KnowledgeIndex, Error> => {
  return useQuery({
    queryKey: knowledgeIndexQueryKey(tenantId, knowledgeID, id),
    queryFn: () =>
      knowledgeIndexService.findByID(tenantId, knowledgeID, id, opts),
    enabled: Boolean(tenantId) && Boolean(knowledgeID) && Boolean(id),
  });
};

export const createKnowledgeIndexQueryKey = () => [INDEX_QUERY_KEY, "create"] as const;
export const useCreateKnowledgeIndex = (
  opts?: ServicePostOptions<CreateKnowledgeIndex, KnowledgeIndex>,
): UseMutationResult<
  KnowledgeIndex,
  Error,
  { tenantId: string; knowledgeID: string; data: CreateKnowledgeIndex }
> => {
  return useMutation({
    mutationKey: createKnowledgeIndexQueryKey(),
    mutationFn: ({ tenantId, knowledgeID, data }) =>
      knowledgeIndexService.create(tenantId, knowledgeID, data, opts),
    onSuccess: (_data, { tenantId, knowledgeID }) => {
      toast({
        title: "Index created successfully",
        description: "The knowledge index has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listKnowledgeIndexesQueryKey(tenantId, knowledgeID) });
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
  { tenantId: string; knowledgeID: string; id: string; data: UpdateKnowledgeIndex }
> => {
  return useMutation({
    mutationKey: updateKnowledgeIndexQueryKey(),
    mutationFn: ({ tenantId, knowledgeID, id, data }) =>
      knowledgeIndexService.update(tenantId, knowledgeID, id, data, opts),
    onSuccess: (_data, { tenantId, knowledgeID, id }) => {
      toast({
        title: "Index updated successfully",
        description: "The knowledge index has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: knowledgeIndexQueryKey(tenantId, knowledgeID, id) });
      queryClient.invalidateQueries({ queryKey: listKnowledgeIndexesQueryKey(tenantId, knowledgeID) });
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
): UseMutationResult<void, Error, { tenantId: string; knowledgeID: string; id: string }> => {
  return useMutation({
    mutationKey: deleteKnowledgeIndexQueryKey(),
    mutationFn: ({ tenantId, knowledgeID, id }) =>
      knowledgeIndexService.delete(tenantId, knowledgeID, id, opts),
    onSuccess: (_data, { tenantId, knowledgeID }) => {
      toast({
        title: "Index deleted successfully",
        description: "The knowledge index has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listKnowledgeIndexesQueryKey(tenantId, knowledgeID) });
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

export const listKnowledgeIndexDriversQueryKey = (tenantId: string) =>
  [INDEX_QUERY_KEY, "drivers", tenantId] as const;
export const useListKnowledgeIndexDrivers = (
  tenantId: string,
  opts?: ServiceGetOptions<ListKnowledgeIndexDrivers>,
): UseQueryResult<DriverInfo[], Error> => {
  return useQuery({
    queryKey: listKnowledgeIndexDriversQueryKey(tenantId),
    queryFn: () => knowledgeIndexService.listDrivers(tenantId, opts),
    enabled: !!tenantId,
  });
};
