import { authenticatedClient } from "$lib/http/client";
import { queryClient } from "$lib/query-client";
import { createMutation, createQuery } from "@tanstack/svelte-query";
import { toast } from "svelte-sonner";

import {
	agentPaginatedSchema,
	agentSchema,
	createAgentSchema,
	filterAgentSchema,
	updateAgentSchema,
} from "./schemas";

import type {
	Agent,
	AgentPaginated,
	CreateAgent,
	FilterAgent,
	UpdateAgent,
} from "./schemas";

const BASE_URL = "/api/agent";

const listQueryKey = (filter?: FilterAgent) => [BASE_URL, "list", filter];
const findByIdQueryKey = (id: string) => [BASE_URL, id];

export const agentService = {
	queryKeys: {
		listQueryKey,
		findByIdQueryKey,
	},
	list: (filter?: FilterAgent) =>
		createQuery(() => ({
			queryKey: listQueryKey(filter),
			queryFn: async () => {
				const res = await authenticatedClient().get<AgentPaginated>({
					url: BASE_URL,
					searchParams: filter,
					schemas: {
						response: agentPaginatedSchema,
						searchParams: filterAgentSchema,
					},
				});
				return res.data;
			},
		})),
	findById: (id: string) =>
		createQuery(() => ({
			queryKey: findByIdQueryKey(id),
			queryFn: () =>
				authenticatedClient().get<Agent>({
					url: `${BASE_URL}/${id}`,
					schemas: { response: agentSchema },
				}),
		})),
	create: () =>
		createMutation(() => ({
			mutationFn: (data: CreateAgent) =>
				authenticatedClient().post<Agent>({
					url: BASE_URL,
					body: data,
					schemas: { body: createAgentSchema, response: agentSchema },
				}),
			onSuccess: () => {
				toast.success("Agent created.");
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
	update: () =>
		createMutation(() => ({
			mutationFn: ({ data, id }: { id: string; data: UpdateAgent }) =>
				authenticatedClient().put({
					url: `${BASE_URL}/${id}`,
					body: data,
					schemas: { body: updateAgentSchema },
				}),
			onSuccess: (_, { id }) => {
				toast.success("Agent updated.");
				queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
	delete: () =>
		createMutation(() => ({
			mutationFn: (id: string) =>
				authenticatedClient().delete({
					url: `${BASE_URL}/${id}`,
				}),
			onSuccess: (_, id) => {
				toast.success("Agent deleted.");
				queryClient.invalidateQueries({ queryKey: findByIdQueryKey(id) });
				queryClient.invalidateQueries({ queryKey: listQueryKey() });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
};
