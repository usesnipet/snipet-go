import type { ROUTES } from "@/routes";
import { useNavigate as useNavigateReactRouter } from "react-router";

import type { NavigateOptions } from "react-router";

export const useNavigate = () => {
  const navigate = useNavigateReactRouter();
  return (path: (typeof ROUTES)[keyof typeof ROUTES], options?: NavigateOptions) => navigate(path, options);
}