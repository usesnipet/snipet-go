import { paginatedSchema } from "@/schemas/paginated";
import { z } from "zod";

import { agentSchema } from "../agent/schemas";

export const sessionMetadataSchema = z.object({
  name: z.string().optional(),
}).catchall(z.unknown());

export type SessionMetadata = z.infer<typeof sessionMetadataSchema>;

export const sessionSchema = z
  .object({
    id: z.uuid(),
    client_id: z.uuid(),
    agent_id: z.uuid(),
    metadata: sessionMetadataSchema,
    agent: agentSchema.optional(),
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

export const messageRoleSchema = z.enum(["system", "user", "assistant", "tool"]);

export type MessageRole = z.infer<typeof messageRoleSchema>;

export const toolCallSchema = z
  .object({
    id: z.string(),
    tool: z.string(),
    arguments: z.record(z.string(), z.unknown()),
  })
  .strict();

export type ToolCall = z.infer<typeof toolCallSchema>;

export const messageSchema = z
  .object({
    id: z.uuid(),
    sequence: z.number(),
    role: messageRoleSchema,
    content: z.string(),
    final: z.boolean().optional(),
    tool_calls: z.array(toolCallSchema).optional(),
    tool_call_id: z.string().optional(),
    timestamp: z.coerce.date(),
  })
  .strict();

export type Message = z.infer<typeof messageSchema>;

export const executionMessageSchema = messageSchema
  .extend({
    execution_id: z.uuid(),
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

export const executionStatusSchema = z.enum([
  "pending",
  "running",
  "completed",
  "failed",
  "max_turns",
  "cancelled",
]);

export type ExecutionStatus = z.infer<typeof executionStatusSchema>;

/** SSE event names from internal/module/agent/subscriber/sse.go */
export const SSE_EVENT = {
  EXECUTION_MESSAGE_ADDED: "execution_message_added",
  EXECUTION_STATUS_CHANGED: "execution_status_changed",
  EXECUTION_TURN_COMPLETED: "execution_turn_completed",
  EXECUTION_FINISHED: "execution_finished",
  EXECUTION_MESSAGE_DELTA: "execution_message_delta",
  EXECUTION_TOOL_CALL_STARTED: "execution_tool_call_started",
  EXECUTION_TOOL_RESULT: "execution_tool_result",
  ERROR: "error",
} as const;

export type SseEventName = (typeof SSE_EVENT)[keyof typeof SSE_EVENT];

export const executionStatusChangedEventSchema = z
  .object({
    status: executionStatusSchema,
    error_message: z.string().optional(),
  })
  .strict();

export type ExecutionStatusChangedEvent = z.infer<
  typeof executionStatusChangedEventSchema
>;

export const executionMessageAddedEventSchema = z
  .object({
    message: messageSchema,
  })
  .strict();

export type ExecutionMessageAddedEvent = z.infer<
  typeof executionMessageAddedEventSchema
>;

export const executionTurnCompletedEventSchema = z
  .object({
    turn: z.number(),
  })
  .strict();

export type ExecutionTurnCompletedEvent = z.infer<
  typeof executionTurnCompletedEventSchema
>;

export const executionFinishedEventSchema = z
  .object({
    status: z.string(),
  })
  .strict();

export type ExecutionFinishedEvent = z.infer<typeof executionFinishedEventSchema>;

export const executionMessageDeltaEventSchema = z
  .object({
    message_id: z.uuid(),
    content: z.string(),
  })
  .strict();

export type ExecutionMessageDeltaEvent = z.infer<
  typeof executionMessageDeltaEventSchema
>;

export const executionToolCallStartedEventSchema = z
  .object({
    tool_call_id: z.string(),
    tool: z.string(),
    arguments: z.record(z.string(), z.unknown()),
  })
  .strict();

export type ExecutionToolCallStartedEvent = z.infer<
  typeof executionToolCallStartedEventSchema
>;

export const executionToolResultEventSchema = z
  .object({
    tool_call_id: z.string(),
    tool: z.string(),
    result: z.string().optional(),
    error: z.string().optional(),
  })
  .strict();

export type ExecutionToolResultEvent = z.infer<
  typeof executionToolResultEventSchema
>;

export const sseErrorEventSchema = z
  .object({
    message: z.string(),
  })
  .strict();

export type SseErrorEvent = z.infer<typeof sseErrorEventSchema>;

