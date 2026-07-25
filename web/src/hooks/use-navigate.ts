import { applyPathParams } from "@/lib/http";
import { useNavigate as useNavigateReactRouter } from "react-router";

import type { ROUTES } from "@/routes";
import type { NavigateOptions } from "react-router";

type Options = NavigateOptions & {
  params?: Record<string, string>;
}

export const useNavigate = () => {
  const navigate = useNavigateReactRouter();
  return (path: (typeof ROUTES)[keyof typeof ROUTES], options?: Options) => {
    const {params, ...rest} = options ?? {};
    const pathWithParams = applyPathParams(path, params ?? {});
    navigate(pathWithParams, rest)
  };
}