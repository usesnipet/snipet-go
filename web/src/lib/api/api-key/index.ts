import { httpGet, httpPost, httpPut } from "$lib/http";
import z from "zod";

const API_KEY_URL = "/api/api-key";

export type APIKey = {
  id: string;
  name: string;
  key_id: string;
  active: boolean;
  expires_at: Date | null;
  created_at: Date;
  updated_at: Date;
};

export type APIKeyWithSecret = APIKey & {
  key: string;
};

const createSchema = z.object({
  name: z.string().min(1).max(255),
  expires_at: z.date().optional(),
});
type CreateSchema = z.infer<typeof createSchema>;

const updateExpirationSchema = z.object({
  expires_at: z.date().optional(),
});
type UpdateExpirationSchema = z.infer<typeof updateExpirationSchema>;

export default {
  me: {
    schema: z.object({
      apiKey: z.string(),
    }),
    run: async (apiKey: string) => {
      return httpGet<APIKeyWithSecret>({
        url: `${API_KEY_URL}/me`,
        headers: {
          "X-API-Key": apiKey,
        },
      });
    }
  },
  list: {
    run: async () => {
      return httpGet<APIKey[]>({
        url: API_KEY_URL,
      });
    }
  },
  create: {
    schema: createSchema,
    run: async (data: CreateSchema) => {
      return httpPost<APIKeyWithSecret>({
        url: API_KEY_URL,
        body: data,
      });
    }
  },
  updateExpiration: {
    schema: updateExpirationSchema,
    run: async (id: string, body: UpdateExpirationSchema) => {
      return httpPut<APIKey>({
        url: `${API_KEY_URL}/{id}/expiration`,
        params: { id },
        body: body,
      });
    }
  }
}