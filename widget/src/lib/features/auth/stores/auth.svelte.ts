import type { AuthenticateResponse, User } from "../schemas";

const AUTH_STORAGE_KEY = "snipet@auth";

type StoredAuth = {
	accessToken: string;
	refreshToken: string;
	user: User;
};

export type Auth = {
	isAuthenticated: boolean;
	accessToken: string | null;
	refreshToken: string | null;
	user: User | null;
	loading: boolean;
};

export const auth = $state<Auth>({
	isAuthenticated: false,
	accessToken: null,
	refreshToken: null,
	user: null,
	loading: true,
});

function persist(data: StoredAuth | null) {
	if (!data) {
		sessionStorage.removeItem(AUTH_STORAGE_KEY);
		return;
	}
	sessionStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(data));
}

function apply(data: StoredAuth | null) {
	auth.accessToken = data?.accessToken ?? null;
	auth.refreshToken = data?.refreshToken ?? null;
	auth.user = data?.user ?? null;
	auth.isAuthenticated = !!data?.accessToken;
	auth.loading = false;
}

export function loadAuth() {
	const raw = sessionStorage.getItem(AUTH_STORAGE_KEY);
	if (!raw) {
		apply(null);
		return;
	}

	try {
		const parsed = JSON.parse(raw) as StoredAuth;
		if (!parsed.accessToken || !parsed.refreshToken) {
			persist(null);
			apply(null);
			return;
		}
		apply(parsed);
	} catch {
		persist(null);
		apply(null);
	}
}

export function setAuth(response: AuthenticateResponse) {
	const data: StoredAuth = {
		accessToken: response.access_token,
		refreshToken: response.refresh_token,
		user: response.user,
	};
	persist(data);
	apply(data);
}

export function logout() {
	persist(null);
	apply(null);
}
