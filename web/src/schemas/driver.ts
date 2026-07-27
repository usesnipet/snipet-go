import { z } from "zod";

export const driverInfoSchema = z.object({
  name: z.string(),
  description: z.string(),
  icon: z.string().optional(),
  tags: z.array(z.string()).optional(),
  configurationSchema: z.record(z.string(), z.any()).optional(),
});

export type DriverInfo = z.infer<typeof driverInfoSchema>;