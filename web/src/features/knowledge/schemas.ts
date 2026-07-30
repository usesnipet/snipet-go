import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const syncStatusSchema = z.enum(["in_progress", "failed", "success"]);
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
