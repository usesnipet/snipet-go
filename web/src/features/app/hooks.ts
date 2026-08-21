import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { appService } from "./service";

import type {
  App, AppWithSecret, CreateApp, PaginatedApp, UpdateApp, UpdateAppAuthConfig
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "app";

export const listAppQueryKey = (tenantId: string) => [BASE_QUERY_KEY, "list", tenantId] as const;
export const useListApp = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedApp>
): UseQueryResult<PaginatedApp, Error> => {
  return useQuery({
    queryKey: listAppQueryKey(tenantId),
    queryFn: () => appService.list(tenantId, opts),
    enabled: !!tenantId,
  })
}

export const findByCodeAppQueryKey = (tenantId: string, code: string) =>
  [BASE_QUERY_KEY, "findByCode", tenantId, code];
export const useFindByCodeApp = (
  tenantId: string,
  code: string,
  opts?: ServiceGetOptions<App>
): UseQueryResult<App, Error> => {
  return useQuery({
    queryKey: findByCodeAppQueryKey(tenantId, code),
    queryFn: (): Promise<App> =>
      appService.findByCode(tenantId, code, opts),
    enabled: !!tenantId && !!code,
  })
}

export const createAppQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateApp = (
  opts?: ServicePostOptions<CreateApp, AppWithSecret>
): UseMutationResult<AppWithSecret, Error, { tenantId: string; data: CreateApp }> => {
  return useMutation({
    mutationKey: createAppQueryKey(),
    mutationFn: ({ tenantId, data }: { tenantId: string; data: CreateApp }) =>
      appService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "App created successfully",
        description: "The app has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to create app",
        description: "The app has not been created successfully",
        variant: "destructive",
      });
    }
  })
}

export const updateAppQueryKey = () => [BASE_QUERY_KEY, "update"];
export const useUpdateApp = (
  opts: ServicePutOptions<UpdateApp, void> = {}
): UseMutationResult<void, Error, { tenantId: string; code: string; data: UpdateApp }> => {
  return useMutation({
    mutationKey: updateAppQueryKey(),
    mutationFn: ({ tenantId, code, data }: { tenantId: string; code: string; data: UpdateApp }) =>
      appService.update(tenantId, code, data, opts),
    onSuccess: (_data, { tenantId, code }) => {
      toast({
        title: "App updated successfully",
        description: "The app has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to update app",
        description: "The app has not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const updateAppAuthConfigQueryKey = () => [BASE_QUERY_KEY, "update-auth-config"];
export const useUpdateAppAuthConfig = (
  opts: ServicePutOptions<UpdateAppAuthConfig, void> = {}
): UseMutationResult<void, Error, { tenantId: string; code: string; data: UpdateAppAuthConfig }> => {
  return useMutation({
    mutationKey: updateAppAuthConfigQueryKey(),
    mutationFn: ({ tenantId, code, data }: { tenantId: string; code: string; data: UpdateAppAuthConfig }) =>
      appService.updateAuthConfig(tenantId, code, data, opts),
    onSuccess: (_data, { tenantId, code }) => {
      toast({
        title: "App updated successfully",
        description: "The app has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to update app",
        description: "The app has not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const rollAppQueryKey = () => [BASE_QUERY_KEY, "roll"];
export const useRollApp = (
  opts?: ServicePostOptions<undefined, AppWithSecret>
): UseMutationResult<AppWithSecret, Error, { tenantId: string; code: string }> => {
  return useMutation({
    mutationKey: rollAppQueryKey(),
    mutationFn: ({ tenantId, code }: { tenantId: string; code: string }) =>
      appService.roll(tenantId, code, opts),
    onSuccess: (_data, { tenantId, code }) => {
      toast({
        title: "App key rolled successfully",
        description: "The app key has been rolled successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to roll app key",
        description: "The app key has not been rolled successfully",
        variant: "destructive",
      });
    }
  })
}

export const setActiveAppQueryKey = () => [BASE_QUERY_KEY, "set-active"];
export const useSetActiveApp = (
  opts: ServicePutOptions<void, void> = {}
): UseMutationResult<void, Error, { tenantId: string; code: string; active: boolean }> => {
  return useMutation({
    mutationKey: setActiveAppQueryKey(),
    mutationFn: ({ tenantId, code, active }: { tenantId: string; code: string; active: boolean }) =>
      appService.setActive(tenantId, code, active, opts),
    onSuccess: (_data, { tenantId, code, active }) => {
      toast({
        title: active ? "App activated" : "App deactivated",
        description: active
          ? "The app has been activated successfully"
          : "The app has been deactivated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to update app status",
        description: "The app status has not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const deleteAppQueryKey = () => [BASE_QUERY_KEY, "delete"];
export const useDeleteApp = (
  opts: ServiceDeleteOptions<void> = {}
): UseMutationResult<void, Error, { tenantId: string; code: string }> => {
  return useMutation({
    mutationKey: deleteAppQueryKey(),
    mutationFn: ({ tenantId, code }: { tenantId: string; code: string }) =>
      appService.delete(tenantId, code, opts),
    onSuccess: (_data, { tenantId, code }) => {
      toast({
        title: "App deleted successfully",
        description: "The app has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey(tenantId) });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(tenantId, code) });
    },
    onError: () => {
      toast({
        title: "Failed to delete app",
        description: "The app has not been deleted successfully",
        variant: "destructive",
      });
    }
  })
}
