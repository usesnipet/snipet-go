import { useMutation, useQuery } from "@tanstack/react-query";

import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";

import { invitationService } from "./service";

import type {
  AcceptInvitation, CreateInvitation, DeclineInvitation, Invitation, InvitationInfo,
  ListInvitationSearchParams, PaginatedInvitation
} from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions
} from "@/lib/services";
import type { Member } from "@/models/member";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "invitation";

export const filterInvitationQueryKey = (tenantId: string) =>
  [BASE_QUERY_KEY, "filter", tenantId] as const;
export const useFilterInvitation = (
  tenantId: string,
  opts?: ServiceGetOptions<PaginatedInvitation, ListInvitationSearchParams>,
): UseQueryResult<PaginatedInvitation, Error> => {
  return useQuery({
    queryKey: [...filterInvitationQueryKey(tenantId), opts?.searchParams],
    queryFn: () => invitationService.filter(tenantId, opts),
    enabled: !!tenantId,
  })
}

export const createInvitationQueryKey = () => [BASE_QUERY_KEY, "create"];
export const useCreateInvitation = (
  opts?: ServicePostOptions<CreateInvitation, Invitation>,
): UseMutationResult<Invitation, Error, { tenantId: string; data: CreateInvitation }> => {
  return useMutation({
    mutationKey: createInvitationQueryKey(),
    mutationFn: ({ tenantId, data }) =>
      invitationService.create(tenantId, data, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Invitation sent successfully",
        description: "The invitation has been sent successfully",
      });
      queryClient.invalidateQueries({ queryKey: filterInvitationQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to send invitation",
        description: "The invitation has not been sent successfully",
        variant: "destructive",
      });
    }
  })
}

export const removeInvitationQueryKey = () => [BASE_QUERY_KEY, "remove"];
export const useRemoveInvitation = (
  opts?: ServiceDeleteOptions<void>,
): UseMutationResult<void, Error, { tenantId: string; id: string }> => {
  return useMutation({
    mutationKey: removeInvitationQueryKey(),
    mutationFn: ({ tenantId, id }) =>
      invitationService.remove(tenantId, id, opts),
    onSuccess: (_data, { tenantId }) => {
      toast({
        title: "Invitation canceled successfully",
        description: "The invitation has been canceled successfully",
      });
      queryClient.invalidateQueries({ queryKey: filterInvitationQueryKey(tenantId) });
    },
    onError: () => {
      toast({
        title: "Failed to cancel invitation",
        description: "The invitation has not been canceled successfully",
        variant: "destructive",
      });
    }
  })
}

export const getByTokenInvitationQueryKey = (token: string) =>
  [BASE_QUERY_KEY, "getByToken", token] as const;
export const useGetInvitationByToken = (
  token: string,
  opts?: ServiceGetOptions<InvitationInfo>,
): UseQueryResult<InvitationInfo, Error> => {
  return useQuery({
    queryKey: getByTokenInvitationQueryKey(token),
    queryFn: () => invitationService.getByToken(token, opts),
    enabled: !!token,
    retry: false,
  })
}

export const acceptInvitationQueryKey = () => [BASE_QUERY_KEY, "accept"];
export const useAcceptInvitation = (
  opts?: ServicePostOptions<AcceptInvitation, Member>,
): UseMutationResult<Member, Error, AcceptInvitation> => {
  return useMutation({
    mutationKey: acceptInvitationQueryKey(),
    mutationFn: (data: AcceptInvitation) =>
      invitationService.accept(data, opts),
    onSuccess: () => {
      toast({
        title: "Invitation accepted successfully",
        description: "You have joined the tenant successfully",
      });
    },
    onError: () => {
      toast({
        title: "Failed to accept invitation",
        description: "The invitation has not been accepted successfully",
        variant: "destructive",
      });
    }
  })
}

export const declineInvitationQueryKey = () => [BASE_QUERY_KEY, "decline"];
export const useDeclineInvitation = (
  opts?: ServicePostOptions<DeclineInvitation, void>,
): UseMutationResult<void, Error, DeclineInvitation> => {
  return useMutation({
    mutationKey: declineInvitationQueryKey(),
    mutationFn: (data: DeclineInvitation) =>
      invitationService.decline(data, opts),
    onSuccess: () => {
      toast({
        title: "Invitation declined",
        description: "You have declined the invitation",
      });
    },
    onError: () => {
      toast({
        title: "Failed to decline invitation",
        description: "The invitation has not been declined successfully",
        variant: "destructive",
      });
    }
  })
}
