import { llmBaseSchema, llmSchema } from "@/models/llm";
import { driverInfoSchema } from "@/schemas/driver";
import { paginatedSchema, paginationParamsSchema } from "@/schemas/paginated";
import { z } from "zod";

export { llmSchema } from "@/models/llm";
export type { Llm } from "@/models/llm";

export const paginatedLlmSchema = paginatedSchema(llmSchema);
export type PaginatedLlm = z.infer<typeof paginatedLlmSchema>;

export const listLlmSearchParamsSchema = paginationParamsSchema;
export type ListLlmSearchParams = z.infer<typeof listLlmSearchParamsSchema>;

export const createLlmSchema = llmBaseSchema.pick({
  name: true,
  provider: true,
  configuration: true,
}).strict();

export type CreateLlm = z.infer<typeof createLlmSchema>;

export const updateLlmSchema = createLlmSchema.partial().strict();
export type UpdateLlm = z.infer<typeof updateLlmSchema>;

export const listDriversSchema = z.array(driverInfoSchema);
export type ListDrivers = z.infer<typeof listDriversSchema>;
