import { toast } from "@/hooks/use-toast";
import { useMutation } from "@tanstack/react-query";
import { useShallow } from "zustand/react/shallow";

import { authService } from "./service";
import { useAuthStore } from "./store";

import type {
  AuthenticateAnonymous, AuthenticateResponse, AuthProvider, AuthTokens, Refresh
} from "./schemas";
import type { ServicePostOptions } from "@/lib/services";
import type { UseMutationResult } from "@tanstack/react-query";
const BASE_QUERY_KEY = "auth";

function toAuthTokens(data: AuthenticateResponse): AuthTokens {
  return {
    access_token: data.access_token,
    access_token_expires_at: data.access_token_expires_at,
    refresh_token: data.refresh_token,
    refresh_token_expires_at: data.refresh_token_expires_at,
  };
}

export type AuthenticateVariables = {
  clientCode: string;
  provider: AuthProvider;
  headers?: Record<string, string>;
  body?: unknown;
};

export const authenticateQueryKey = () => [BASE_QUERY_KEY, "authenticate"] as const;
export const useAuthenticate = (
  opts?: ServicePostOptions<unknown, AuthenticateResponse>,
): UseMutationResult<AuthenticateResponse, Error, AuthenticateVariables> => {
  const setTokens = useAuthStore((state) => state.setTokens);

  return useMutation({
    mutationKey: authenticateQueryKey(),
    mutationFn: ({ clientCode, provider, headers, body }) =>
      authService.authenticate(clientCode, provider, {
        ...opts,
        auth: false,
        headers: { ...opts?.headers, ...headers },
        body: body ?? opts?.body,
      }),
    onSuccess: (data) => {
      setTokens(toAuthTokens(data));
      toast({
        title: "Signed in successfully",
        description: "You have been authenticated successfully",
      });
    },
    onError: () => {
      toast({
        title: "Failed to sign in",
        description: "Authentication was not successful",
        variant: "destructive",
      });
    },
  })
}

export type AuthenticateAnonymousVariables = {
  clientCode: string;
  data?: AuthenticateAnonymous;
};

export const authenticateAnonymousQueryKey = () =>
  [BASE_QUERY_KEY, "authenticateAnonymous"] as const;
export const useAuthenticateAnonymous = (
  opts?: ServicePostOptions<AuthenticateAnonymous, AuthenticateResponse>,
): UseMutationResult<AuthenticateResponse, Error, AuthenticateAnonymousVariables> => {
  const setTokens = useAuthStore((state) => state.setTokens);

  return useMutation({
    mutationKey: authenticateAnonymousQueryKey(),
    mutationFn: ({ clientCode, data = {} }) =>
      authService.authenticateAnonymous(clientCode, data, { ...opts, auth: false }),
    onSuccess: (response) => {
      setTokens(toAuthTokens(response));
      toast({
        title: "Signed in successfully",
        description: "You have been authenticated anonymously",
      });
    },
    onError: () => {
      toast({
        title: "Failed to sign in",
        description: "Anonymous authentication was not successful",
        variant: "destructive",
      });
    },
  })
}

export type RefreshVariables = {
  clientCode: string;
  data?: Refresh;
};

export const refreshQueryKey = () => [BASE_QUERY_KEY, "refresh"] as const;
export const useRefresh = (
  opts?: ServicePostOptions<Refresh, AuthenticateResponse>,
): UseMutationResult<AuthenticateResponse, Error, RefreshVariables> => {
  const setTokens = useAuthStore((state) => state.setTokens);
  const clear = useAuthStore((state) => state.clear);

  return useMutation({
    mutationKey: refreshQueryKey(),
    mutationFn: ({ clientCode, data }) => {
      const refresh_token = data?.refresh_token ?? useAuthStore.getState().refreshToken;
      if (!refresh_token) {
        throw new Error("No refresh token available");
      }
      return authService.refresh(
        clientCode,
        { refresh_token },
        { ...opts, auth: false },
      );
    },
    onSuccess: (response) => {
      setTokens(toAuthTokens(response));
    },
    onError: () => {
      clear();
      toast({
        title: "Session expired",
        description: "Please sign in again",
        variant: "destructive",
      });
    },
  })
}

export const useAuth = () => {
  return useAuthStore(
    useShallow((s) => ({
      accessToken: s.accessToken,
      accessTokenExpiresAt: s.accessTokenExpiresAt,
      refreshToken: s.refreshToken,
      refreshTokenExpiresAt: s.refreshTokenExpiresAt,
    }))
  )
}