import { z } from "zod";

import { clientToClientUserSchema, type ClientToClientUser } from "@/models/client-user";
import { sessionSchema, type Session } from "@/models/session";
import { tenantSchema, type Tenant } from "@/models/tenant";

/** Accepts "" from form inputs and coerces to undefined; input/output stay `string | undefined`. */
const urlSchema = z.union([
  z.url(),
  z.literal("").transform(() => undefined),
  z.undefined(),
]);

export const clientConfigSchema = z
  .object({
    oidc: z
      .object({
        issuer: urlSchema,
        audience: urlSchema,
        enabled: z.boolean(),
      })
      .strict(),
    webhook: z
      .object({
        url: urlSchema,
        enabled: z.boolean(),
      })
      .strict(),
    anonymous: z
      .object({
        enabled: z.boolean(),
      })
      .strict(),
  })
  .strict();

export type ClientConfig = z.infer<typeof clientConfigSchema>;

export interface Client {
  id: string;
  tenant_id: string;
  code: string;
  name: string;
  config: ClientConfig;
  tenant?: Tenant | null;
  sessions: Session[] | null;
  client_to_users: ClientToClientUser[] | null;
}

/** Own fields only, no relations — pick/extend/partial from this in feature schemas (create/update DTOs). */
export const clientBaseSchema = z
  .object({
    id: z.string(),
    tenant_id: z.string(),
    code: z.string(),
    name: z.string().min(1).max(255),
    config: clientConfigSchema,
  })
  .strict();

export const clientSchema: z.ZodType<Client> = z.lazy(() =>
  clientBaseSchema
    .extend({
      tenant: tenantSchema.nullable().optional(),
      sessions: z.array(sessionSchema).nullable(),
      client_to_users: z.array(clientToClientUserSchema).nullable(),
    })
    .strict(),
);
