import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

export const sessionMetadataSchema = z.object({
  name: z.string().optional(),
}).catchall(z.unknown());

export type SessionMetadata = z.infer<typeof sessionMetadataSchema>;

/** Nested agent when `include=agent` is requested. */
export const sessionAgentSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  description: z.string(),
  instructions: z.string(),
  configuration: z.unknown(),
}).loose();

export type SessionAgent = z.infer<typeof sessionAgentSchema>;

export const sessionSchema = z
  .object({
    id: z.uuid(),
    client_id: z.uuid(),
    agent_id: z.uuid(),
    metadata: sessionMetadataSchema,
    agent: sessionAgentSchema.optional(),
  })
  .strict();

export type Session = z.infer<typeof sessionSchema>;

export const paginatedSessionSchema = paginatedSchema(sessionSchema);
export type PaginatedSession = z.infer<typeof paginatedSessionSchema>;

export const createSessionSchema = z
  .object({
    agent_id: z.uuid(),
    metadata: sessionMetadataSchema.optional(),
  })
  .strict();

export type CreateSession = z.infer<typeof createSessionSchema>;

export const updateSessionSchema = z
  .object({
    agent_id: z.uuid().optional(),
    metadata: sessionMetadataSchema.optional(),
  })
  .strict();

export type UpdateSession = z.infer<typeof updateSessionSchema>;

export const runSessionSchema = z
  .object({
    message: z.string().min(1),
  })
  .strict();

export type RunSession = z.infer<typeof runSessionSchema>;

export const messageRoleSchema = z.enum([
  "system",
  "user",
  "assistant",
  "tool",
  "final",
]);

export type MessageRole = z.infer<typeof messageRoleSchema>;

export const executionMessageSchema = z
  .object({
    id: z.uuid(),
    execution_id: z.uuid(),
    sequence: z.number(),
    role: messageRoleSchema,
    content: z.string(),
    tool_calls: z.array(z.unknown()).optional(),
    tool_result: z.unknown().nullable().optional(),
    created_at: z.coerce.date(),
  })
  .strict();

export type ExecutionMessage = z.infer<typeof executionMessageSchema>;

export const paginatedExecutionMessageSchema = paginatedSchema(executionMessageSchema);
export type PaginatedExecutionMessage = z.infer<typeof paginatedExecutionMessageSchema>;

export const sessionIncludeSchema = z.enum(["agent"]);
export type SessionInclude = z.infer<typeof sessionIncludeSchema>;

export const listSessionSearchParamsSchema = z
  .object({
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
    include: z.union([sessionIncludeSchema, z.array(sessionIncludeSchema)]).optional(),
  })
  .strict();

export type ListSessionSearchParams = z.infer<typeof listSessionSearchParamsSchema>;

export const findSessionSearchParamsSchema = z
  .object({
    include: z.union([sessionIncludeSchema, z.array(sessionIncludeSchema)]).optional(),
  })
  .strict();

export type FindSessionSearchParams = z.infer<typeof findSessionSearchParamsSchema>;

export const listMessagesSearchParamsSchema = z
  .object({
    sort: z.enum(["asc", "desc"]).optional(),
    take: z.number().min(1).optional(),
    skip: z.number().min(0).optional(),
  })
  .strict();

export type ListMessagesSearchParams = z.infer<typeof listMessagesSearchParamsSchema>;
