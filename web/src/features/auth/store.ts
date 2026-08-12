import { create } from "zustand";

import { authenticateResponseSchema } from "./schemas";

import type { AuthenticateResponse } from "./schemas";

const KEYS = {
  accessToken: "snipet@access-token",
  accessTokenExpiresAt: "snipet@access-token-expires-at",
  refreshToken: "snipet@refresh-token",
  refreshTokenExpiresAt: "snipet@refresh-token-expires-at",
} as const;

function readDate(key: string): Date | null {
  const raw = localStorage.getItem(key);
  if (!raw) return null;
  const date = new Date(raw);
  return Number.isNaN(date.getTime()) ? null : date;
}

function writeDate(key: string, value: Date) {
  localStorage.setItem(key, value.toISOString());
}

type AuthStore = {
  accessToken: string | null;
  accessTokenExpiresAt: Date | null;
  refreshToken: string | null;
  refreshTokenExpiresAt: Date | null;
  setTokens: (tokens: AuthenticateResponse) => void;
  clear: () => void;
};

export const useAuthStore = create<AuthStore>((set) => ({
  accessToken: localStorage.getItem(KEYS.accessToken),
  accessTokenExpiresAt: readDate(KEYS.accessTokenExpiresAt),
  refreshToken: localStorage.getItem(KEYS.refreshToken),
  refreshTokenExpiresAt: readDate(KEYS.refreshTokenExpiresAt),

  setTokens: (tokens) => {
    const parsed = authenticateResponseSchema.parse(tokens);

    localStorage.setItem(KEYS.accessToken, parsed.access_token);
    writeDate(KEYS.accessTokenExpiresAt, parsed.access_token_expires_at);
    localStorage.setItem(KEYS.refreshToken, parsed.refresh_token);
    writeDate(KEYS.refreshTokenExpiresAt, parsed.refresh_token_expires_at);

    set({
      accessToken: parsed.access_token,
      accessTokenExpiresAt: parsed.access_token_expires_at,
      refreshToken: parsed.refresh_token,
      refreshTokenExpiresAt: parsed.refresh_token_expires_at,
    });
  },

  clear: () => {
    for (const key of Object.values(KEYS)) {
      localStorage.removeItem(key);
    }

    set({
      accessToken: null,
      accessTokenExpiresAt: null,
      refreshToken: null,
      refreshTokenExpiresAt: null,
    });
  },
}));
