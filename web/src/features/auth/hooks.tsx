import { ToastAction } from "@/components/ui/toast";
import { useNavigate } from "@/hooks/use-navigate";
import { toast } from "@/hooks/use-toast";
import { ApiError } from "@/lib/http";
import { logger } from "@/lib/logger";
import { queryClient } from "@/lib/query-client";
import { ROUTES } from "@/routes";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect } from "react";

import { meUserQueryKey } from "../user/hooks";

import { authService } from "./service";
import { useAuthStore } from "./store";

import type {
  ActivateAccount,
  AuthenticateResponse,
  AuthProvider,
  AuthorizationUrlResponse,
  ForgotPassword,
  Login,
  ProviderCallbackQuery,
  Refresh,
  Register,
  RegisterResponse,
  ResendActivation,
  ResetPassword,
  SetPassword,
} from "./schemas";
import type {
  ServiceGetOptions,
  ServicePostOptions,
  ServicePutOptions,
} from "@/lib/services";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";
import type { RoutePath } from "@/routes";
const BASE_QUERY_KEY = "auth";

export const registerQueryKey = () => [BASE_QUERY_KEY, "register"] as const;
export const useRegister = (
  opts?: ServicePostOptions<Register, RegisterResponse>,
): UseMutationResult<RegisterResponse, Error, Register> => {
  const navigate = useNavigate();

  return useMutation({
    mutationKey: registerQueryKey(),
    mutationFn: (data) => authService.register(data, opts),
    onSuccess: () => {
      toast({
        title: "Account created",
        description: "Check your email to activate your account",
      });
      navigate(ROUTES.authLogin);
    },
    onError: (e) => {
      logger.error("Failed to register", { error: e });
      toast({
        title: "Registration failed",
        description: e.message ?? "Could not create your account",
        variant: "destructive",
      });
    },
  });
};

export const loginQueryKey = () => [BASE_QUERY_KEY, "login"] as const;
export const useLogin = (
  redirect?: RoutePath,
  opts?: ServicePostOptions<Login, AuthenticateResponse>,
): UseMutationResult<AuthenticateResponse, Error, Login> => {
  const navigate = useNavigate();
  const setTokens = useAuthStore((state) => state.setTokens);
  const { mutate: resendActivation } = useResendActivation();

  return useMutation({
    mutationKey: loginQueryKey(),
    mutationFn: (data) => authService.login(data, opts),
    onSuccess: async (data) => {
      setTokens(data);
      toast({
        title: "Signed in successfully",
        description: "You have been authenticated successfully",
      });
      await queryClient.invalidateQueries({ queryKey: meUserQueryKey() });
      await navigate(redirect ?? ROUTES.tenantHome, { replace: true });
    },
    onError: (e, variables) => {
      logger.error("Failed to login", { error: e });

      if (ApiError.is(e) && e.statusCode === 403) {
        toast({
          title: "Account not activated",
          description: "Check your email to activate your account.",
          variant: "destructive",
          action: (
            <ToastAction
              altText="Resend activation email"
              onClick={() => resendActivation({ email: variables.email })}
            >
              Resend email
            </ToastAction>
          ),
        });
        return;
      }

      toast({
        title: "Failed to sign in",
        description: "Invalid credentials",
        variant: "destructive",
      });
    },
  });
};

export type AuthorizationUrlVariables = {
  provider: AuthProvider;
};

export const authorizationUrlQueryKey = () =>
  [BASE_QUERY_KEY, "authorizationUrl"] as const;
export const useAuthorizationUrl = (
  opts?: ServiceGetOptions<AuthorizationUrlResponse>,
): UseMutationResult<AuthorizationUrlResponse, Error, AuthorizationUrlVariables> => {
  return useMutation({
    mutationKey: authorizationUrlQueryKey(),
    mutationFn: ({ provider }) =>
      authService.getAuthorizationUrl(provider, opts),
    onError: () => {
      toast({
        title: "OAuth failed",
        description: "Could not start provider login",
        variant: "destructive",
      });
    },
  });
};

export type OAuthCallbackVariables = {
  provider: AuthProvider;
  searchParams: ProviderCallbackQuery;
};

export const oauthCallbackQueryKey = () =>
  [BASE_QUERY_KEY, "oauthCallback"] as const;
export const useOAuthCallback = (
  opts?: ServiceGetOptions<AuthenticateResponse, ProviderCallbackQuery>,
): UseMutationResult<AuthenticateResponse, Error, OAuthCallbackVariables> => {
  return useMutation({
    mutationKey: oauthCallbackQueryKey(),
    mutationFn: ({ provider, searchParams }) =>
      authService.callback(provider, searchParams, opts),
    onSuccess: () => {
      toast({
        title: "Signed in successfully",
        description: "You have been authenticated successfully",
      });
    },
    onError: () => {
      toast({
        title: "Failed to sign in",
        description: "OAuth callback was not successful",
        variant: "destructive",
      });
    },
  });
};

export const refreshQueryKey = () => [BASE_QUERY_KEY, "refresh"] as const;
export const useRefresh = (
  opts?: ServicePostOptions<Refresh, AuthenticateResponse>,
): UseMutationResult<AuthenticateResponse, Error, Refresh> => {
  return useMutation({
    mutationKey: refreshQueryKey(),
    mutationFn: (data) =>
      authService.refresh(data, opts),
    onError: () => {
      toast({
        title: "Session expired",
        description: "Please sign in again",
        variant: "destructive",
      });
    },
  });
};

export const logoutQueryKey = () => [BASE_QUERY_KEY, "logout"] as const;
export const useLogout = (
  opts?: ServicePostOptions<Refresh, void>,
): UseMutationResult<void, Error> => {
  const clearAuth = useAuthStore((state) => state.clear);
  const refreshToken = useAuthStore((state) => state.refreshToken);
  const navigate = useNavigate();

  return useMutation({
    mutationKey: logoutQueryKey(),
    mutationFn: async () => {
      if (!refreshToken) return;
      return authService.logout({ refresh_token: refreshToken}, opts);
    },
    onSuccess: () => {
      toast({
        title: "Signed out",
        description: "You have been signed out successfully",
      });
      clearAuth()
      navigate(ROUTES.authLogin, { replace: true })
    },
    onError: () => {
      toast({
        title: "Failed to sign out",
        description: "Could not revoke your session",
        variant: "destructive",
      });
    },
  });
};

export const setPasswordQueryKey = () => [BASE_QUERY_KEY, "setPassword"] as const;
export const useSetPassword = (
  opts?: ServicePutOptions<SetPassword, void>,
): UseMutationResult<void, Error, SetPassword> => {
  return useMutation({
    mutationKey: setPasswordQueryKey(),
    mutationFn: (data) =>
      authService.setPassword(data, opts),
    onSuccess: () => {
      toast({
        title: "Password updated",
        description: "Your password has been set successfully",
      });
    },
    onError: () => {
      toast({
        title: "Failed to set password",
        description: "Could not update your password",
        variant: "destructive",
      });
    },
  });
};

export const forgotPasswordQueryKey = () =>
  [BASE_QUERY_KEY, "forgotPassword"] as const;
export const useForgotPassword = (
  opts?: ServicePostOptions<ForgotPassword, void>,
): UseMutationResult<void, Error, ForgotPassword> => {
  return useMutation({
    mutationKey: forgotPasswordQueryKey(),
    mutationFn: (data) =>
      authService.forgotPassword(data, opts),
    onSuccess: () => {
      toast({
        title: "Check your email",
        description: "If an account exists, a reset link has been sent",
      });
    },
    onError: () => {
      toast({
        title: "Request failed",
        description: "Could not process password reset request",
        variant: "destructive",
      });
    },
  });
};

export const resetPasswordQueryKey = () =>
  [BASE_QUERY_KEY, "resetPassword"] as const;
export const useResetPassword = (
  opts?: ServicePostOptions<ResetPassword, void>,
): UseMutationResult<void, Error, ResetPassword> => {
  return useMutation({
    mutationKey: resetPasswordQueryKey(),
    mutationFn: (data) =>
      authService.resetPassword(data, opts),
    onSuccess: () => {
      toast({
        title: "Password reset",
        description: "You can now sign in with your new password",
      });
    },
    onError: () => {
      toast({
        title: "Reset failed",
        description: "Invalid or expired reset token",
        variant: "destructive",
      });
    },
  });
};

export const activateQueryKey = (token: string) =>
  [BASE_QUERY_KEY, "activate", token] as const;
export const useActivate = (
  token: string,
  opts?: ServicePostOptions<ActivateAccount, void>,
): UseQueryResult<null, Error> => {
  const navigate = useNavigate();
  const query = useQuery({
    queryKey: activateQueryKey(token),
    queryFn: async () => {
      await authService.activate({ token }, opts);
      return null;
    },
    enabled: Boolean(token),
    retry: false,
    staleTime: Infinity,
  });

  useEffect(() => {
    if (query.isSuccess) {
      toast({
        title: "Account activated",
        description: "You can now sign in",
      });
      navigate(ROUTES.authLogin);
    }
    if (query.isError) {
      toast({
        title: "Activation failed",
        description: "Invalid or expired activation token",
        variant: "destructive",
      });
    }
  }, [query.isSuccess, query.isError, navigate]);

  return query;
};

export const resendActivationQueryKey = () =>
  [BASE_QUERY_KEY, "resendActivation"] as const;
export const useResendActivation = (
  opts?: ServicePostOptions<ResendActivation, void>,
): UseMutationResult<void, Error, ResendActivation> => {
  return useMutation({
    mutationKey: resendActivationQueryKey(),
    mutationFn: (data) =>
      authService.resendActivation(data, opts),
    onSuccess: () => {
      toast({
        title: "Check your email",
        description: "If an account exists, an activation link has been sent",
      });
    },
    onError: () => {
      toast({
        title: "Request failed",
        description: "Could not resend activation email",
        variant: "destructive",
      });
    },
  });
};
