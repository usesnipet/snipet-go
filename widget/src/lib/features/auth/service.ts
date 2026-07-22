import { publicClient } from "$lib/http/client";

import {
	authenticateAnonymousSchema,
	authenticateResponseSchema,
	refreshSchema,
} from "./schemas";

import type { AuthenticateAnonymous, AuthenticateResponse, Refresh } from "./schemas";

const baseUrl = (clientCode: string) => `/api/client/${clientCode}`;

export const authService = {
	authenticateAnonymous: (clientCode: string, data: AuthenticateAnonymous = {}) =>
		publicClient().post<AuthenticateResponse, AuthenticateAnonymous>({
			url: `${baseUrl(clientCode)}/authenticate/anonymous`,
			body: data,
			schemas: {
				body: authenticateAnonymousSchema,
				response: authenticateResponseSchema,
			},
		}),
	refresh: (clientCode: string, data: Refresh) =>
		publicClient().post<AuthenticateResponse, Refresh>({
			url: `${baseUrl(clientCode)}/refresh`,
			body: data,
			schemas: {
				body: refreshSchema,
				response: authenticateResponseSchema,
			},
		}),
};
