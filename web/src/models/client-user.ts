import { z } from "zod";

export const clientUserMetadataSchema = z.record(z.string(), z.unknown());
export type ClientUserMetadata = z.infer<typeof clientUserMetadataSchema>;

export const clientUserSchema = z
  .object({
    id: z.uuid(),
    name: z.string(),
    picture: z.string().nullable().optional(),
    email: z.string().nullable().optional(),
    metadata: clientUserMetadataSchema,
  })
  .strict();

export type ClientUser = z.infer<typeof clientUserSchema>;
