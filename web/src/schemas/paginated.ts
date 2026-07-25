import { z } from "zod";

export const paginatedSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.object({
    data: z.array(dataSchema),
    total: z.number(),
    skip: z.number(),
    take: z.number(),
  });

export type Paginated<T> = {
  data: T[];
  total: number;
  skip: number;
  take: number;
};