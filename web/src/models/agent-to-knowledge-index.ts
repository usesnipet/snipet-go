import { z } from "zod";

import { agentSchema, type Agent } from "@/models/agent";
import { knowledgeIndexSchema, type KnowledgeIndex } from "@/models/knowledge";

/**
 * AgentID/IndexID have no `json` tag on the Go struct (internal/model/agent_to_knowledge_index.go),
 * so they serialize as Go's default field-name casing instead of snake_case.
 */
export interface AgentToKnowledgeIndex {
  AgentID: string;
  IndexID: string;
  agent?: Agent | null;
  index?: KnowledgeIndex | null;
}

export const agentToKnowledgeIndexSchema: z.ZodType<AgentToKnowledgeIndex> = z
  .object({
    AgentID: z.uuid(),
    IndexID: z.uuid(),
    agent: agentSchema.nullable().optional(),
    index: knowledgeIndexSchema.nullable().optional(),
  })
  .strict();
