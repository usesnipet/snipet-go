import { z } from "zod";

export const sessionSchema = z.object({
  id: z.string(),
  clientId: z.string(),
  agentId: z.string(),
  metadata: z.object({
    name: z.string().optional(),
  }).catchall(z.unknown()),
}).strict();

export type Session = z.infer<typeof sessionSchema>;