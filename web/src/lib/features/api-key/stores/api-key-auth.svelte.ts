const API_KEY_STORAGE_KEY = "snipet@api-key";

export const apiKeyAuth = $state({
  isAuthenticated: false,
  apiKey: null as string | null,
  loading: true,
})

export function loadAuth() {
  const apiKey = sessionStorage.getItem(API_KEY_STORAGE_KEY);

  apiKeyAuth.loading = false;
  apiKeyAuth.apiKey = apiKey;
  apiKeyAuth.isAuthenticated = !!apiKey;
}

export function login(apiKey: string) {
  sessionStorage.setItem(API_KEY_STORAGE_KEY, apiKey);

  apiKeyAuth.loading = false;
  apiKeyAuth.apiKey = apiKey;
  apiKeyAuth.isAuthenticated = true;
}

export function logout() {
  sessionStorage.removeItem(API_KEY_STORAGE_KEY);

  apiKeyAuth.apiKey = null;
  apiKeyAuth.isAuthenticated = false;
}
