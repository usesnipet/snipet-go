import { useMeUser } from "@/features/user/hooks";
import { ROUTES } from "@/routes";
import { Navigate, Outlet, useLocation } from "react-router";

export const AuthGuard = () => {
  const { pathname } = useLocation();
  const { isLoading, error } = useMeUser();
  const redirectURL = `${ROUTES.authLogin}?redirect=${pathname}`;

  if (isLoading) return null;
  if (error) return <Navigate to={redirectURL} />
  return <Outlet />
}
