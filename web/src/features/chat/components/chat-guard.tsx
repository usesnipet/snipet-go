import { LoadingFallback } from "@/components/loading-fallback";
import { useAuth, useRefresh } from "@/features/auth/hooks";
import { useFindPublicByCodeClient } from "@/features/client/hooks";
import { useEffect } from "react";
import { Outlet, useLocation, useParams } from "react-router";

import { ChatAuthGate } from "./chat-auth-gate";

let refreshInFlight: Promise<unknown> | null = null;

function isExpired(expiresAt: Date | null): boolean {
  return !expiresAt || expiresAt < new Date();
}

export function ChatGuard() {
  const { accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt } = useAuth();
  const { clientCode: clientCodeParam } = useParams<{ clientCode: string }>();
  const clientCode = clientCodeParam ?? "";
  const { pathname } = useLocation();
  const { data: client, isLoading: isLoadingClient } = useFindPublicByCodeClient(clientCode);
  const { mutateAsync: refresh } = useRefresh();

  const accessValid = !!accessToken && !isExpired(accessTokenExpiresAt);
  const refreshValid = !!refreshToken && !isExpired(refreshTokenExpiresAt);
  const needsRefresh = !accessValid && refreshValid;

  useEffect(() => {
    if (isLoadingClient || !clientCode || !client || accessValid || !refreshValid) return;
    if (refreshInFlight || !refreshToken) return;

    refreshInFlight = refresh({
      clientCode,
      data: { refresh_token: refreshToken },
    }).finally(() => {
      refreshInFlight = null;
    });
  }, [
    accessValid,
    client,
    clientCode,
    isLoadingClient,
    refresh,
    refreshToken,
    refreshValid,
  ]);

  if (isLoadingClient || needsRefresh) {
    return <LoadingFallback />;
  }

  if (!accessValid) {
    return (
      <ChatAuthGate
        clientCode={clientCode}
        allowAnonymous={!!client?.allow_anonymous}
        redirectPath={pathname}
      />
    );
  }

  return <Outlet />;
}
