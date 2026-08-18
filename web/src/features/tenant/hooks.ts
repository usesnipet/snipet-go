import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useParams } from "react-router";

import { tenantService } from "./service";
import { useTenantStore } from "./store";

import type { CreateTenant, PaginatedTenant, Tenant, UpdateTenant } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";
const BASE_QUERY_KEY = "tenant";

export const findMineTenantQueryKey = () => [BASE_QUERY_KEY, "findMine"] as const;
export const useFindMineTenant = (
  opts?: ServiceGetOptions<PaginatedTenant>
): UseQueryResult<PaginatedTenant, Error> => {
  return useQuery({
    queryKey: findMineTenantQueryKey(),
    queryFn: () => tenantService.findMine(opts),
  });
};

export const findByIDTenantQueryKey = (id: string) => [BASE_QUERY_KEY, "findByID", id] as const;
export const useFindByIDTenant = (
  id: string,
  opts?: ServiceGetOptions<Tenant>
): UseQueryResult<Tenant, Error> => {
  return useQuery({
    queryKey: findByIDTenantQueryKey(id),
    queryFn: () => tenantService.findByID(id, opts),
    enabled: !!id,
  });
};

export const findBySlugTenantQueryKey = (slug: string) => [BASE_QUERY_KEY, "findBySlug", slug] as const;
export const useFindBySlugTenant = (
  slug: string,
  opts?: ServiceGetOptions<Tenant>
): UseQueryResult<Tenant, Error> => {
  return useQuery({
    queryKey: findBySlugTenantQueryKey(slug),
    queryFn: () => tenantService.findBySlug(slug, opts),
    enabled: !!slug,
  });
};

export const createTenantQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateTenant = (
  opts?: ServicePostOptions<CreateTenant, Tenant>
): UseMutationResult<Tenant, Error, CreateTenant> => {
  return useMutation({
    mutationKey: createTenantQueryKey(),
    mutationFn: (data: CreateTenant) => tenantService.create(data, opts),
    onSuccess: () => {
      toast({
        title: "Tenant created successfully",
        description: "The tenant has been created successfully",
      });
      queryClient.invalidateQueries({ queryKey: findMineTenantQueryKey() });
    },
    onError: () => {
      toast({
        title: "Failed to create tenant",
        description: "The tenant has not been created successfully",
        variant: "destructive",
      });
    },
  });
};

export const updateTenantQueryKey = () => [BASE_QUERY_KEY, "update"];
export const useUpdateTenant = (
  opts?: ServicePutOptions<UpdateTenant, void>
): UseMutationResult<void, Error, { id: string; data: UpdateTenant }> => {
  return useMutation({
    mutationKey: updateTenantQueryKey(),
    mutationFn: ({ id, data }: { id: string; data: UpdateTenant }) =>
      tenantService.update(id, data, opts),
    onSuccess: (_data, { id }) => {
      toast({
        title: "Tenant updated successfully",
        description: "The tenant has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: findMineTenantQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIDTenantQueryKey(id) });
    },
    onError: () => {
      toast({
        title: "Failed to update tenant",
        description: "The tenant has not been updated successfully",
        variant: "destructive",
      });
    },
  });
};

export const deleteTenantQueryKey = () => [BASE_QUERY_KEY, "delete"];
export const useDeleteTenant = (
  opts?: ServiceDeleteOptions<void>
): UseMutationResult<void, Error, string> => {
  return useMutation({
    mutationKey: deleteTenantQueryKey(),
    mutationFn: (id: string) => tenantService.delete(id, opts),
    onSuccess: (_data, id) => {
      toast({
        title: "Tenant deleted successfully",
        description: "The tenant has been deleted successfully",
      });
      queryClient.invalidateQueries({ queryKey: findMineTenantQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIDTenantQueryKey(id) });
    },
    onError: () => {
      toast({
        title: "Failed to delete tenant",
        description: "The tenant has not been deleted successfully",
        variant: "destructive",
      });
    },
  });
};

export const leaveTenantQueryKey = () => [BASE_QUERY_KEY, "leave"];
export const useLeaveTenant = (
  opts?: ServicePostOptions<void, void>
): UseMutationResult<void, Error, string> => {
  return useMutation({
    mutationKey: leaveTenantQueryKey(),
    mutationFn: (id: string) => tenantService.leave(id, opts),
    onSuccess: (_data, id) => {
      toast({
        title: "Left tenant",
        description: "You have left the tenant successfully",
      });
      queryClient.invalidateQueries({ queryKey: findMineTenantQueryKey() });
      queryClient.invalidateQueries({ queryKey: findByIDTenantQueryKey(id) });
    },
    onError: () => {
      toast({
        title: "Failed to leave tenant",
        description: "Could not leave the tenant",
        variant: "destructive",
      });
    },
  });
};

export const useCurrentTenant = (): UseQueryResult<Tenant, Error> => {
  const { tenantSlug } = useParams<{ tenantSlug: string }>();
  if (!tenantSlug) {
    throw new Error("Tenant slug is required");
  }

  const query = useFindBySlugTenant(tenantSlug);
  const setTenant = useTenantStore((state) => state.setTenant);
  const clearTenant = useTenantStore((state) => state.clearTenant);

  useEffect(() => {
    if (query.data) setTenant(query.data);
  }, [query.data, setTenant]);

  useEffect(() => {
    if (query.isError) clearTenant();
  }, [query.isError, clearTenant]);

  return query;
}