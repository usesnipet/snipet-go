import { authenticatedClient } from "$lib/http/client";
import { queryClient } from "$lib/query-client";
import { createMutation, createQuery } from "@tanstack/svelte-query";
import { toast } from "svelte-sonner";

import {
  createSessionSchema, filterSessionByIDSchema, filterSessionSchema, sessionPaginatedSchema, sessionSchema
} from "./schemas";

import type {
	CreateSession,
	FilterSession,
	FilterSessionByID,
	Session,
	SessionPaginated,
} from "./schemas";

const baseUrl = (clientCode: string) => `/api/client/${clientCode}/session`;

const listQueryKey = (clientCode: string) => [
	baseUrl(clientCode),
	"list",
];
const findByIdQueryKey = (clientCode: string, id: string) => [
	baseUrl(clientCode),
	id,
];

export const sessionService = {
	queryKeys: {
		listQueryKey,
		findByIdQueryKey,
	},
	list: (clientCode: string, filter?: FilterSession) =>
		createQuery(() => ({
			queryKey: listQueryKey(clientCode),
			queryFn: async () => {
				const res = await authenticatedClient().get<SessionPaginated>({
					url: baseUrl(clientCode),
					searchParams: filter,
					schemas: {
						response: sessionPaginatedSchema,
						searchParams: filterSessionSchema,
					},
				});
				return res.data;
			},
		})),
	findById: (clientCode: string, id: string, filter?: FilterSessionByID) =>
		createQuery(() => ({
			queryKey: findByIdQueryKey(clientCode, id),
			queryFn: () =>
				authenticatedClient().get<Session>({
					url: `${baseUrl(clientCode)}/${id}`,
					searchParams: filter,
					schemas: { response: sessionSchema, searchParams: filterSessionByIDSchema },
				}),
		})),
	create: (clientCode: string) =>
		createMutation(() => ({
			mutationFn: (data: CreateSession) =>
				authenticatedClient().post<Session>({
					url: baseUrl(clientCode),
					body: data,
					schemas: { body: createSessionSchema, response: sessionSchema },
				}),
			onSuccess: () => {
				toast.success("Session created.");
				queryClient.invalidateQueries({ queryKey: listQueryKey(clientCode) });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
	delete: (clientCode: string) =>
		createMutation(() => ({
			mutationFn: (id: string) =>
				authenticatedClient().delete({
					url: `${baseUrl(clientCode)}/${id}`,
				}),
			onSuccess: (_, id) => {
				toast.success("Session deleted.");
				queryClient.invalidateQueries({ queryKey: findByIdQueryKey(clientCode, id) });
				queryClient.invalidateQueries({ queryKey: listQueryKey(clientCode) });
			},
			onError: (error) => {
				toast.error(error.message);
			},
		})),
};
