import { z } from "zod";

export const driverInfoSchema = z.object({
  key: z.string(),
  name: z.string(),
  description: z.string(),
  icon: z.string().optional(),
  tags: z.array(z.string()).optional(),
  configuration_schema: z.record(z.string(), z.unknown()).optional(),
});

export type DriverInfo = z.infer<typeof driverInfoSchema>;
