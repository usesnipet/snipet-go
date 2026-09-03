import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { appService } from "./service";

import type {
  App, AppWithSecret, CreateApp, LinkAppAgents, PaginatedApp, UpdateApp, UpdateAppAuthConfig
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "app";

export const listAppQueryKey = () => [BASE_QUERY_KEY, "list"] as const;
export const useListApp = (
  opts?: ServiceGetOptions<PaginatedApp>
): UseQueryResult<PaginatedApp, Error> => {
  return useQuery({
    queryKey: listAppQueryKey(),
    queryFn: () => appService.list(opts),
  })
}

export const findByCodeAppQueryKey = (code: string) =>
  [BASE_QUERY_KEY, "findByCode", code];
export const useFindByCodeApp = (
  code: string,
  opts?: ServiceGetOptions<App>
): UseQueryResult<App, Error> => {
  return useQuery({
    queryKey: findByCodeAppQueryKey(code),
    queryFn: (): Promise<App> =>
      appService.findByCode(code, opts),
    enabled: !!code,
  })
}

export const createAppQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateApp = (
  opts?: ServicePostOptions<CreateApp, AppWithSecret>
): UseMutationResult<AppWithSecret, Error, { data: CreateApp }> => {
  return useMutation({
    mutationKey: createAppQueryKey(),
    mutationFn: ({ data }: { data: CreateApp }) =>
      appService.create(data, opts),
    onSuccess: () => {
      toast({
        title: "App created successfully",
        description: "The app has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
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
): UseMutationResult<void, Error, { code: string; data: UpdateApp }> => {
  return useMutation({
    mutationKey: updateAppQueryKey(),
    mutationFn: ({ code, data }: { code: string; data: UpdateApp }) =>
      appService.update(code, data, opts),
    onSuccess: (_data, { code }) => {
      toast({
        title: "App updated successfully",
        description: "The app has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(code) });
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
): UseMutationResult<void, Error, { code: string; data: UpdateAppAuthConfig }> => {
  return useMutation({
    mutationKey: updateAppAuthConfigQueryKey(),
    mutationFn: ({ code, data }: { code: string; data: UpdateAppAuthConfig }) =>
      appService.updateAuthConfig(code, data, opts),
    onSuccess: (_data, { code }) => {
      toast({
        title: "App updated successfully",
        description: "The app has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(code) });
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

export const linkAppAgentsQueryKey = () => [BASE_QUERY_KEY, "link-agents"];
export const useLinkAppAgents = (
  opts: ServicePutOptions<LinkAppAgents, void> = {}
): UseMutationResult<void, Error, { code: string; data: LinkAppAgents }> => {
  return useMutation({
    mutationKey: linkAppAgentsQueryKey(),
    mutationFn: ({ code, data }: { code: string; data: LinkAppAgents }) =>
      appService.linkAgents(code, data, opts),
    onSuccess: (_data, { code }) => {
      toast({
        title: "App agents updated successfully",
        description: "The app's linked agents have been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(code) });
    },
    onError: () => {
      toast({
        title: "Failed to update app agents",
        description: "The app's linked agents have not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const rollAppQueryKey = () => [BASE_QUERY_KEY, "roll"];
export const useRollApp = (
  opts?: ServicePostOptions<undefined, AppWithSecret>
): UseMutationResult<AppWithSecret, Error, { code: string }> => {
  return useMutation({
    mutationKey: rollAppQueryKey(),
    mutationFn: ({ code }: { code: string }) =>
      appService.roll(code, opts),
    onSuccess: (_data, { code }) => {
      toast({
        title: "App key rolled successfully",
        description: "The app key has been rolled successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(code) });
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
): UseMutationResult<void, Error, { code: string; active: boolean }> => {
  return useMutation({
    mutationKey: setActiveAppQueryKey(),
    mutationFn: ({ code, active }: { code: string; active: boolean }) =>
      appService.setActive(code, active, opts),
    onSuccess: (_data, { code, active }) => {
      toast({
        title: active ? "App activated" : "App deactivated",
        description: active
          ? "The app has been activated successfully"
          : "The app has been deactivated successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(code) });
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
): UseMutationResult<void, Error, { code: string }> => {
  return useMutation({
    mutationKey: deleteAppQueryKey(),
    mutationFn: ({ code }: { code: string }) =>
      appService.delete(code, opts),
    onSuccess: (_data, { code }) => {
      toast({
        title: "App deleted successfully",
        description: "The app has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: listAppQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByCodeAppQueryKey(code) });
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
