import { publicClient } from "$lib/http/client";
import { createQuery } from "@tanstack/svelte-query";

import { appConfigSchema } from "./schemas";

import type { AppConfig } from "./schemas";

const BASE_URL = "/api/app";

const configQueryKey = () => [BASE_URL, "config"];

export const appService = {
	queryKeys: {
		configQueryKey,
	},
	config: () =>
		createQuery(() => ({
			queryKey: configQueryKey(),
			queryFn: () =>
				publicClient().get<AppConfig>({
					url: `${BASE_URL}/config`,
					schemas: { response: appConfigSchema },
				}),
		})),
};
