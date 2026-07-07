import { auth } from "$lib/store/auth.svelte";

import { ApiError, HttpClient } from "./http";

export const authenticatedClient = (): HttpClient => {
  const token = auth.apiKey;
  if (!token) {
    throw new ApiError({
      message: "Unauthorized",
      statusCode: 401
    });
  }

  return new HttpClient({ "X-API-Key": token })
}

export const publicClient = (): HttpClient => {
  return new HttpClient();
}