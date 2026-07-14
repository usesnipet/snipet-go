import { authenticatedClient, publicClient } from "$lib/http/client";
import { queryClient } from "$lib/query-client";
import { createMutation, createQuery } from "@tanstack/svelte-query";
import { toast } from "svelte-sonner";

import {
	apiKeyPaginatedSchema,
	apiKeySchema,
	apiKeyWithSecretSchema,
	createApiKeySchema,
	updateApiKeyExpirationSchema,
} from "./schemas";

import type {
	APIKey,
	APIKeyPaginated,
	APIKeyWithSecret,
	CreateApiKey,
	UpdateApiKeyExpiration,
} from "./schemas";

const BASE_URL = "/api/api-key";

const listQueryKey = () => ["api-keys", "list"];
const meQueryKey = () => ["api-keys", "me"];

export const apiKeyService = {
	queryKeys: {
		listQueryKey,
		meQueryKey,
	},
	check: () =>
		createMutation(() => ({
			mutationFn: async (apiKey: string) => {
				const res = await publicClient().get<APIKey>({
					url: `${BASE_URL}/me`,
					headers: { "X-API-Key": apiKey },
					schemas: { response: apiKeySchema },
				});
				return !!res;
			},
		})),
	me: () =>
		createQuery(() => ({
			queryKey: meQueryKey(),
			queryFn: () =>
				authenticatedClient().get<APIKey>({
					url: `${BASE_URL}/me`,
					schemas: { response: apiKeySchema },
				}),
		})),
	list: () =>
		createQuery(() => ({
			queryKey: listQueryKey(),
			queryFn: async () => {
				const res = await authenticatedClient().get<APIKeyPaginated>({
					url: `${BASE_URL}`,
					schemas: { response: apiKeyPaginatedSchema },
				});
				return res.data;
			},
		})),
	create: () =>
		createMutation(() => ({
			mutationFn: (data: CreateApiKey) =>
				authenticatedClient().post<APIKeyWithSecret>({
					url: `${BASE_URL}`,
					body: data,
					schemas: { body: createApiKeySchema, response: apiKeyWithSecretSchema },
				}),
			onSuccess: () => {
				toast.success("API key created.");
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
	updateExpiration: () =>
		createMutation(() => ({
			mutationFn: ({ id, data }: { id: string; data: UpdateApiKeyExpiration }) =>
				authenticatedClient().put({
					url: `${BASE_URL}/${id}/expiration`,
					body: data,
					schemas: { body: updateApiKeyExpirationSchema },
				}),
			onSuccess: () => {
				toast.success("API key expiration updated.");
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
				queryClient.invalidateQueries({ queryKey: meQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
};
