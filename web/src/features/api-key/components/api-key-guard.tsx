import { ROUTES } from "@/routes";
import { Navigate, Outlet, useLocation } from "react-router";

import { useMeApiKey } from "../hooks";
import { useApiKeyStore } from "../store";

export const ApiKeyGuard = () => {
  const { pathname } = useLocation();

  const key = useApiKeyStore((state) => state.key);
  const { isLoading, error } = useMeApiKey();
  const redirectURL = `${ROUTES.authApiKey}?redirect=${pathname}`;
  if (!key) return <Navigate to={redirectURL} />
  if (isLoading) return null;
  if (error) return <Navigate to={redirectURL} />
  return <Outlet />
}
