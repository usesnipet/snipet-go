import { authenticatedClient } from "$lib/http/client";
import { queryClient } from "$lib/query-client";
import { createMutation, createQuery } from "@tanstack/svelte-query";
import { toast } from "svelte-sonner";

import {
	clientPaginatedSchema,
	clientSchema,
	createClientSchema,
	filterClientSchema,
	updateClientSchema,
} from "./schemas";

import type {
	Client,
	ClientPaginated,
	CreateClient,
	FilterClient,
	UpdateClient,
} from "./schemas";

const BASE_URL = "/api/clients";

const listQueryKey = (filter?: FilterClient) => [BASE_URL, "list", filter];
const findByCodeQueryKey = (code: string) => [BASE_URL, code];

export const clientService = {
	queryKeys: {
		listQueryKey,
		findByCodeQueryKey,
	},
	list: (filter?: FilterClient) =>
		createQuery(() => ({
			queryKey: listQueryKey(filter),
			queryFn: async () => {
				const res = await authenticatedClient().get<ClientPaginated>({
					url: BASE_URL,
					searchParams: filter,
					schemas: {
						response: clientPaginatedSchema,
						searchParams: filterClientSchema,
					},
				});
				return res.data;
			},
		})),
	findByCode: (code: string) =>
		createQuery(() => ({
			queryKey: findByCodeQueryKey(code),
			queryFn: () =>
				authenticatedClient().get<Client>({
					url: `${BASE_URL}/${code}`,
					schemas: { response: clientSchema },
				}),
		})),
	create: () =>
		createMutation(() => ({
			mutationFn: (data: CreateClient) =>
				authenticatedClient().post<Client>({
					url: BASE_URL,
					body: data,
					schemas: { body: createClientSchema, response: clientSchema },
				}),
			onSuccess: () => {
				toast.success("Client created.");
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
	update: () =>
		createMutation(() => ({
			mutationFn: ({ data, code }: { code: string; data: UpdateClient }) =>
				authenticatedClient().put({
					url: `${BASE_URL}/${code}`,
					body: data,
					schemas: { body: updateClientSchema },
				}),
			onSuccess: (_, { code }) => {
				toast.success("Client updated.");
				queryClient.invalidateQueries({ queryKey: findByCodeQueryKey(code) });
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
	delete: () =>
		createMutation(() => ({
			mutationFn: (code: string) =>
				authenticatedClient().delete({
					url: `${BASE_URL}/${code}`,
				}),
			onSuccess: (_, code) => {
				toast.success("Client deleted.");
				queryClient.invalidateQueries({ queryKey: findByCodeQueryKey(code) });
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
};
