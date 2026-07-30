import { http } from "@/lib/http";

import { paginatedKnowledgeSchema } from "./schemas";

import type { PaginatedKnowledge } from "./schemas";
import type { ServiceGetOptions } from "@/lib/services";

const KNOWLEDGE_URL = "/api/knowledge";

const list = async (
  opts?: ServiceGetOptions<PaginatedKnowledge>,
): Promise<PaginatedKnowledge> => {
  return http.get({
    url: KNOWLEDGE_URL,
    schemas: {
      response: paginatedKnowledgeSchema,
    },
    ...opts,
  });
};

export const knowledgeService = {
  list,
};
