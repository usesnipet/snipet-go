import { z } from "zod";

import { agentSchema, type Agent } from "@/models/agent";
import { executionMessageSchema, sessionSchema, type ExecutionMessage, type Session } from "@/models/session";

export const executionStatusSchema = z.enum([
  "pending",
  "running",
  "completed",
  "failed",
  "max_turns",
  "cancelled",
]);
export type ExecutionStatus = z.infer<typeof executionStatusSchema>;

export interface Execution {
  id: string;
  session_id?: string;
  agent_id: string;
  status: ExecutionStatus;
  error_message?: string;
  turns: number;
  metadata: Record<string, unknown>;
  created_at: Date;
  updated_at: Date;
  session?: Session | null;
  agent?: Agent | null;
  messages: ExecutionMessage[] | null;
}

export const executionSchema: z.ZodType<Execution> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      session_id: z.uuid().optional(),
      agent_id: z.uuid(),
      status: executionStatusSchema,
      error_message: z.string().optional(),
      turns: z.number().int(),
      metadata: z.record(z.string(), z.unknown()),
      created_at: z.coerce.date(),
      updated_at: z.coerce.date(),
      session: sessionSchema.nullable().optional(),
      agent: agentSchema.nullable().optional(),
      messages: z.array(executionMessageSchema).nullable(),
    })
    .strict(),
);
