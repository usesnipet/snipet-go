import z from "zod";

export const knowledgeSyncStatusSchema = z.enum(["in_progress", "failed", "success"]);
export type KnowledgeSyncStatus = z.infer<typeof knowledgeSyncStatusSchema>;

export const knowledgeSchema = z.object({
  id: z.uuid(),
  name: z.string().min(1).max(255),
  description: z.string().optional(),
  driver: z.string().min(1).max(255),
  configuration: z.json(),
  last_synced_at: z.coerce.date().nullable(),
  sync_status: knowledgeSyncStatusSchema.nullable(),
  sync_error: z.string().nullable(),
});
export type Knowledge = z.infer<typeof knowledgeSchema>;

export const knowledgePaginatedSchema = z.object({
  data: knowledgeSchema.array(),
  total: z.number(),
  take: z.number(),
  skip: z.number(),
});
export type KnowledgePaginated = z.infer<typeof knowledgePaginatedSchema>;

export const filterKnowledgeSchema = z.object({
  take: z.number().min(1).optional(),
  skip: z.number().min(0).optional(),
}).default({});
export type FilterKnowledge = z.infer<typeof filterKnowledgeSchema>;

export const createKnowledgeSchema = knowledgeSchema.pick({
  name: true,
  description: true,
  driver: true,
  configuration: true,
});
export type CreateKnowledge = z.infer<typeof createKnowledgeSchema>;

export const updateKnowledgeSchema = z.object({
  name: z.string().min(1).max(255).optional(),
  description: z.string().optional(),
});
export type UpdateKnowledge = z.infer<typeof updateKnowledgeSchema>;

export const testConnectionSchema = z.object({
  driver: z.string().min(1).max(255),
  configuration: z.json(),
});
export type TestConnection = z.infer<typeof testConnectionSchema>;

export const syncKnowledgeQuerySchema = z.object({
  force: z.boolean().optional(),
});
export type SyncKnowledgeQuery = z.infer<typeof syncKnowledgeQuerySchema>;

export const knowledgeItemSchema = z.object({
  id: z.uuid(),
  external_id: z.string(),
  name: z.string(),
  hash: z.string(),
  metadata: z.json(),
  kinds: z.array(z.unknown()),
  last_modified: z.coerce.date().nullable().optional(),
  knowledge_id: z.uuid(),
});
export type KnowledgeItem = z.infer<typeof knowledgeItemSchema>;

export const knowledgeItemPaginatedSchema = z.object({
  data: knowledgeItemSchema.array(),
  total: z.number(),
  take: z.number(),
  skip: z.number(),
});
export type KnowledgeItemPaginated = z.infer<typeof knowledgeItemPaginatedSchema>;

export const filterKnowledgeItemSchema = filterKnowledgeSchema;
export type FilterKnowledgeItem = z.infer<typeof filterKnowledgeItemSchema>;
