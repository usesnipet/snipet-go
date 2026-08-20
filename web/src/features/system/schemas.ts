import z from "zod";

export const systemInfoSchema = z.object({
  version: z.string(),
  multi_tenant_enabled: z.boolean(),
})
export type SystemInfo = z.infer<typeof systemInfoSchema>