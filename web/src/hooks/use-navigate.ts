import type { ROUTES } from "@/routes";
import { useNavigate as useNavigateReactRouter } from "react-router";

export const useNavigate = () => {
  const navigate = useNavigateReactRouter();
  return (path: (typeof ROUTES)[number]) => navigate(path);
}