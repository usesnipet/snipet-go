import { z } from "zod";

import { knowledgeIndexSchema, knowledgeItemSchema, type KnowledgeIndex, type KnowledgeItem } from "@/models/knowledge";

export const indexStatusSchema = z.enum(["pending", "syncing", "indexed", "skipped", "error"]);
export type IndexStatus = z.infer<typeof indexStatusSchema>;

export interface IndexedKnowledgeItem {
  id: string;
  hash: string;
  indexed_at?: Date;
  metadata: Record<string, unknown>;
  status: IndexStatus;
  reason?: string;
  last_error?: string;
  index_id: string;
  knowledge_item_id?: string;
  index?: KnowledgeIndex | null;
  knowledge_item?: KnowledgeItem | null;
}

export const indexedKnowledgeItemSchema: z.ZodType<IndexedKnowledgeItem> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      hash: z.string(),
      indexed_at: z.coerce.date().optional(),
      metadata: z.record(z.string(), z.unknown()),
      status: indexStatusSchema,
      reason: z.string().optional(),
      last_error: z.string().optional(),
      index_id: z.uuid(),
      knowledge_item_id: z.uuid().optional(),
      index: knowledgeIndexSchema.nullable().optional(),
      knowledge_item: knowledgeItemSchema.nullable().optional(),
    })
    .strict(),
);
