import { z } from "zod";

import { tenantSchema, type Tenant } from "@/models/tenant";

export interface Llm {
  id: string;
  tenant_id: string;
  name: string;
  provider: string;
  configuration: Record<string, unknown>;
  tenant?: Tenant | null;
}

export const llmSchema: z.ZodType<Llm> = z.lazy(() =>
  z
    .object({
      id: z.uuid(),
      tenant_id: z.uuid(),
      name: z.string().min(1).max(255),
      provider: z.string().min(1).max(255),
      configuration: z.record(z.string(), z.unknown()),
      tenant: tenantSchema.nullable().optional(),
    })
    .strict(),
);
