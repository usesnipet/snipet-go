import { useMutation, useQuery } from "@tanstack/react-query";

import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";

import { memberService } from "./service";

import type { CreateMember, ListMemberSearchParams, PaginatedMember, UpdateMemberRole } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";
import type { Member } from "@/models/member";

const BASE_QUERY_KEY = "member";

export const filterMemberQueryKey = (tenantId: string) =>
  [BASE_QUERY_KEY, "filter", tenantId] as const;
export const useFilterMember = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedMember, ListMemberSearchParams>,
): UseQueryResult<PaginatedMember, Error> => {
  return useQuery({
    queryKey: [...filterMemberQueryKey(tenantId), opts?.searchParams],
    queryFn: () => memberService.filter(tenantId, opts),
    enabled: !!tenantId,
  })
}

export const createMemberQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateMember = (
  opts?: ServicePostOptions<CreateMember, Member>,
): UseMutationResult<Member, Error, { tenantId: string; data: CreateMember }> => {
  return useMutation({
    mutationKey: createMemberQueryKey(),
    mutationFn: ({ tenantId, data }) =>
      memberService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Member created successfully",
        description: "The member has been created and can now sign in",
      });
      queryClient.invalidateQueries({ queryKey: filterMemberQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to create member",
        description: "The member has not been created successfully",
        variant: "destructive",
      });
    }
  })
}

export const updateRoleMemberQueryKey = () => [BASE_QUERY_KEY, "updateRole"];
export const useUpdateRoleMember = (
  opts?: ServicePutOptions<UpdateMemberRole, void>,
): UseMutationResult<void, Error, { tenantId: string; id: string; data: UpdateMemberRole }> => {
  return useMutation({
    mutationKey: updateRoleMemberQueryKey(),
    mutationFn: ({ tenantId, id, data }) =>
      memberService.updateRole(tenantId, id, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Member role updated successfully",
        description: "The member role has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: filterMemberQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to update member role",
        description: "The member role has not been updated successfully",
        variant: "destructive",
      });
    }
  })
}

export const removeMemberQueryKey = () => [BASE_QUERY_KEY, "remove"];
export const useRemoveMember = (
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: removeMemberQueryKey(),
    mutationFn: ({ tenantId, id }) =>
      memberService.remove(tenantId, id, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Member removed successfully",
        description: "The member has been removed from the tenant successfully",
      });
      queryClient.invalidateQueries({ queryKey: filterMemberQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to remove member",
        description: "The member has not been removed successfully",
        variant: "destructive",
      });
    }
  })
}
