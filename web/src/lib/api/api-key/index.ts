import { authenticatedClient, publicClient } from "$lib/http/client";
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

const meSchema = z.object({ apiKey: z.string().trim().min(30).startsWith("sn_") });
type MeSchema = z.infer<typeof meSchema>;

const createSchema = z.object({
  name: z.string().min(1).max(255),
  expires_at: z.date().optional(),
});
type CreateSchema = z.infer<typeof createSchema>;

const updateExpirationSchema = z.object({
  expires_at: z.date().optional(),
});
type UpdateExpSchema = z.infer<typeof updateExpirationSchema>;

export default {
  me: {
    schema: meSchema,
    run: ({ apiKey }: MeSchema) => publicClient().get<APIKeyWithSecret>({ url: `${API_KEY_URL}/me`, headers: { "X-API-Key": apiKey } }),
  },
  list: {
    run: () => authenticatedClient().get<APIKey[]>({ url: API_KEY_URL }),
  },
  create: {
    schema: createSchema,
    run: (data: CreateSchema) => authenticatedClient().post<APIKeyWithSecret>({ url: API_KEY_URL, body: data }),
  },
  updateExpiration: {
    schema: updateExpirationSchema,
    run: (id: string, body: UpdateExpSchema) => authenticatedClient().put<APIKey>({ url: `${API_KEY_URL}/{id}/expiration`, params: { id }, body: body }),
  }
}