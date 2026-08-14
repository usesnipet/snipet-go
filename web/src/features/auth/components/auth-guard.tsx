import { useMeUser } from "@/features/user/hooks";
import { ROUTES } from "@/routes";
import { Navigate, Outlet, useLocation } from "react-router";

type AuthGuardProps = {
  // "private" (default) requires the user to be logged in, redirecting to login otherwise.
  // "public" requires the user to be logged out, redirecting to a logged-in route otherwise.
  mode?: "private" | "public";
}

export const AuthGuard = ({ mode = "private" }: AuthGuardProps) => {
  const { pathname } = useLocation();
  const { isLoading, error } = useMeUser();
  const isAuthenticated = !error;

  if (isLoading) return null;

  if (mode === "public") {
    if (isAuthenticated) return <Navigate to={ROUTES.selectTenant} replace />;
    return <Outlet />
  }

  if (!isAuthenticated) {
    const redirectURL = `${ROUTES.authLogin}?redirect=${pathname}`;
    return <Navigate to={redirectURL} replace />;
  }
  return <Outlet />
}
