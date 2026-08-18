import { z } from "zod";

export const syncStatusSchema = z.enum(["pending", "in_progress", "failed", "success"]);
export type SyncStatus = z.infer<typeof syncStatusSchema>;

export const knowledgeSchema = z
  .object({
    id: z.uuid(),
    tenant_id: z.uuid(),
    name: z.string().min(1).max(255),
    description: z.string(),
    driver: z.string().min(1).max(100),
    configuration: z.record(z.string(), z.unknown()),
    last_synced_at: z.coerce.date().nullable(),
    sync_status: syncStatusSchema.nullable(),
    sync_error: z.string().nullable(),
  })
  .strict();

export type Knowledge = z.infer<typeof knowledgeSchema>;

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

export const knowledgeItemSchema = z
  .object({
    id: z.uuid(),
    tenant_id: z.uuid(),
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
  })
  .strict();

export type KnowledgeItem = z.infer<typeof knowledgeItemSchema>;

export const knowledgeIndexSchema = z
  .object({
    id: z.uuid(),
    tenant_id: z.uuid(),
    name: z.string().min(1).max(255),
    driver: z.string().min(1).max(100),
    configuration: z.record(z.string(), z.unknown()),
    knowledge_id: z.uuid(),
  })
  .strict();

export type KnowledgeIndex = z.infer<typeof knowledgeIndexSchema>;
