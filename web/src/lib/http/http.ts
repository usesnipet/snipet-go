export type ApiMethod = "GET" | "POST" | "PUT" | "DELETE";
export type SearchParamsRecord = Record<
  string,
  string | number | boolean | undefined | null
>;

export type PathParamsRecord = Record<string, string | number | boolean>;

export type ApiRequestOptions<TBody = unknown> = {
  method: ApiMethod;
  url: string;
  body?: TBody;
  headers?: Record<string, string>;
  params?: PathParamsRecord;
  searchParams?: SearchParamsRecord;
};

export { ApiError, isApiErrorBody, parseApiErrorResponse } from "./errors.js";
export type { ApiErrorBody, ApiErrorDetails } from "./errors.js";

import { handleApiError } from "./errors.js";

export function applyPathParams(url: string, params: PathParamsRecord): string {
  return Object.entries(params).reduce(
    (result, [key, value]) =>
      result.replaceAll(`{${key}}`, encodeURIComponent(String(value))),
    url,
  );
}

export function buildSearchParams(params: SearchParamsRecord): string {
  const searchParams = new URLSearchParams();

  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null) continue;
    searchParams.append(key, String(value));
  }

  return searchParams.toString();
}

export function applySearchParams(
  url: string,
  params: SearchParamsRecord,
): string {
  const query = buildSearchParams(params);
  if (!query) return url;

  const separator = url.includes("?") ? "&" : "?";
  return `${url}${separator}${query}`;
}

export async function http<TResponse = unknown, TBody = unknown>(
  options: ApiRequestOptions<TBody>,
): Promise<TResponse> {
  const { url, method, body, headers, params, searchParams } = options;
  const pathUrl = params ? applyPathParams(url, params) : url;
  const requestUrl = searchParams
    ? applySearchParams(pathUrl, searchParams)
    : pathUrl;

  const response = await fetch(requestUrl, {
    method,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    headers: {
      ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
      ...headers,
    },
  });

  if (!response.ok) {
    await handleApiError(response);
  }

  if (response.status === 204) {
    return undefined as TResponse;
  }

  return response.json() as Promise<TResponse>;
}


export type HttpGetOptions = Omit<ApiRequestOptions<undefined>, "method" | "body">;
export function httpGet<TResponse = unknown>(options: HttpGetOptions): Promise<TResponse> {
  return http<TResponse, undefined>({ ...options, method: "GET" });
}

export type HttpPostOptions<TBody = unknown> = Omit<ApiRequestOptions<TBody>, "method">;
export function httpPost<TResponse = unknown, TBody = unknown>(options: HttpPostOptions<TBody>): Promise<TResponse> {
  return http<TResponse, TBody>({ ...options, method: "POST" });
}

export type HttpPutOptions<TBody = unknown> = Omit<ApiRequestOptions<TBody>, "method">;
export function httpPut<TResponse = unknown, TBody = unknown>(options: HttpPutOptions<TBody>): Promise<TResponse> {
  return http<TResponse, TBody>({ ...options, method: "PUT" });
}

export type HttpDeleteOptions = Omit<ApiRequestOptions<undefined>, "method" | "body">;
export function httpDelete<TResponse = unknown>(options: HttpDeleteOptions): Promise<TResponse> {
  return http<TResponse, undefined>({ ...options, method: "DELETE" });
}