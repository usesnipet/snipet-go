import { useQuery } from "@tanstack/react-query";

import { knowledgeService } from "./service";

import type { PaginatedKnowledge } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";
import type { UseQueryResult } from "@tanstack/react-query";

const BASE_QUERY_KEY = "knowledge";

export const listKnowledgeQueryKey = () => [BASE_QUERY_KEY] as const;
export const useListKnowledge = (
  opts?: ServiceGetOptions<PaginatedKnowledge>,
): UseQueryResult<PaginatedKnowledge, Error> => {
  return useQuery({
    queryKey: listKnowledgeQueryKey(),
    queryFn: () => knowledgeService.list({ ...opts, auth: "api-key" }),
  });
};
