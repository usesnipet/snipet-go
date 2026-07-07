

const API_KEY_STORAGE_KEY = "snipet@api-key";

export const auth = $state({
  isAuthenticated: false,
  apiKey: null as string | null,
  loading: true,
})

export function loadAuth() {
  const apiKey = sessionStorage.getItem(API_KEY_STORAGE_KEY);

  auth.loading = false;
  auth.apiKey = apiKey;
  auth.isAuthenticated = !!apiKey;
}

export function login(apiKey: string) {
  sessionStorage.setItem(API_KEY_STORAGE_KEY, apiKey);

  auth.loading = false;
  auth.apiKey = apiKey;
  auth.isAuthenticated = true;
}

export function logout() {
  sessionStorage.removeItem(API_KEY_STORAGE_KEY);

  auth.apiKey = null;
  auth.isAuthenticated = false;
}