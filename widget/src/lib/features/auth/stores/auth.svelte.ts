const AUTH_STORAGE_KEY = "snipet@auth-token";
export type Auth = {
  isAuthenticated: boolean;
  token: string | null;
  loading: boolean;
}
export const auth = $state<Auth>({
  isAuthenticated: false,
  token: null,
  loading: true,
});

export function loadAuth() {
  const token = sessionStorage.getItem(AUTH_STORAGE_KEY);

  auth.loading = false;
  auth.token = token;
  auth.isAuthenticated = !!token;
}

export function login(token: string) {
  sessionStorage.setItem(AUTH_STORAGE_KEY, token);

  auth.loading = false;
  auth.token = token;
  auth.isAuthenticated = true;
}

export function logout() {
  sessionStorage.removeItem(AUTH_STORAGE_KEY);

  auth.token = null;
  auth.isAuthenticated = false;
}
