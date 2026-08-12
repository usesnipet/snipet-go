import { toast } from "@/hooks/use-toast";
import { queryClient } from "@/lib/query-client";
import { useMutation, useQuery } from "@tanstack/react-query";

import { userService } from "./service";

import type { UpdateProfilePicture, User } from "./schemas";
import type { ServiceGetOptions, ServicePutOptions } from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "user";

export const meUserQueryKey = () => [BASE_QUERY_KEY, "me"] as const;
export const useMeUser = (
  opts?: ServiceGetOptions<User>,
): UseQueryResult<User, Error> => {
  return useQuery({
    queryKey: meUserQueryKey(),
    queryFn: () => userService.me(opts),
  });
};

export const updatePictureUserQueryKey = () =>
  [BASE_QUERY_KEY, "updatePicture"] as const;
export const useUpdatePictureUser = (
  opts?: ServicePutOptions<UpdateProfilePicture, void>,
): UseMutationResult<void, Error, UpdateProfilePicture> => {
  return useMutation({
    mutationKey: updatePictureUserQueryKey(),
    mutationFn: (data) => userService.updatePicture(data, opts),
    onSuccess: () => {
      toast({
        title: "Picture updated",
        description: "Your profile picture has been updated successfully",
      });
      queryClient.invalidateQueries({ queryKey: meUserQueryKey() });
    },
    onError: () => {
      toast({
        title: "Failed to update picture",
        description: "Your profile picture could not be updated",
        variant: "destructive",
      });
    },
  });
};
