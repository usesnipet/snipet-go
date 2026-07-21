import { logger } from "$lib/logger.js";
import z from "zod";

import { handleApiError, parseZodErrors } from "./errors.js";

export type ApiMethod = "GET" | "POST" | "PUT" | "DELETE";
export type SearchParamsRecord = Record<
  string,
  string | number | boolean | string[] | number[] | boolean[] | undefined | null
>;

export type PathParamsRecord = Record<string, string | number | boolean>;

export type ApiRequestOptions<TBody = unknown, TResponse = unknown> = {
  method: ApiMethod;
  url: string;
  body?: TBody;
  headers?: Record<string, string>;
  params?: PathParamsRecord;
  searchParams?: SearchParamsRecord;
  schemas?: {
    body?: z.ZodSchema<TBody>;
    response?: z.ZodSchema<TResponse>;
    searchParams?: z.ZodSchema<SearchParamsRecord>;
    pathParams?: z.ZodSchema<PathParamsRecord>;
    headers?: z.ZodSchema<Record<string, string>>;
  }
};

export { ApiError, isApiErrorBody, parseApiErrorResponse } from "./errors.js";
export type { ApiErrorBody, ApiErrorDetails } from "./errors.js";


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
  options: ApiRequestOptions<TBody, TResponse>,
): Promise<TResponse> {
  const { url, method, schemas } = options;
  let { body, headers, params, searchParams } = options;
  const pathUrl = params ? applyPathParams(url, params) : url;
  try {
    if (schemas?.pathParams && params) params = schemas.pathParams.parse(params);
    if (schemas?.searchParams && searchParams) searchParams = schemas.searchParams.parse(searchParams);
    if (schemas?.headers && headers) headers = schemas.headers.parse(headers);
    if (schemas?.body) body = schemas.body.parse(body);
  } catch (error) {
    console.error(error);
    if (error instanceof z.ZodError) throw parseZodErrors(error);
    throw error;
  }

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

  const json = await response.json();
  if (schemas?.response) {
    const parsed = schemas.response.safeParse(json);

    if (parsed.success) return parsed.data;
    logger.warn("Invalid response", parsed.error);
    throw parseZodErrors(parsed.error, "Invalid response");
  }
  return json as TResponse;
}


export type HttpGetOptions<TResponse = unknown> = Omit<ApiRequestOptions<undefined, TResponse>, "method" | "body">;
export function httpGet<TResponse = unknown>(options: HttpGetOptions<TResponse>): Promise<TResponse> {
  return http<TResponse, undefined>({ ...options, method: "GET" });
}

export type HttpPostOptions<TBody = unknown, TResponse = unknown> = Omit<ApiRequestOptions<TBody, TResponse>, "method">;
export function httpPost<TResponse = unknown, TBody = unknown>(options: HttpPostOptions<TBody, TResponse>): Promise<TResponse> {
  return http<TResponse, TBody>({ ...options, method: "POST" });
}

export type HttpPutOptions<TBody = unknown, TResponse = unknown> = Omit<ApiRequestOptions<TBody, TResponse>, "method">;
export function httpPut<TResponse = unknown, TBody = unknown>(options: HttpPutOptions<TBody, TResponse>): Promise<TResponse> {
  return http<TResponse, TBody>({ ...options, method: "PUT" });
}

export type HttpDeleteOptions<TResponse = unknown> = Omit<ApiRequestOptions<undefined, TResponse>, "method" | "body">;
export function httpDelete<TResponse = unknown>(options: HttpDeleteOptions<TResponse>): Promise<TResponse> {
  return http<TResponse, undefined>({ ...options, method: "DELETE" });
}

export class HttpClient {
  private readonly defaultHeaders: Record<string, string>;

  constructor(defaultHeaders: Record<string, string> = {}) {
    this.defaultHeaders = defaultHeaders;
  }

  get = <T = unknown>(options: HttpGetOptions<T>) => httpGet<T>({ ...options, headers: { ...this.defaultHeaders, ...options.headers } });
  post = <T = unknown, B = unknown>(options: HttpPostOptions<B, T>) => httpPost<T, B>({ ...options, headers: { ...this.defaultHeaders, ...options.headers } });
  put = <T = unknown, B = unknown>(options: HttpPutOptions<B, T>) => httpPut<T, B>({ ...options, headers: { ...this.defaultHeaders, ...options.headers } });
  delete = <T = unknown>(options: HttpDeleteOptions<T>) => httpDelete<T>({ ...options, headers: { ...this.defaultHeaders, ...options.headers } });

  withHeaders(headers: Record<string, string>): HttpClient {
    return new HttpClient({ ...this.defaultHeaders, ...headers });
  }
}