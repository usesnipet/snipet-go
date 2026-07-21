import { auth } from "../features/auth/stores/auth.svelte";

import { ApiError, HttpClient } from "./http";

export const authenticatedClient = (): HttpClient => {
  const token = auth.token;
  if (!token) {
    throw new ApiError({
      message: "Unauthorized",
      statusCode: 401
    });
  }

  return new HttpClient({ "Authorization": `Bearer ${token}` })
}

export const publicClient = (): HttpClient => {
  return new HttpClient();
}