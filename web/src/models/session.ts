import { z } from "zod";

import { agentSchema } from "@/models/agent";

export const sessionMetadataSchema = z
  .object({ name: z.string().optional() })
  .catchall(z.unknown());

export type SessionMetadata = z.infer<typeof sessionMetadataSchema>;

export const sessionSchema = z
  .object({
    id: z.uuid(),
    tenant_id: z.uuid(),
    client_id: z.uuid(),
    agent_id: z.uuid(),
    metadata: sessionMetadataSchema,
    agent: agentSchema.optional(),
  })
  .strict();

export type Session = z.infer<typeof sessionSchema>;

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
    tenant_id: z.uuid(),
    execution_id: z.uuid(),
  })
  .strict();

export type ExecutionMessage = z.infer<typeof executionMessageSchema>;
