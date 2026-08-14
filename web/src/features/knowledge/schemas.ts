import {
  knowledgeIndexSchema,
  knowledgeItemSchema,
  knowledgeSchema,
} from "@/models/knowledge";
import { driverInfoSchema } from "@/schemas/driver";
import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export {
  knowledgeIndexSchema,
  knowledgeItemKindSchema,
  knowledgeItemSchema,
  knowledgeSchema,
  syncStatusSchema,
} from "@/models/knowledge";
export type {
  Knowledge,
  KnowledgeIndex,
  KnowledgeItem,
  KnowledgeItemKind,
  SyncStatus,
} from "@/models/knowledge";

export const paginatedKnowledgeSchema = paginatedSchema(knowledgeSchema);
export type PaginatedKnowledge = z.infer<typeof paginatedKnowledgeSchema>;

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
  })
  .strict();

export type CreateKnowledgeResponse = z.infer<typeof createKnowledgeResponseSchema>;

export const listKnowledgeDriversSchema = z
  .object({
    source_drivers: z.array(driverInfoSchema),
  })
  .strict();

export type ListKnowledgeDrivers = z.infer<typeof listKnowledgeDriversSchema>;

export const paginatedKnowledgeIndexSchema = paginatedSchema(knowledgeIndexSchema);
export type PaginatedKnowledgeIndex = z.infer<typeof paginatedKnowledgeIndexSchema>;

export const createKnowledgeIndexSchema = z
  .object({
    name: z.string().min(1).max(255),
    driver: z.string().min(1).max(100),
    configuration: z.record(z.string(), z.unknown()),
  })
  .strict();

export type CreateKnowledgeIndex = z.infer<typeof createKnowledgeIndexSchema>;

export const updateKnowledgeIndexSchema = z
  .object({
    name: z.string().min(1).max(255),
  })
  .strict();

export type UpdateKnowledgeIndex = z.infer<typeof updateKnowledgeIndexSchema>;

export const listKnowledgeIndexDriversSchema = z
  .object({
    index_drivers: z.array(driverInfoSchema),
  })
  .strict();

export type ListKnowledgeIndexDrivers = z.infer<typeof listKnowledgeIndexDriversSchema>;
