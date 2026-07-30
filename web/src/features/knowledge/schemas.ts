import { driverInfoSchema } from "@/schemas/driver";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const syncStatusSchema = z.enum(["pending", "in_progress", "failed", "success"]);
export type SyncStatus = z.infer<typeof syncStatusSchema>;

export const knowledgeSchema = z
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

export type Knowledge = z.infer<typeof knowledgeSchema>;

export const paginatedKnowledgeSchema = paginatedSchema(knowledgeSchema);
export type PaginatedKnowledge = z.infer<typeof paginatedKnowledgeSchema>;

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

export const paginatedKnowledgeItemSchema = paginatedSchema(knowledgeItemSchema);
export type PaginatedKnowledgeItem = z.infer<typeof paginatedKnowledgeItemSchema>;

export const createKnowledgeSchema = z
  .object({
    name: z.string().min(1).max(255),
    description: z.string(),
    driver: z.string().min(1).max(255),
    configuration: z.record(z.string(), z.unknown()),
  })
  .strict();

export type CreateKnowledge = z.infer<typeof createKnowledgeSchema>;

export const updateKnowledgeSchema = z
  .object({
    name: z.string().min(1).max(255),
    description: z.string(),
  })
  .strict();

export type UpdateKnowledge = z.infer<typeof updateKnowledgeSchema>;

export const createKnowledgeResponseSchema = z
  .object({
    knowledge: knowledgeSchema,
    job_id: z.number().int(),
  })
  .strict();

export type CreateKnowledgeResponse = z.infer<typeof createKnowledgeResponseSchema>;

export const listKnowledgeDriversSchema = z
  .object({
    source_drivers: z.array(driverInfoSchema),
  })
  .strict();

export type ListKnowledgeDrivers = z.infer<typeof listKnowledgeDriversSchema>;

export const syncKnowledgeResponseSchema = z
  .object({
    job_id: z.number().int(),
  })
  .strict();

export type SyncKnowledgeResponse = z.infer<typeof syncKnowledgeResponseSchema>;
