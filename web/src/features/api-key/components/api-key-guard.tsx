import { ROUTES } from "@/routes";
import { Navigate, Outlet } from "react-router";

import { useMeApiKey } from "../hooks";
import { useApiKeyStore } from "../store";

export const ApiKeyGuard = () => {
  const key = useApiKeyStore((state) => state.key);
  const { isLoading, error } = useMeApiKey();
  if (!key) return <Navigate to={ROUTES.authApiKey} />
  if (isLoading) return null;
  if (error) return <Navigate to={ROUTES.authApiKey} />
  return <Outlet />
}
