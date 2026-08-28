import { z } from "zod";

import { indexedKnowledgeItemSchema, type IndexedKnowledgeItem } from "@/models/indexed-knowledge-item";

export const syncStatusSchema = z.enum(["pending", "in_progress", "failed", "success"]);
export type SyncStatus = z.infer<typeof syncStatusSchema>;

export interface Knowledge {
  id: string;
  name: string;
  description: string;
  driver: string;
  configuration: Record<string, unknown>;
  last_synced_at: Date | null;
  sync_status: SyncStatus | null;
  sync_error: string | null;
  items: KnowledgeItem[] | null;
  indexes: KnowledgeIndex[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const knowledgeBaseSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string(),
    driver: z.string().min(1).max(100),
    configuration: z.record(z.string(), z.unknown()),
    last_synced_at: z.coerce.date().nullable(),
    sync_status: syncStatusSchema.nullable(),
    sync_error: z.string().nullable(),
  })
  .strict();

export const knowledgeSchema: z.ZodType<Knowledge> = z.lazy(() =>
  knowledgeBaseSchema
    .extend({
      items: z.array(knowledgeItemSchema).nullable(),
      indexes: z.array(knowledgeIndexSchema).nullable(),
    })
    .strict(),
);

export const knowledgeItemKindSchema = z.enum([
  "text",
  "document",
  "image",
  "audio",
  "video",
  "structured",
  "unknown",
]);
export type KnowledgeItemKind = z.infer<typeof knowledgeItemKindSchema>;

export interface KnowledgeItem {
  id: string;
  external_id: string;
  name: string;
  hash: string;
  metadata: Record<string, unknown>;
  attributes: Record<string, unknown>;
  kind: KnowledgeItemKind;
  last_modified?: Date | null;
  knowledge_id: string;
  knowledge?: Knowledge | null;
  indexes: IndexedKnowledgeItem[] | null;
}

export const knowledgeItemSchema: z.ZodType<KnowledgeItem> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      external_id: z.string(),
      name: z.string(),
      hash: z.string(),
      metadata: z
        .record(z.string(), z.unknown())
        .nullish()
        .transform((value) => value ?? {}),
      attributes: z
        .record(z.string(), z.unknown())
        .nullish()
        .transform((value) => value ?? {}),
      kind: knowledgeItemKindSchema,
      last_modified: z.coerce.date().nullish(),
      knowledge_id: z.uuid(),
      knowledge: knowledgeSchema.nullable().optional(),
      indexes: z.array(indexedKnowledgeItemSchema).nullable(),
    })
    .strict(),
);

export interface KnowledgeIndex {
  id: string;
  name: string;
  driver: string;
  configuration: Record<string, unknown>;
  knowledge_id: string;
  knowledge?: Knowledge | null;
  items: IndexedKnowledgeItem[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const knowledgeIndexBaseSchema = z
  .object({
    id: z.uuid(),
    name: z.string().min(1).max(255),
    driver: z.string().min(1).max(100),
    configuration: z.record(z.string(), z.unknown()),
    knowledge_id: z.uuid(),
  })
  .strict();

export const knowledgeIndexSchema: z.ZodType<KnowledgeIndex> = z.lazy(() =>
  knowledgeIndexBaseSchema
    .extend({
      knowledge: knowledgeSchema.nullable().optional(),
      items: z.array(indexedKnowledgeItemSchema).nullable(),
    })
    .strict(),
);
