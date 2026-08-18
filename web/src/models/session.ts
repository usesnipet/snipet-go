import { z } from "zod";

import { agentSchema, type Agent } from "@/models/agent";
import { clientSchema, type Client } from "@/models/client";
import { clientUserToSessionSchema, type ClientUserToSession } from "@/models/client-user";
import { executionSchema, type Execution } from "@/models/execution";
import { tenantSchema, type Tenant } from "@/models/tenant";

export const sessionMetadataSchema = z
  .object({ name: z.string().optional() })
  .catchall(z.unknown());

export type SessionMetadata = z.infer<typeof sessionMetadataSchema>;

export interface Session {
  id: string;
  tenant_id: string;
  client_id: string;
  agent_id: string;
  metadata: SessionMetadata;
  tenant?: Tenant | null;
  client?: Client | null;
  agent?: Agent | null;
  client_user_to_sessions: ClientUserToSession[] | null;
  executions: Execution[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const sessionBaseSchema = z
  .object({
    id: z.uuid(),
    tenant_id: z.uuid(),
    client_id: z.uuid(),
    agent_id: z.uuid(),
    metadata: sessionMetadataSchema,
  })
  .strict();

export const sessionSchema: z.ZodType<Session> = z.lazy(() =>
  sessionBaseSchema
    .extend({
      tenant: tenantSchema.nullable().optional(),
      client: clientSchema.nullable().optional(),
      agent: agentSchema.nullable().optional(),
      client_user_to_sessions: z.array(clientUserToSessionSchema).nullable(),
      executions: z.array(executionSchema).nullable(),
    })
    .strict(),
);

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

export interface ExecutionMessage extends Message {
  tenant_id: string;
  execution_id: string;
  tenant?: Tenant | null;
  execution?: Execution | null;
}

export const executionMessageSchema: z.ZodType<ExecutionMessage> = z.lazy(() =>
  messageSchema
    .extend({
      tenant_id: z.uuid(),
      execution_id: z.uuid(),
      tenant: tenantSchema.nullable().optional(),
      execution: executionSchema.nullable().optional(),
    })
    .strict(),
);
