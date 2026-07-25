import { z } from "zod";

export const jwtSchema = z.string().startsWith("Bearer ");
export type Jwt = z.infer<typeof jwtSchema>;