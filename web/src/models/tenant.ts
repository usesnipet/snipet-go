import { z } from "zod";

import { agentSchema, type Agent } from "@/models/agent";
import { apiKeySchema, type ApiKey } from "@/models/api-key";
import { appSchema, type App } from "@/models/app";
import { invitationSchema, type Invitation } from "@/models/invitation";
import { knowledgeSchema, type Knowledge } from "@/models/knowledge";
import { llmSchema, type Llm } from "@/models/llm";
import { memberSchema, type Member } from "@/models/member";

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  icon: string | null;
  created_at: Date;
  updated_at: Date;
  members: Member[] | null;
  invitations: Invitation[] | null;
  agents: Agent[] | null;
  api_keys: ApiKey[] | null;
  clients: App[] | null;
  knowledges: Knowledge[] | null;
  llms: Llm[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const tenantBaseSchema = z
  .object({
    id: z.string(),
    name: z.string().min(1).max(255),
    slug: z.string().min(1).max(255),
    icon: z.string().max(255).nullable(),
    created_at: z.coerce.date(),
    updated_at: z.coerce.date(),
  })
  .strict();

export const tenantSchema: z.ZodType<Tenant> = z.lazy(() =>
  tenantBaseSchema
    .extend({
      members: z.array(memberSchema).nullable(),
      invitations: z.array(invitationSchema).nullable(),
      agents: z.array(agentSchema).nullable(),
      api_keys: z.array(apiKeySchema).nullable(),
      clients: z.array(appSchema).nullable(),
      knowledges: z.array(knowledgeSchema).nullable(),
      llms: z.array(llmSchema).nullable(),
    })
    .strict(),
);
