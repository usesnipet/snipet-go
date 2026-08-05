import { LoadingFallback } from "@/components/loading-fallback";
import { useAuth, useAuthenticateAnonymous } from "@/features/auth/hooks";
import { ChatAuthGate } from "@/features/chat/components/chat-auth-gate";
import { useFindPublicByCodeClient } from "@/features/client/hooks";
import { applyPathParams } from "@/lib/http";
import { ROUTES } from "@/routes";
import { useEffect } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router";

let authInFlight: Promise<unknown> | null = null;

function isExpired(expiresAt: Date | null): boolean {
  return !expiresAt || expiresAt < new Date();
}

export function LoginAnonymous() {
  const { clientCode = "" } = useParams<{ clientCode: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { accessToken, accessTokenExpiresAt } = useAuth();
  const { data: client, isLoading: isLoadingClient } = useFindPublicByCodeClient(clientCode);
  const { mutateAsync: authenticateAnonymous } = useAuthenticateAnonymous();

  const accessValid = !!accessToken && !isExpired(accessTokenExpiresAt);
  const redirectTo = searchParams.get("redirect") ?? applyPathParams(ROUTES.clientChat, { clientCode });

  useEffect(() => {
    if (isLoadingClient || !clientCode || !client) return;
    if (!client.allow_anonymous) return;

    if (accessValid) {
      navigate(redirectTo, { replace: true });
      return;
    }

    if (authInFlight) {
      void authInFlight
        .then(() => navigate(redirectTo, { replace: true }))
        .catch(() => {
          navigate(applyPathParams(ROUTES.clientChat, { clientCode }), { replace: true });
        });
      return;
    }

    authInFlight = authenticateAnonymous({ clientCode, data: {} }).finally(() => {
      authInFlight = null;
    });

    void authInFlight
      .then(() => navigate(redirectTo, { replace: true }))
      .catch(() => {
        navigate(applyPathParams(ROUTES.clientChat, { clientCode }), { replace: true });
      });
  }, [
    accessValid,
    authenticateAnonymous,
    client,
    clientCode,
    isLoadingClient,
    navigate,
    redirectTo,
  ]);

  if (isLoadingClient) return <LoadingFallback />;

  if (client && !client.allow_anonymous) {
    return (
      <ChatAuthGate
        clientCode={clientCode}
        allowAnonymous={false}
        redirectPath={redirectTo}
      />
    );
  }

  return <LoadingFallback />;
}
